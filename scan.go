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

type scanWarning struct {
	Path string
	Err  error
}

type scanOutcome struct {
	Results  []scanResult
	Warnings []scanWarning
}

func findNodeModulesDirectories(rootDirectory string) (scanOutcome, error) {
	rootDirectory, err := filepath.Abs(rootDirectory)
	if err != nil {
		return scanOutcome{}, fmt.Errorf("resolve root directory: %w", err)
	}

	rootInfo, err := os.Stat(rootDirectory)
	if err != nil {
		return scanOutcome{}, fmt.Errorf("inspect root directory: %w", err)
	}

	if !rootInfo.IsDir() {
		return scanOutcome{}, fmt.Errorf(
			"root path is not a directory: %s",
			rootDirectory,
		)
	}

	var outcome scanOutcome

	err = filepath.WalkDir(
		rootDirectory,
		func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				if path == rootDirectory {
					return walkErr
				}

				outcome.Warnings = append(
					outcome.Warnings,
					scanWarning{
						Path: path,
						Err:  walkErr,
					},
				)

				return nil
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

			outcome.Results = append(
				outcome.Results,
				scanResult{
					Path:      path,
					SizeBytes: sizeBytes,
					SizeError: sizeErr,
				},
			)

			return fs.SkipDir
		},
	)
	if err != nil {
		return scanOutcome{}, fmt.Errorf("scan root directory: %w", err)
	}

	return outcome, nil
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
