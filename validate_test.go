package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestValidateCleanupTargetAcceptsNodeModulesWithinRoot(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	targetDirectory := filepath.Join(
		rootDirectory,
		"project",
		"node_modules",
	)

	if err := os.MkdirAll(targetDirectory, 0o755); err != nil {
		t.Fatalf(
			"create cleanup target: %v",
			err,
		)
	}

	if err := validateCleanupTarget(
		rootDirectory,
		targetDirectory,
	); err != nil {
		t.Fatalf(
			"expected cleanup target to be valid: %v",
			err,
		)
	}
}

func TestValidateCleanupTargetRejectsWrongDirectoryName(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	targetDirectory := filepath.Join(
		rootDirectory,
		"project",
		"dependencies",
	)

	if err := os.MkdirAll(targetDirectory, 0o755); err != nil {
		t.Fatalf(
			"create test directory: %v",
			err,
		)
	}

	err := validateCleanupTarget(
		rootDirectory,
		targetDirectory,
	)
	if err == nil {
		t.Fatal(
			"expected directory with wrong name to be rejected",
		)
	}
}

func TestValidateCleanupTargetRejectsMissingTarget(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	targetDirectory := filepath.Join(
		rootDirectory,
		"project",
		"node_modules",
	)

	err := validateCleanupTarget(
		rootDirectory,
		targetDirectory,
	)
	if err == nil {
		t.Fatal(
			"expected missing cleanup target to be rejected",
		)
	}
}

func TestValidateCleanupTargetRejectsFile(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	projectDirectory := filepath.Join(
		rootDirectory,
		"project",
	)

	if err := os.MkdirAll(projectDirectory, 0o755); err != nil {
		t.Fatalf(
			"create project directory: %v",
			err,
		)
	}

	targetPath := filepath.Join(
		projectDirectory,
		"node_modules",
	)

	if err := os.WriteFile(
		targetPath,
		[]byte("not a directory"),
		0o600,
	); err != nil {
		t.Fatalf(
			"create test file: %v",
			err,
		)
	}

	err := validateCleanupTarget(
		rootDirectory,
		targetPath,
	)
	if err == nil {
		t.Fatal(
			"expected file cleanup target to be rejected",
		)
	}
}

func TestValidateCleanupTargetRejectsTargetOutsideRoot(
	t *testing.T,
) {
	parentDirectory := t.TempDir()

	rootDirectory := filepath.Join(
		parentDirectory,
		"projects",
	)

	outsideTarget := filepath.Join(
		parentDirectory,
		"outside",
		"node_modules",
	)

	if err := os.MkdirAll(rootDirectory, 0o755); err != nil {
		t.Fatalf(
			"create root directory: %v",
			err,
		)
	}

	if err := os.MkdirAll(outsideTarget, 0o755); err != nil {
		t.Fatalf(
			"create outside target: %v",
			err,
		)
	}

	err := validateCleanupTarget(
		rootDirectory,
		outsideTarget,
	)
	if err == nil {
		t.Fatal(
			"expected target outside root to be rejected",
		)
	}
}

func TestValidateCleanupTargetRejectsPrefixConfusion(
	t *testing.T,
) {
	parentDirectory := t.TempDir()

	rootDirectory := filepath.Join(
		parentDirectory,
		"projects",
	)

	targetDirectory := filepath.Join(
		parentDirectory,
		"projects-other",
		"node_modules",
	)

	if err := os.MkdirAll(rootDirectory, 0o755); err != nil {
		t.Fatalf(
			"create root directory: %v",
			err,
		)
	}

	if err := os.MkdirAll(targetDirectory, 0o755); err != nil {
		t.Fatalf(
			"create cleanup target: %v",
			err,
		)
	}

	err := validateCleanupTarget(
		rootDirectory,
		targetDirectory,
	)
	if err == nil {
		t.Fatal(
			"expected prefix-confusion target outside root to be rejected",
		)
	}
}

func TestValidateCleanupTargetRejectsSymlinkTarget(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	projectDirectory := filepath.Join(
		rootDirectory,
		"project",
	)

	actualDirectory := filepath.Join(
		rootDirectory,
		"shared-dependencies",
	)

	targetDirectory := filepath.Join(
		projectDirectory,
		"node_modules",
	)

	if err := os.MkdirAll(projectDirectory, 0o755); err != nil {
		t.Fatalf(
			"create project directory: %v",
			err,
		)
	}

	if err := os.MkdirAll(actualDirectory, 0o755); err != nil {
		t.Fatalf(
			"create symlink destination: %v",
			err,
		)
	}

	if err := os.Symlink(
		actualDirectory,
		targetDirectory,
	); err != nil {
		if runtime.GOOS == "windows" {
			t.Skipf(
				"symlink creation unavailable on this Windows host: %v",
				err,
			)
		}

		t.Fatalf(
			"create symlink cleanup target: %v",
			err,
		)
	}

	err := validateCleanupTarget(
		rootDirectory,
		targetDirectory,
	)
	if err == nil {
		t.Fatal(
			"expected symlink cleanup target to be rejected",
		)
	}
}

func TestValidateCleanupTargetRejectsSymlinkAncestor(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	actualWorkspace := filepath.Join(
		rootDirectory,
		"actual-workspace",
	)

	projectDirectory := filepath.Join(
		actualWorkspace,
		"app",
	)

	targetDirectory := filepath.Join(
		projectDirectory,
		"node_modules",
	)

	linkedWorkspace := filepath.Join(
		rootDirectory,
		"workspace",
	)

	if err := os.MkdirAll(targetDirectory, 0o755); err != nil {
		t.Fatalf(
			"create cleanup target: %v",
			err,
		)
	}

	if err := os.Symlink(
		actualWorkspace,
		linkedWorkspace,
	); err != nil {
		if runtime.GOOS == "windows" {
			t.Skipf(
				"symlink creation unavailable on this Windows host: %v",
				err,
			)
		}

		t.Fatalf(
			"create symlink ancestor: %v",
			err,
		)
	}

	linkedTarget := filepath.Join(
		linkedWorkspace,
		"app",
		"node_modules",
	)

	err := validateCleanupTarget(
		rootDirectory,
		linkedTarget,
	)
	if err == nil {
		t.Fatal(
			"expected cleanup target beneath symlink ancestor to be rejected",
		)
	}
}

func TestValidatePathContainedWithinRejectsRootItself(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	err := validatePathContainedWithin(
		rootDirectory,
		rootDirectory,
	)
	if err == nil {
		t.Fatal(
			"expected scan root itself to be rejected",
		)
	}
}
