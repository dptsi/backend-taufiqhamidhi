package handler

import (
	"testing"
)

func TestHexToDecimal(t *testing.T) {
	tests := []struct {
		hexNumber   string
		expectedDec string
		expectedErr bool
	}{
		{"ABCDEF", "11259375", false},
		{"FF", "255", false},
		{"123", "", true},
		{"GHIJK", "", true},
	}

	for _, test := range tests {
		decNumber := hexToDecimal(test.hexNumber)

		if decNumber != test.expectedDec {
			if !test.expectedErr {
				t.Errorf("Expected decimal: %s, but got: %s", test.expectedDec, decNumber)
			}
		}
	}
}
