package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateCleanupTarget(
	rootDirectory string,
	targetDirectory string,
) error {
	absoluteRoot, err := filepath.Abs(rootDirectory)
	if err != nil {
		return fmt.Errorf("resolve root directory: %w", err)
	}

	absoluteTarget, err := filepath.Abs(targetDirectory)
	if err != nil {
		return fmt.Errorf("resolve cleanup target: %w", err)
	}

	if filepath.Base(absoluteTarget) != nodeModulesDirectoryName {
		return fmt.Errorf(
			"cleanup target is not named %q: %s",
			nodeModulesDirectoryName,
			absoluteTarget,
		)
	}

	if err := validatePathContainedWithin(
		absoluteRoot,
		absoluteTarget,
	); err != nil {
		return fmt.Errorf(
			"validate cleanup target location: %w",
			err,
		)
	}

	if err := validateOrdinaryDirectoryPath(
		absoluteRoot,
		absoluteTarget,
	); err != nil {
		return fmt.Errorf(
			"validate cleanup target path: %w",
			err,
		)
	}

	targetInfo, err := os.Stat(absoluteTarget)
	if err != nil {
		return fmt.Errorf(
			"inspect cleanup target: %w",
			err,
		)
	}

	if !targetInfo.IsDir() {
		return fmt.Errorf(
			"cleanup target is not a directory: %s",
			absoluteTarget,
		)
	}

	resolvedRoot, err := filepath.EvalSymlinks(absoluteRoot)
	if err != nil {
		return fmt.Errorf(
			"resolve root directory links: %w",
			err,
		)
	}

	resolvedTarget, err := filepath.EvalSymlinks(absoluteTarget)
	if err != nil {
		return fmt.Errorf(
			"resolve cleanup target links: %w",
			err,
		)
	}

	if err := validatePathContainedWithin(
		resolvedRoot,
		resolvedTarget,
	); err != nil {
		return fmt.Errorf(
			"validate resolved cleanup target location: %w",
			err,
		)
	}

	return nil
}

func validatePathContainedWithin(
	rootDirectory string,
	targetDirectory string,
) error {
	relativePath, err := filepath.Rel(
		rootDirectory,
		targetDirectory,
	)
	if err != nil {
		return fmt.Errorf(
			"determine relative path: %w",
			err,
		)
	}

	if relativePath == "." {
		return fmt.Errorf(
			"cleanup target must not be the scan root",
		)
	}

	if !filepath.IsLocal(relativePath) {
		return fmt.Errorf(
			"cleanup target is outside the scan root: %s",
			targetDirectory,
		)
	}

	return nil
}

func validateOrdinaryDirectoryPath(
	rootDirectory string,
	targetDirectory string,
) error {
	relativePath, err := filepath.Rel(
		rootDirectory,
		targetDirectory,
	)
	if err != nil {
		return fmt.Errorf(
			"determine relative cleanup path: %w",
			err,
		)
	}

	if relativePath == "." {
		return fmt.Errorf(
			"cleanup target must not be the scan root",
		)
	}

	if !filepath.IsLocal(relativePath) {
		return fmt.Errorf(
			"cleanup target is outside the scan root: %s",
			targetDirectory,
		)
	}

	currentPath := rootDirectory

	for _, component := range strings.Split(
		relativePath,
		string(filepath.Separator),
	) {
		currentPath = filepath.Join(
			currentPath,
			component,
		)

		entryInfo, err := os.Lstat(currentPath)
		if err != nil {
			return fmt.Errorf(
				"inspect path entry %q: %w",
				currentPath,
				err,
			)
		}

		if !entryInfo.IsDir() {
			return fmt.Errorf(
				"path entry is not an ordinary directory: %s",
				currentPath,
			)
		}
	}

	return nil
}
