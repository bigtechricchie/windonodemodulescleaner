package main

import (
	"flag"
	"fmt"
	"os"
)

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
		"Scanning %s for node_modules directories...\n",
		*rootDirectory,
	)

	scanOutcome, err := findNodeModulesDirectories(*rootDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(scanOutcome.Results) == 0 {
		fmt.Println("No node_modules directories found.")
	} else {
		directoryWord := "directories"

		if len(scanOutcome.Results) == 1 {
			directoryWord = "directory"
		}

		fmt.Printf(
			"Found %d node_modules %s:\n",
			len(scanOutcome.Results),
			directoryWord,
		)

		for _, result := range scanOutcome.Results {
			if result.SizeError != nil {
				fmt.Printf(
					"%s    size unavailable\n",
					result.Path,
				)

				fmt.Fprintf(
					os.Stderr,
					"warning: %v\n",
					result.SizeError,
				)

				continue
			}

			fmt.Printf(
				"%s    %s\n",
				result.Path,
				formatByteSize(result.SizeBytes),
			)
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
