package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const nodeModulesDirectoryName = "node_modules"

func findNodeModulesDirectories(rootDirectory string) ([]string, error) {
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

	var nodeModulesDirectories []string

	err = filepath.WalkDir(rootDirectory, func(path string, entry fs.DirEntry, walkErr error) error {
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

		nodeModulesDirectories = append(nodeModulesDirectories, path)

		return fs.SkipDir
	})
	if err != nil {
		return nil, fmt.Errorf("scan root directory: %w", err)
	}

	return nodeModulesDirectories, nil
}
