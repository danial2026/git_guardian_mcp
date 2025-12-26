package testcases

import (
	"fmt"
	"os"
	"testing"
)

// TestUnusedVariables has REAL unused variables
func TestUnusedVariables(t *testing.T) {
	// These variables are declared but never used
	unusedVar := 42                // ERROR: declared and not used
	anotherUnused := "hello world" // ERROR: declared and not used
	yetAnother := 3.14159          // ERROR: declared and not used

	// Only using one variable
	used := "I am used"
	fmt.Println(used)

	// Function result ignored
	calculateSomething() // Return value ignored

	// Variable declared but never read
	var neverRead int
	neverRead = 100 // ERROR: value assigned but never used

	// Import os but barely use it
	_ = os.Stdout // Suspicious use just to avoid error
}

// calculateSomething returns a value that's ignored
func calculateSomething() int {
	return 42
}

// TestUnusedParameters has unused function parameters
func TestUnusedParameters(t *testing.T) {
	result := addButIgnoreSecond(10, 20)
	if result != 10 {
		t.Error("Wrong result")
	}
}

// addButIgnoreSecond has unused parameter
func addButIgnoreSecond(a, b int) int {
	// ERROR: parameter 'b' is never used
	return a
}

// NeverCalledFunction is completely unused
func NeverCalledFunction() { // ERROR: function is never called
	fmt.Println("I am never called")
}
