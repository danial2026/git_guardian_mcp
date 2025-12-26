package testcases

import (
	"testing"
)

// TestPanicError demonstrates code that will panic
func TestPanicError(t *testing.T) {
	// This will panic with nil pointer dereference
	var ptr *string
	t.Log("Length:", len(*ptr)) // PANIC: nil pointer dereference

	// This will also panic
	slice := []int{1, 2, 3}
	t.Log("Value:", slice[10]) // PANIC: index out of range
}

// TestDivideByZero demonstrates division by zero panic
func TestDivideByZero(t *testing.T) {
	numerator := 100
	denominator := 0
	result := numerator / denominator // PANIC: division by zero
	t.Log("Result:", result)
}

// TestNilMapWrite demonstrates writing to nil map
func TestNilMapWrite(t *testing.T) {
	var m map[string]int
	m["key"] = 42 // PANIC: assignment to entry in nil map
	t.Log("Map:", m)
}

// TestNilSliceIndex demonstrates nil slice access
func TestNilSliceIndex(t *testing.T) {
	var slice []int
	value := slice[0] // PANIC: index out of range
	t.Log("Value:", value)
}

// TestTypeAssertionPanic demonstrates bad type assertion
func TestTypeAssertionPanic(t *testing.T) {
	var i interface{} = "hello"
	num := i.(int) // PANIC: interface conversion
	t.Log("Number:", num)
}

// TestClosedChannelWrite demonstrates writing to closed channel
func TestClosedChannelWrite(t *testing.T) {
	ch := make(chan int)
	close(ch)
	ch <- 42 // PANIC: send on closed channel
	t.Log("Sent value")
}
