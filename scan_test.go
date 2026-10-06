package main

import (
	"os"
	"path/filepath"
	"reflect"
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
			t.Fatalf("create test directory %q: %v", directory, err)
		}
	}

	actualDirectories, err := findNodeModulesDirectories(rootDirectory)
	if err != nil {
		t.Fatalf("find node_modules directories: %v", err)
	}

	expectedDirectories := []string{
		firstNodeModules,
		secondNodeModules,
	}

	if !reflect.DeepEqual(actualDirectories, expectedDirectories) {
		t.Fatalf(
			"unexpected directories:\nwant: %v\ngot:  %v",
			expectedDirectories,
			actualDirectories,
		)
	}
}

func TestFindNodeModulesDirectoriesReturnsEmptyResult(t *testing.T) {
	rootDirectory := t.TempDir()

	actualDirectories, err := findNodeModulesDirectories(rootDirectory)
	if err != nil {
		t.Fatalf("find node_modules directories: %v", err)
	}

	if len(actualDirectories) != 0 {
		t.Fatalf(
			"expected no node_modules directories, got %v",
			actualDirectories,
		)
	}
}

func TestFindNodeModulesDirectoriesRejectsFileAsRoot(t *testing.T) {
	rootDirectory := t.TempDir()
	filePath := filepath.Join(rootDirectory, "example.txt")

	if err := os.WriteFile(filePath, []byte("test"), 0o600); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	_, err := findNodeModulesDirectories(filePath)
	if err == nil {
		t.Fatal("expected an error when root path is a file")
	}
}
