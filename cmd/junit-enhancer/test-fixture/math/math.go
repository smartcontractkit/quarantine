package math

import "errors"

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func Power(base, exp int) int {
	if exp == 0 {
		return 1
	}
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}
