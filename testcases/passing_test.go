package testcases

import (
	"testing"
)

// TestPassing is a well-written test that passes all checks
func TestPassing(t *testing.T) {
	// Arrange
	input := 42
	expected := 42

	// Act
	result := ProcessNumber(input)

	// Assert
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// ProcessNumber is a simple helper function
func ProcessNumber(n int) int {
	return n
}

// TestAnotherPassing demonstrates proper test structure
func TestAnotherPassing(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive", 5, 5},
		{"zero", 0, 0},
		{"negative", -3, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessNumber(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessNumber(%d) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}
