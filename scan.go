package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const nodeModulesDirectoryName = "node_modules"

type scanResult struct {
	Path      string
	SizeBytes int64
	SizeError error
}

func findNodeModulesDirectories(rootDirectory string) ([]scanResult, error) {
	rootDirectory, err := filepath.Abs(rootDirectory)
	if err != nil {
		return nil, fmt.Errorf("resolve root directory: %w", err)
	}

	rootInfo, err := os.Stat(rootDirectory)
	if err != nil {
		return nil, fmt.Errorf("inspect root directory: %w", err)
	}

	if !rootInfo.IsDir() {
		return nil, fmt.Errorf("root path is not a directory: %s", rootDirectory)
	}

	var results []scanResult

	err = filepath.WalkDir(
		rootDirectory,
		func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if !entry.IsDir() {
				return nil
			}

			if path == rootDirectory {
				return nil
			}

			if entry.Name() != nodeModulesDirectoryName {
				return nil
			}

			sizeBytes, sizeErr := calculateDirectorySize(path)

			results = append(results, scanResult{
				Path:      path,
				SizeBytes: sizeBytes,
				SizeError: sizeErr,
			})

			return fs.SkipDir
		},
	)
	if err != nil {
		return nil, fmt.Errorf("scan root directory: %w", err)
	}

	return results, nil
}

func calculateDirectorySize(rootDirectory string) (int64, error) {
	var sizeBytes int64

	err := filepath.WalkDir(
		rootDirectory,
		func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() {
				return nil
			}

			fileInfo, err := entry.Info()
			if err != nil {
				return fmt.Errorf("inspect %q: %w", path, err)
			}

			sizeBytes += fileInfo.Size()

			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf(
			"calculate size of %q: %w",
			rootDirectory,
			err,
		)
	}

	return sizeBytes, nil
}
