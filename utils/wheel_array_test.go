package utils

import (
	"slices"
	"testing"
)

func TestWheelArray_Add(t *testing.T) {
	// Initialize a new WheelArray
	wa := NewWheelArray[int](5)

	// Add some elements to the WheelArray
	for i := 0; i < wa.size; i++ {
		wa.Add(i)
	}

	// Check if the elements were added correctly
	for i, v := range wa.data {
		if v != i {
			t.Errorf("Expected %d, but got %d", i, v)
		}
	}

	// Test the rotation of the WheelArray
	wa.Add(5)
	if wa.data[0] != 5 {
		t.Errorf("Expected %d, but got %d", 5, wa.data[0])
	}

	if slices.Compare(wa.data, []int{5, 1, 2, 3, 4}) != 0 {
		t.Errorf("Expected %v, but got %v", []int{5, 1, 2, 3, 4}, wa.data)
	}

	if slices.Compare(wa.Get(), []int{1, 2, 3, 4, 5}) != 0 {
		t.Errorf("Expected %v, but got %v", []int{1, 2, 3, 4, 5}, wa.Get())
	}
}
