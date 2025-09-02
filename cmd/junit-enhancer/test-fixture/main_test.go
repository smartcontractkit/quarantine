package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}

func TestMultiply(t *testing.T) {
	result := Multiply(4, 5)
	if result != 20 {
		t.Errorf("Expected 20, got %d", result)
	}
}

func TestAddSubtests(t *testing.T) {
	t.Run("positive_numbers", func(t *testing.T) {
		result := Add(1, 2)
		if result != 3 {
			t.Errorf("Expected 3, got %d", result)
		}
	})
	
	t.Run("negative_numbers", func(t *testing.T) {
		result := Add(-1, -2)
		if result != -3 {
			t.Errorf("Expected -3, got %d", result)
		}
	})
	
	t.Run("mixed_numbers", func(t *testing.T) {
		result := Add(5, -3)
		if result != 2 {
			t.Errorf("Expected 2, got %d", result)
		}
	})
}

func FuzzAdd(f *testing.F) {
	f.Add(1, 2)
	f.Add(-1, -2)
	f.Add(0, 0)
	
	f.Fuzz(func(t *testing.T, a, b int) {
		result := Add(a, b)
		expected := a + b
		if result != expected {
			t.Errorf("Add(%d, %d) = %d, expected %d", a, b, result, expected)
		}
	})
}

func FuzzMultiply(f *testing.F) {
	f.Add(2, 3)
	f.Add(-2, 3)
	f.Add(0, 5)
	
	f.Fuzz(func(t *testing.T, a, b int) {
		result := Multiply(a, b)
		expected := a * b
		if result != expected {
			t.Errorf("Multiply(%d, %d) = %d, expected %d", a, b, result, expected)
		}
	})
}
