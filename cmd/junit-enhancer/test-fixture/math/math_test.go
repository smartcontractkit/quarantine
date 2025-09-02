package math

import (
	"testing"
)

func TestDivide(t *testing.T) {
	result, err := Divide(10, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5, got %f", result)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(10, 0)
	if err == nil {
		t.Error("Expected error for division by zero")
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		base     int
		exp      int
		expected int
	}{
		{2, 3, 8},
		{5, 0, 1},
		{3, 2, 9},
		{1, 10, 1},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := Power(test.base, test.exp)
			if result != test.expected {
				t.Errorf("Power(%d, %d) = %d, expected %d", test.base, test.exp, result, test.expected)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	t.Run("base_cases", func(t *testing.T) {
		if Factorial(0) != 1 {
			t.Error("Factorial(0) should be 1")
		}
		if Factorial(1) != 1 {
			t.Error("Factorial(1) should be 1")
		}
	})
	
	t.Run("positive_numbers", func(t *testing.T) {
		tests := map[int]int{
			2: 2,
			3: 6,
			4: 24,
			5: 120,
		}
		
		for input, expected := range tests {
			result := Factorial(input)
			if result != expected {
				t.Errorf("Factorial(%d) = %d, expected %d", input, result, expected)
			}
		}
	})
}

func FuzzDivide(f *testing.F) {
	f.Add(10.0, 2.0)
	f.Add(100.0, 4.0)
	f.Add(-10.0, 2.0)
	
	f.Fuzz(func(t *testing.T, a, b float64) {
		if b == 0 {
			return // Skip division by zero
		}
		result, err := Divide(a, b)
		if err != nil {
			t.Errorf("Unexpected error for Divide(%f, %f): %v", a, b, err)
		}
		expected := a / b
		if result != expected {
			t.Errorf("Divide(%f, %f) = %f, expected %f", a, b, result, expected)
		}
	})
}
