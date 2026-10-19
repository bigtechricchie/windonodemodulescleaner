package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type cleanupResult struct {
	Index int
	Path  string
	Err   error
}

func cleanupCandidate(
	rootDirectory string,
	candidate selectedCandidate,
) cleanupResult {
	targetDirectory := candidate.Candidate.ScanResult.Path

	result := cleanupResult{
		Index: candidate.Index,
		Path:  targetDirectory,
	}

	if err := validateCleanupTarget(
		rootDirectory,
		targetDirectory,
	); err != nil {
		result.Err = fmt.Errorf(
			"final cleanup validation failed: %w",
			err,
		)

		return result
	}

	if err := os.RemoveAll(targetDirectory); err != nil {
		result.Err = fmt.Errorf(
			"remove cleanup target: %w",
			err,
		)

		return result
	}

	if err := verifyCleanupTargetRemoved(
		targetDirectory,
	); err != nil {
		result.Err = err
	}

	return result
}

func cleanupSelectedCandidates(
	rootDirectory string,
	selectedCandidates []selectedCandidate,
) []cleanupResult {
	results := make(
		[]cleanupResult,
		0,
		len(selectedCandidates),
	)

	for _, candidate := range selectedCandidates {
		results = append(
			results,
			cleanupCandidate(
				rootDirectory,
				candidate,
			),
		)
	}

	return results
}

func verifyCleanupTargetRemoved(
	targetDirectory string,
) error {
	_, err := os.Lstat(targetDirectory)

	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf(
			"verify cleanup target removal: %w",
			err,
		)
	}

	return fmt.Errorf(
		"cleanup target still exists after deletion: %s",
		targetDirectory,
	)
}

func cleanupFailureCount(
	results []cleanupResult,
) int {
	failureCount := 0

	for _, result := range results {
		if result.Err != nil {
			failureCount++
		}
	}

	return failureCount
}

func printCleanupResults(
	results []cleanupResult,
) {
	fmt.Println("Cleanup results:")
	fmt.Println()

	for _, result := range results {
		fmt.Printf(
			"[%d] %s\n",
			result.Index,
			result.Path,
		)

		if result.Err != nil {
			fmt.Printf(
				"    Failed: %v\n",
				result.Err,
			)
		} else {
			fmt.Println("    Deleted")
		}

		fmt.Println()
	}

	failureCount := cleanupFailureCount(results)

	if failureCount == 0 {
		fmt.Println("Cleanup completed successfully.")
		return
	}

	failureWord := "failures"

	if failureCount == 1 {
		failureWord = "failure"
	}

	fmt.Printf(
		"Cleanup completed with %d %s.\n",
		failureCount,
		failureWord,
	)
}
