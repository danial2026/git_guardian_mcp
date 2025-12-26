package testcases

import (
	"testing"
)

// TestFailingAssertion ACTUALLY FAILS
func TestFailingAssertion(t *testing.T) {
	input := 5
	expected := 100 // WRONG! Should be 15

	result := tripleNumber(input)

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result) // This WILL fail
	}
}

func tripleNumber(n int) int {
	return n * 3 // Returns 15, but test expects 100
}

// TestAlwaysFails is designed to always fail
func TestAlwaysFails(t *testing.T) {
	t.Fatal("This test intentionally fails to demonstrate test failure handling")
}

// TestBadMath has wrong logic
func TestBadMath(t *testing.T) {
	result := add(2, 2)
	expected := 5 // WRONG! 2 + 2 = 4, not 5

	if result != expected {
		t.Fatalf("Math is broken: %d + %d should be %d, got %d", 2, 2, expected, result)
	}
}

func add(a, b int) int {
	return a + b // Correct, but test expectation is wrong
}

// TestStringComparison fails comparison
func TestStringComparison(t *testing.T) {
	actual := "hello"
	expected := "goodbye" // WRONG!

	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual) // WILL fail
	}
}
