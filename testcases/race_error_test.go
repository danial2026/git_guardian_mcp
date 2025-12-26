package testcases

import (
	"testing"
	"time"
)

// TestRaceCondition has a REAL race condition
func TestRaceCondition(t *testing.T) {
	counter := 0

	// Multiple goroutines writing to same variable - RACE!
	for i := 0; i < 100; i++ {
		go func() {
			counter++ // ERROR: concurrent write without synchronization
		}()
	}

	// Also reading while writing - RACE!
	go func() {
		for i := 0; i < 100; i++ {
			_ = counter // ERROR: concurrent read without synchronization
			time.Sleep(time.Microsecond)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	t.Logf("Counter value: %d (unpredictable due to race)", counter)
}

// TestMapRace has unprotected map access - RACE!
func TestMapRace(t *testing.T) {
	data := make(map[string]int)

	// Goroutine 1: writes to map
	go func() {
		for i := 0; i < 1000; i++ {
			data["key"] = i // ERROR: concurrent map write
		}
	}()

	// Goroutine 2: also writes to map
	go func() {
		for i := 0; i < 1000; i++ {
			data["key"] = i * 2 // ERROR: concurrent map write
		}
	}()

	// Goroutine 3: reads from map
	go func() {
		for i := 0; i < 1000; i++ {
			_ = data["key"] // ERROR: concurrent map read
		}
	}()

	time.Sleep(100 * time.Millisecond)
	t.Log("Map race completed")
}

// TestSliceRace has slice append races
func TestSliceRace(t *testing.T) {
	slice := []int{}

	// Multiple goroutines appending - RACE!
	for i := 0; i < 50; i++ {
		go func(val int) {
			slice = append(slice, val) // ERROR: concurrent slice modification
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	t.Logf("Slice length: %d (unpredictable)", len(slice))
}

// TestSharedStructRace modifies shared struct - RACE!
func TestSharedStructRace(t *testing.T) {
	type Stats struct {
		Count int
		Total int
	}

	stats := &Stats{}

	// Multiple goroutines modifying - RACE!
	for i := 0; i < 100; i++ {
		go func(val int) {
			stats.Count++      // ERROR: concurrent write
			stats.Total += val // ERROR: concurrent write
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	t.Logf("Stats: Count=%d, Total=%d", stats.Count, stats.Total)
}
