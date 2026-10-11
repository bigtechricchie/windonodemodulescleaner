package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindNodeModulesDirectories(t *testing.T) {
	rootDirectory := t.TempDir()

	firstNodeModules := filepath.Join(
		rootDirectory,
		"project-one",
		"node_modules",
	)

	secondNodeModules := filepath.Join(
		rootDirectory,
		"workspace",
		"apps",
		"dashboard",
		"node_modules",
	)

	nestedNodeModules := filepath.Join(
		firstNodeModules,
		"package-a",
		"node_modules",
	)

	for _, directory := range []string{
		firstNodeModules,
		secondNodeModules,
		nestedNodeModules,
	} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatalf(
				"create test directory %q: %v",
				directory,
				err,
			)
		}
	}

	outcome, err := findNodeModulesDirectories(rootDirectory)
	if err != nil {
		t.Fatalf(
			"find node_modules directories: %v",
			err,
		)
	}

	if len(outcome.Results) != 2 {
		t.Fatalf(
			"expected 2 node_modules directories, got %d",
			len(outcome.Results),
		)
	}

	if outcome.Results[0].Path != firstNodeModules {
		t.Fatalf(
			"expected first directory %q, got %q",
			firstNodeModules,
			outcome.Results[0].Path,
		)
	}

	if outcome.Results[1].Path != secondNodeModules {
		t.Fatalf(
			"expected second directory %q, got %q",
			secondNodeModules,
			outcome.Results[1].Path,
		)
	}

	if len(outcome.Warnings) != 0 {
		t.Fatalf(
			"expected no scan warnings, got %v",
			outcome.Warnings,
		)
	}
}

func TestFindNodeModulesDirectoriesReturnsEmptyResult(t *testing.T) {
	rootDirectory := t.TempDir()

	outcome, err := findNodeModulesDirectories(rootDirectory)
	if err != nil {
		t.Fatalf(
			"find node_modules directories: %v",
			err,
		)
	}

	if len(outcome.Results) != 0 {
		t.Fatalf(
			"expected no node_modules directories, got %v",
			outcome.Results,
		)
	}

	if len(outcome.Warnings) != 0 {
		t.Fatalf(
			"expected no scan warnings, got %v",
			outcome.Warnings,
		)
	}
}

func TestFindNodeModulesDirectoriesRejectsFileAsRoot(t *testing.T) {
	rootDirectory := t.TempDir()

	filePath := filepath.Join(
		rootDirectory,
		"example.txt",
	)

	if err := os.WriteFile(
		filePath,
		[]byte("test"),
		0o600,
	); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	_, err := findNodeModulesDirectories(filePath)
	if err == nil {
		t.Fatal(
			"expected an error when root path is a file",
		)
	}
}

func TestFindNodeModulesDirectoriesRejectsMissingRoot(t *testing.T) {
	rootDirectory := t.TempDir()

	missingDirectory := filepath.Join(
		rootDirectory,
		"does-not-exist",
	)

	_, err := findNodeModulesDirectories(missingDirectory)
	if err == nil {
		t.Fatal(
			"expected an error when root directory does not exist",
		)
	}
}

func TestCalculateDirectorySize(t *testing.T) {
	rootDirectory := t.TempDir()

	firstFile := filepath.Join(
		rootDirectory,
		"first.txt",
	)

	secondDirectory := filepath.Join(
		rootDirectory,
		"nested",
	)

	secondFile := filepath.Join(
		secondDirectory,
		"second.txt",
	)

	if err := os.MkdirAll(secondDirectory, 0o755); err != nil {
		t.Fatalf(
			"create nested test directory: %v",
			err,
		)
	}

	if err := os.WriteFile(
		firstFile,
		[]byte("12345"),
		0o600,
	); err != nil {
		t.Fatalf(
			"create first test file: %v",
			err,
		)
	}

	if err := os.WriteFile(
		secondFile,
		[]byte("1234567890"),
		0o600,
	); err != nil {
		t.Fatalf(
			"create second test file: %v",
			err,
		)
	}

	sizeBytes, err := calculateDirectorySize(rootDirectory)
	if err != nil {
		t.Fatalf(
			"calculate directory size: %v",
			err,
		)
	}

	const expectedSize int64 = 15

	if sizeBytes != expectedSize {
		t.Fatalf(
			"expected size %d bytes, got %d",
			expectedSize,
			sizeBytes,
		)
	}
}
