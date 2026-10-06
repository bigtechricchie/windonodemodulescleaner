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
		fmt.Println(directory)
	}
}
