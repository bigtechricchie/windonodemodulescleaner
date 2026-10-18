package main

import "fmt"

func revalidateSelectedCandidates(
	rootDirectory string,
	selectedCandidates []selectedCandidate,
) error {
	if len(selectedCandidates) == 0 {
		return fmt.Errorf(
			"cannot revalidate an empty candidate selection",
		)
	}

	for _, selectedCandidate := range selectedCandidates {
		err := validateCleanupTarget(
			rootDirectory,
			selectedCandidate.Candidate.ScanResult.Path,
		)
		if err != nil {
			return fmt.Errorf(
				"candidate index %d failed fresh cleanup validation: %w",
				selectedCandidate.Index,
				err,
			)
		}
	}

	return nil
}
