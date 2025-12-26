package testcases

import (
	"fmt"
	"testing"
)

// TestGoVetError has REAL go vet issues
func TestGoVetError(t *testing.T) {
	// Wrong Printf format - expects %d but gets string
	name := "John"
	fmt.Printf("Age: %d\n", name) // ERROR: format %d expects int, got string

	// Unreachable code after return
	if true {
		return
		t.Log("This line is unreachable") // ERROR: unreachable code
	}

	// Append result not used
	slice := []int{1, 2, 3}
	append(slice, 4) // ERROR: result of append not used

	// Comparing function to nil (always true)
	if TestGoVetError != nil {
		t.Log("Function pointer is always non-nil") // ERROR: comparison always true
	}
}

// TestMoreVetIssues has additional problems
func TestMoreVetIssues(t *testing.T) {
	// Printf with wrong number of args
	fmt.Printf("Name: %s, Age: %d\n", "Alice") // ERROR: missing argument for %d

	// Self-assignment
	x := 5
	x = x // ERROR: self-assignment

	_ = x
}
