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

	flag.Parse()

	if *rootDirectory == "" {
		fmt.Fprintln(os.Stderr, "error: -root is required")
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

	if len(candidateResults) == 0 {
		fmt.Println("No node_modules directories found.")
	} else {
		directoryWord := "directories"

		if len(candidateResults) == 1 {
			directoryWord = "directory"
		}

		fmt.Printf(
			"Found %d node_modules %s:\n\n",
			len(candidateResults),
			directoryWord,
		)

		for _, result := range candidateResults {
			printCandidateResult(result)
		}
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

func printCandidateResult(result candidateResult) {
	scanResult := result.ScanResult

	if scanResult.SizeError != nil {
		fmt.Printf(
			"%s    size unavailable\n",
			scanResult.Path,
		)

		fmt.Fprintf(
			os.Stderr,
			"warning: %v\n",
			scanResult.SizeError,
		)
	} else {
		fmt.Printf(
			"%s    %s\n",
			scanResult.Path,
			formatByteSize(scanResult.SizeBytes),
		)
	}

	if result.ValidationError != nil {
		fmt.Println("  Cleanup eligibility: rejected")
		fmt.Printf(
			"  Reason: %v\n",
			result.ValidationError,
		)
		fmt.Println()
		return
	}

	fmt.Println("  Cleanup eligibility: currently eligible")
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
