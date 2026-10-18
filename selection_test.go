package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseCandidateIndexesAcceptsSingleIndex(
	t *testing.T,
) {
	indexes, err := parseCandidateIndexes("1")
	if err != nil {
		t.Fatalf(
			"parse candidate indexes: %v",
			err,
		)
	}

	expected := []int{1}

	if !reflect.DeepEqual(indexes, expected) {
		t.Fatalf(
			"expected indexes %v, got %v",
			expected,
			indexes,
		)
	}
}

func TestParseCandidateIndexesAcceptsMultipleIndexes(
	t *testing.T,
) {
	indexes, err := parseCandidateIndexes("1,2,5")
	if err != nil {
		t.Fatalf(
			"parse candidate indexes: %v",
			err,
		)
	}

	expected := []int{1, 2, 5}

	if !reflect.DeepEqual(indexes, expected) {
		t.Fatalf(
			"expected indexes %v, got %v",
			expected,
			indexes,
		)
	}
}

func TestParseCandidateIndexesRejectsInvalidInput(
	t *testing.T,
) {
	testCases := []struct {
		name  string
		value string
	}{
		{
			name:  "empty",
			value: "",
		},
		{
			name:  "zero",
			value: "0",
		},
		{
			name:  "negative",
			value: "-1",
		},
		{
			name:  "letters",
			value: "abc",
		},
		{
			name:  "mixed",
			value: "1,x",
		},
		{
			name:  "range",
			value: "1-5",
		},
		{
			name:  "all",
			value: "all",
		},
		{
			name:  "asterisk",
			value: "*",
		},
		{
			name:  "whitespace",
			value: "1, 2",
		},
		{
			name:  "missing component",
			value: "1,,2",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				_, err := parseCandidateIndexes(
					testCase.value,
				)
				if err == nil {
					t.Fatalf(
						"expected %q to be rejected",
						testCase.value,
					)
				}
			},
		)
	}
}

func TestParseCandidateIndexesRejectsDuplicateIndexes(
	t *testing.T,
) {
	_, err := parseCandidateIndexes("1,1,2")
	if err == nil {
		t.Fatal(
			"expected duplicate candidate index to be rejected",
		)
	}
}

func TestSelectCandidatesAcceptsEligibleCandidates(
	t *testing.T,
) {
	candidates := []candidateResult{
		{
			ScanResult: scanResult{
				Path: "first",
			},
		},
		{
			ScanResult: scanResult{
				Path: "second",
			},
		},
		{
			ScanResult: scanResult{
				Path: "third",
			},
		},
	}

	selectedCandidates, err := selectCandidates(
		candidates,
		[]int{1, 3},
	)
	if err != nil {
		t.Fatalf(
			"select candidates: %v",
			err,
		)
	}

	if len(selectedCandidates) != 2 {
		t.Fatalf(
			"expected 2 selected candidates, got %d",
			len(selectedCandidates),
		)
	}

	if selectedCandidates[0].Index != 1 {
		t.Fatalf(
			"expected first selected index 1, got %d",
			selectedCandidates[0].Index,
		)
	}

	if selectedCandidates[0].Candidate.ScanResult.Path != "first" {
		t.Fatalf(
			"expected first selected candidate, got %q",
			selectedCandidates[0].Candidate.ScanResult.Path,
		)
	}

	if selectedCandidates[1].Index != 3 {
		t.Fatalf(
			"expected second selected index 3, got %d",
			selectedCandidates[1].Index,
		)
	}

	if selectedCandidates[1].Candidate.ScanResult.Path != "third" {
		t.Fatalf(
			"expected third candidate, got %q",
			selectedCandidates[1].Candidate.ScanResult.Path,
		)
	}
}

func TestSelectCandidatesRejectsMissingIndex(
	t *testing.T,
) {
	candidates := []candidateResult{
		{
			ScanResult: scanResult{
				Path: "first",
			},
		},
	}

	selectedCandidates, err := selectCandidates(
		candidates,
		[]int{1, 2},
	)
	if err == nil {
		t.Fatal(
			"expected nonexistent candidate index to be rejected",
		)
	}

	if selectedCandidates != nil {
		t.Fatalf(
			"expected no partial selection, got %v",
			selectedCandidates,
		)
	}
}

func TestSelectCandidatesRejectsIneligibleCandidate(
	t *testing.T,
) {
	candidates := []candidateResult{
		{
			ScanResult: scanResult{
				Path: "eligible",
			},
		},
		{
			ScanResult: scanResult{
				Path: "rejected",
			},
			ValidationError: errors.New(
				"candidate rejected",
			),
		},
	}

	selectedCandidates, err := selectCandidates(
		candidates,
		[]int{1, 2},
	)
	if err == nil {
		t.Fatal(
			"expected rejected candidate selection to fail",
		)
	}

	if selectedCandidates != nil {
		t.Fatalf(
			"expected no partial selection, got %v",
			selectedCandidates,
		)
	}
}

func TestSelectCandidatesRejectsIndexWhenNoCandidatesExist(
	t *testing.T,
) {
	selectedCandidates, err := selectCandidates(
		nil,
		[]int{1},
	)
	if err == nil {
		t.Fatal(
			"expected candidate index to be rejected when no candidates exist",
		)
	}

	if selectedCandidates != nil {
		t.Fatalf(
			"expected no selected candidates, got %v",
			selectedCandidates,
		)
	}
}
