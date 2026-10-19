package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanupCandidateRemovesValidTarget(
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

	result := cleanupCandidate(
		rootDirectory,
		selectedCandidate{
			Index: 1,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: targetDirectory,
				},
			},
		},
	)

	if result.Err != nil {
		t.Fatalf(
			"cleanup valid candidate: %v",
			result.Err,
		)
	}

	if result.Index != 1 {
		t.Fatalf(
			"expected candidate index 1, got %d",
			result.Index,
		)
	}

	if result.Path != targetDirectory {
		t.Fatalf(
			"expected result path %q, got %q",
			targetDirectory,
			result.Path,
		)
	}

	if _, err := os.Lstat(targetDirectory); !os.IsNotExist(err) {
		t.Fatalf(
			"expected cleanup target to be absent, got %v",
			err,
		)
	}
}

func TestCleanupCandidateRejectsMissingTargetBeforeDeletion(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	targetDirectory := filepath.Join(
		rootDirectory,
		"project",
		"node_modules",
	)

	result := cleanupCandidate(
		rootDirectory,
		selectedCandidate{
			Index: 2,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: targetDirectory,
				},
			},
		},
	)

	if result.Err == nil {
		t.Fatal(
			"expected missing cleanup target to fail final validation",
		)
	}

	if !strings.Contains(
		result.Err.Error(),
		"final cleanup validation failed",
	) {
		t.Fatalf(
			"expected final validation failure, got %v",
			result.Err,
		)
	}
}

func TestCleanupCandidateRejectsWrongDirectoryName(
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

	result := cleanupCandidate(
		rootDirectory,
		selectedCandidate{
			Index: 3,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: targetDirectory,
				},
			},
		},
	)

	if result.Err == nil {
		t.Fatal(
			"expected wrong-name cleanup target to fail final validation",
		)
	}

	if _, err := os.Stat(targetDirectory); err != nil {
		t.Fatalf(
			"expected rejected target to remain present: %v",
			err,
		)
	}
}

func TestVerifyCleanupTargetRemovedAcceptsMissingPath(
	t *testing.T,
) {
	targetDirectory := filepath.Join(
		t.TempDir(),
		"node_modules",
	)

	if err := verifyCleanupTargetRemoved(
		targetDirectory,
	); err != nil {
		t.Fatalf(
			"verify missing cleanup target: %v",
			err,
		)
	}
}

func TestVerifyCleanupTargetRemovedRejectsExistingPath(
	t *testing.T,
) {
	targetDirectory := filepath.Join(
		t.TempDir(),
		"node_modules",
	)

	if err := os.MkdirAll(targetDirectory, 0o755); err != nil {
		t.Fatalf(
			"create cleanup target: %v",
			err,
		)
	}

	err := verifyCleanupTargetRemoved(targetDirectory)
	if err == nil {
		t.Fatal(
			"expected existing cleanup target to fail removal verification",
		)
	}
}

func TestCleanupSelectedCandidatesProcessesEveryCandidate(
	t *testing.T,
) {
	rootDirectory := t.TempDir()

	firstTarget := filepath.Join(
		rootDirectory,
		"first",
		"node_modules",
	)

	missingTarget := filepath.Join(
		rootDirectory,
		"missing",
		"node_modules",
	)

	thirdTarget := filepath.Join(
		rootDirectory,
		"third",
		"node_modules",
	)

	for _, targetDirectory := range []string{
		firstTarget,
		thirdTarget,
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
					Path: missingTarget,
				},
			},
		},
		{
			Index: 3,
			Candidate: candidateResult{
				ScanResult: scanResult{
					Path: thirdTarget,
				},
			},
		},
	}

	results := cleanupSelectedCandidates(
		rootDirectory,
		selectedCandidates,
	)

	if len(results) != 3 {
		t.Fatalf(
			"expected 3 cleanup results, got %d",
			len(results),
		)
	}

	if results[0].Err != nil {
		t.Fatalf(
			"expected first cleanup to succeed: %v",
			results[0].Err,
		)
	}

	if results[1].Err == nil {
		t.Fatal(
			"expected missing second candidate to fail",
		)
	}

	if results[2].Err != nil {
		t.Fatalf(
			"expected third cleanup to succeed after second failure: %v",
			results[2].Err,
		)
	}

	if cleanupFailureCount(results) != 1 {
		t.Fatalf(
			"expected 1 cleanup failure, got %d",
			cleanupFailureCount(results),
		)
	}

	for _, targetDirectory := range []string{
		firstTarget,
		thirdTarget,
	} {
		if _, err := os.Lstat(targetDirectory); !os.IsNotExist(err) {
			t.Fatalf(
				"expected target %q to be absent, got %v",
				targetDirectory,
				err,
			)
		}
	}
}
