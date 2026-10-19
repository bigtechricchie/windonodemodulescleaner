package main

import (
	"flag"
	"fmt"
	"os"
)

type candidateResult struct {
	ScanResult      scanResult
	ValidationError error
}

func main() {
	rootDirectory := flag.String(
		"root",
		"",
		"root directory to scan for node_modules directories",
	)

	selectValue := flag.String(
		"select",
		"",
		"comma-separated candidate indexes to select",
	)

	confirmRequested := flag.Bool(
		"confirm",
		false,
		"confirm the selected cleanup candidates",
	)

	deleteRequested := flag.Bool(
		"delete",
		false,
		"delete the confirmed cleanup candidates",
	)

	flag.Parse()

	selectProvided := false

	flag.Visit(func(currentFlag *flag.Flag) {
		if currentFlag.Name == "select" {
			selectProvided = true
		}
	})

	if *rootDirectory == "" {
		fmt.Fprintln(os.Stderr, "error: -root is required")
		os.Exit(1)
	}

	if *confirmRequested && !selectProvided {
		fmt.Fprintln(
			os.Stderr,
			"error: -confirm requires -select",
		)
		os.Exit(1)
	}

	if *deleteRequested && !*confirmRequested {
		fmt.Fprintln(
			os.Stderr,
			"error: -delete requires -confirm",
		)
		os.Exit(1)
	}

	fmt.Printf(
		"Scanning %s for node_modules directories...\n\n",
		*rootDirectory,
	)

	scanOutcome, err := findNodeModulesDirectories(*rootDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	candidateResults := make(
		[]candidateResult,
		0,
		len(scanOutcome.Results),
	)

	for _, result := range scanOutcome.Results {
		candidateResults = append(
			candidateResults,
			candidateResult{
				ScanResult: result,
				ValidationError: validateCleanupTarget(
					*rootDirectory,
					result.Path,
				),
			},
		)
	}

	var selectedCandidates []selectedCandidate

	if selectProvided {
		indexes, err := parseCandidateIndexes(*selectValue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		selectedCandidates, err = selectCandidates(
			candidateResults,
			indexes,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	printCandidateResults(candidateResults)

	if selectProvided {
		printSelectedCandidates(selectedCandidates)

		if *confirmRequested {
			fmt.Println("Selection confirmed.")
			fmt.Println()

			if err := revalidateSelectedCandidates(
				*rootDirectory,
				selectedCandidates,
			); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			candidateWord := "candidates"

			if len(selectedCandidates) == 1 {
				candidateWord = "candidate"
			}

			fmt.Printf(
				"Fresh validation passed for %d %s.\n\n",
				len(selectedCandidates),
				candidateWord,
			)

			if *deleteRequested {
				cleanupResults := cleanupSelectedCandidates(
					*rootDirectory,
					selectedCandidates,
				)

				printCleanupResults(cleanupResults)

				if cleanupFailureCount(cleanupResults) > 0 {
					os.Exit(1)
				}

				return
			}
		}

		fmt.Println("No filesystem changes were made.")
	}

	if len(scanOutcome.Warnings) > 0 {
		fmt.Fprintf(
			os.Stderr,
			"\nScan completed with %d warning(s):\n",
			len(scanOutcome.Warnings),
		)

		for _, warning := range scanOutcome.Warnings {
			fmt.Fprintf(
				os.Stderr,
				"%s: %v\n",
				warning.Path,
				warning.Err,
			)
		}
	}
}

func printCandidateResults(candidateResults []candidateResult) {
	if len(candidateResults) == 0 {
		fmt.Println("No node_modules directories found.")
		return
	}

	directoryWord := "directories"

	if len(candidateResults) == 1 {
		directoryWord = "directory"
	}

	fmt.Printf(
		"Found %d node_modules %s:\n\n",
		len(candidateResults),
		directoryWord,
	)

	for index, result := range candidateResults {
		printCandidateResult(index+1, result)
	}
}

func printCandidateResult(
	index int,
	result candidateResult,
) {
	scanResult := result.ScanResult

	if scanResult.SizeError != nil {
		fmt.Printf(
			"[%d] %s    size unavailable\n",
			index,
			scanResult.Path,
		)

		fmt.Fprintf(
			os.Stderr,
			"warning: %v\n",
			scanResult.SizeError,
		)
	} else {
		fmt.Printf(
			"[%d] %s    %s\n",
			index,
			scanResult.Path,
			formatByteSize(scanResult.SizeBytes),
		)
	}

	if result.ValidationError != nil {
		fmt.Println("    Cleanup eligibility: rejected")
		fmt.Printf(
			"    Reason: %v\n",
			result.ValidationError,
		)
		fmt.Println()
		return
	}

	fmt.Println("    Cleanup eligibility: currently eligible")
	fmt.Println()
}

func printSelectedCandidates(
	selectedCandidates []selectedCandidate,
) {
	candidateWord := "candidates"

	if len(selectedCandidates) == 1 {
		candidateWord = "candidate"
	}

	fmt.Printf(
		"Selected %d %s:\n\n",
		len(selectedCandidates),
		candidateWord,
	)

	for _, selectedCandidate := range selectedCandidates {
		fmt.Printf(
			"[%d] %s\n",
			selectedCandidate.Index,
			selectedCandidate.Candidate.ScanResult.Path,
		)
	}

	fmt.Println()
}

func formatByteSize(sizeBytes int64) string {
	const unit = 1024

	if sizeBytes < unit {
		return fmt.Sprintf("%d B", sizeBytes)
	}

	divisor := int64(unit)
	unitIndex := 0

	for value := sizeBytes / unit; value >= unit; value /= unit {
		divisor *= unit
		unitIndex++
	}

	units := "KMGTPE"

	return fmt.Sprintf(
		"%.1f %ciB",
		float64(sizeBytes)/float64(divisor),
		units[unitIndex],
	)
}
