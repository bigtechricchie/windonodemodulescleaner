package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRevalidateSelectedCandidatesAcceptsValidCandidate(
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

	selectedCandidates := []selectedCandidate{
		{
			Index: 1,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: targetDirectory,
				},
			},
		},
	}

	if err := revalidateSelectedCandidates(
		rootDirectory,
		selectedCandidates,
	); err != nil {
		t.Fatalf(
			"revalidate selected candidate: %v",
			err,
		)
	}
}

func TestRevalidateSelectedCandidatesAcceptsMultipleValidCandidates(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	firstTarget := filepath.Join(
		rootDirectory,
		"first",
		"node_modules",
	)

	secondTarget := filepath.Join(
		rootDirectory,
		"second",
		"node_modules",
	)

	for _, targetDirectory := range []string{
		firstTarget,
		secondTarget,
	} {
		if err := os.MkdirAll(targetDirectory, 0o755); err != nil {
			t.Fatalf(
				"create cleanup target %q: %v",
				targetDirectory,
				err,
			)
		}
	}

	selectedCandidates := []selectedCandidate{
		{
			Index: 1,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: firstTarget,
				},
			},
		},
		{
			Index: 2,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: secondTarget,
				},
			},
		},
	}

	if err := revalidateSelectedCandidates(
		rootDirectory,
		selectedCandidates,
	); err != nil {
		t.Fatalf(
			"revalidate selected candidates: %v",
			err,
		)
	}
}

func TestRevalidateSelectedCandidatesRejectsEmptySelection(
	t *testing.T,
) {
	err := revalidateSelectedCandidates(
		t.TempDir(),
		nil,
	)
	if err == nil {
		t.Fatal(
			"expected empty selected-candidate set to be rejected",
		)
	}
}

func TestRevalidateSelectedCandidatesRejectsDisappearedCandidate(
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
			"expected initial cleanup target to be valid: %v",
			err,
		)
	}

	selectedCandidates := []selectedCandidate{
		{
			Index: 4,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: targetDirectory,
				},
			},
		},
	}

	if err := os.RemoveAll(targetDirectory); err != nil {
		t.Fatalf(
			"remove cleanup target before fresh validation: %v",
			err,
		)
	}

	err := revalidateSelectedCandidates(
		rootDirectory,
		selectedCandidates,
	)
	if err == nil {
		t.Fatal(
			"expected disappeared cleanup target to fail fresh validation",
		)
	}

	if !strings.Contains(err.Error(), "candidate index 4") {
		t.Fatalf(
			"expected error to identify candidate index 4, got %v",
			err,
		)
	}
}

func TestRevalidateSelectedCandidatesRejectsChangedPathTopology(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	workspaceDirectory := filepath.Join(
		rootDirectory,
		"workspace",
	)

	targetDirectory := filepath.Join(
		workspaceDirectory,
		"app",
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
			"expected initial cleanup target to be valid: %v",
			err,
		)
	}

	selectedCandidates := []selectedCandidate{
		{
			Index: 2,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: targetDirectory,
				},
			},
		},
	}

	actualWorkspace := filepath.Join(
		rootDirectory,
		"actual-workspace",
	)

	if err := os.Rename(
		workspaceDirectory,
		actualWorkspace,
	); err != nil {
		t.Fatalf(
			"move workspace before fresh validation: %v",
			err,
		)
	}

	if err := os.Symlink(
		actualWorkspace,
		workspaceDirectory,
	); err != nil {
		if runtime.GOOS == "windows" {
			t.Skipf(
				"symlink creation unavailable on this Windows host: %v",
				err,
			)
		}

		t.Fatalf(
			"replace workspace with symlink: %v",
			err,
		)
	}

	err := revalidateSelectedCandidates(
		rootDirectory,
		selectedCandidates,
	)
	if err == nil {
		t.Fatal(
			"expected changed cleanup-path topology to fail fresh validation",
		)
	}

	if !strings.Contains(err.Error(), "candidate index 2") {
		t.Fatalf(
			"expected error to identify candidate index 2, got %v",
			err,
		)
	}
}
