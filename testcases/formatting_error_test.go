package testcases

import (
	"testing"
)

// TestFormattingError has REAL formatting issues
func TestFormattingError(t *testing.T) {
	input := 10
	expected := 20

	result := input * 2

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// BadlyFormatted demonstrates terrible formatting
func BadlyFormatted(a int, b int) (int, error) {
	if a > b {
		return a, nil
	}
	return b, nil
}

// MoreBadFormatting shows inconsistent spacing
func MoreBadFormatting() {
	x := 1 + 2 + 3
	y := 4 + 5
	_ = x
	_ = y
}
