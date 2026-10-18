package main

import (
	"fmt"
	"strconv"
	"strings"
)

type selectedCandidate struct {
	Index     int
	Candidate candidateResult
}

func parseCandidateIndexes(value string) ([]int, error) {
	if value == "" {
		return nil, fmt.Errorf(
			"selection must not be empty",
		)
	}

	components := strings.Split(value, ",")

	indexes := make([]int, 0, len(components))
	seenIndexes := make(map[int]struct{}, len(components))

	for _, component := range components {
		if !isPositiveDecimalIndex(component) {
			return nil, fmt.Errorf(
				"invalid candidate index %q: indexes must be positive decimal integers",
				component,
			)
		}

		index, err := strconv.Atoi(component)
		if err != nil {
			return nil, fmt.Errorf(
				"parse candidate index %q: %w",
				component,
				err,
			)
		}

		if index <= 0 {
			return nil, fmt.Errorf(
				"candidate index must be greater than zero: %d",
				index,
			)
		}

		if _, exists := seenIndexes[index]; exists {
			return nil, fmt.Errorf(
				"candidate index %d was selected more than once",
				index,
			)
		}

		seenIndexes[index] = struct{}{}
		indexes = append(indexes, index)
	}

	return indexes, nil
}

func isPositiveDecimalIndex(value string) bool {
	if value == "" {
		return false
	}

	for index := 0; index < len(value); index++ {
		if value[index] < '0' || value[index] > '9' {
			return false
		}
	}

	return true
}

func selectCandidates(
	candidates []candidateResult,
	indexes []int,
) ([]selectedCandidate, error) {
	selectedCandidates := make(
		[]selectedCandidate,
		0,
		len(indexes),
	)

	for _, index := range indexes {
		if index < 1 || index > len(candidates) {
			return nil, fmt.Errorf(
				"candidate index %d does not exist",
				index,
			)
		}

		candidate := candidates[index-1]

		if candidate.ValidationError != nil {
			return nil, fmt.Errorf(
				"candidate index %d is not eligible for cleanup",
				index,
			)
		}

		selectedCandidates = append(
			selectedCandidates,
			selectedCandidate{
				Index:     index,
				Candidate: candidate,
			},
		)
	}

	return selectedCandidates, nil
}
