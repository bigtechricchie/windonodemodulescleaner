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

	nodeModulesDirectories, err := findNodeModulesDirectories(*rootDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(nodeModulesDirectories) == 0 {
		fmt.Println("No node_modules directories found.")
		return
	}

	fmt.Printf(
		"Found %d node_modules directories:\n",
		len(nodeModulesDirectories),
	)

	for _, directory := range nodeModulesDirectories {
		if directory.SizeError != nil {
			fmt.Printf(
				"%s    size unavailable\n",
				directory.Path,
			)

			fmt.Fprintf(
				os.Stderr,
				"warning: %v\n",
				directory.SizeError,
			)

			continue
		}

		fmt.Printf(
			"%s    %s\n",
			directory.Path,
			formatByteSize(directory.SizeBytes),
		)
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
