package testcases

import (
	"errors"
	"io/ioutil"
	"testing"
)

// TestLintingErrors demonstrates golangci-lint errors
func TestLintingErrors(t *testing.T) {
	// Using deprecated ioutil instead of os
	data, err := ioutil.ReadFile("nonexistent.txt")
	if err != nil {
		t.Log("Expected error")
	}

	// Naked return in function with named returns
	result := processDataBadly(data)
	if result == nil {
		t.Log("Data is nil")
	}

	// Magic number without constant
	if len(data) > 1024 {
		t.Log("Data too large")
	}

	// Error shadowing
	err = errors.New("new error")
	if err != nil {
		err := errors.New("shadowed error")
		t.Log(err)
	}
}

// Naked return - bad practice
func processDataBadly(data []byte) (result []byte) {
	if len(data) == 0 {
		return // naked return
	}
	result = data
	return // another naked return
}

// Function too complex (cyclomatic complexity)
func ComplexFunction(a, b, c, d, e int) int {
	if a > 0 {
		if b > 0 {
			if c > 0 {
				if d > 0 {
					if e > 0 {
						return a + b + c + d + e
					} else {
						return a + b + c + d
					}
				} else {
					return a + b + c
				}
			} else {
				return a + b
			}
		} else {
			return a
		}
	} else {
		return 0
	}
}

// TestShouldFailWithComplexity tests the complex function
func TestShouldFailWithComplexity(t *testing.T) {
	result := ComplexFunction(1, 2, 3, 4, 5)
	if result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}
}
