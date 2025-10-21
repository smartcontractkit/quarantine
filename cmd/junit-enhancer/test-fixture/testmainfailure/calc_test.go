package testmainfailure

import (
	"fmt"
	"os"
	"testing"
)

// TestMain runs once for the entire package. We intentionally fail here
// (after running the tests) to act as a failing fixture for tooling.
func TestMain(m *testing.M) {
	// Example setup
	fmt.Fprintln(os.Stderr, "[fixture] global setup")

	// Run all tests in this package
	code := m.Run()

	// Example teardown
	fmt.Fprintln(os.Stderr, "[fixture] global teardown")

	// Intentionally fail the whole suite, regardless of test results.
	// (If you want to preserve the real exit code sometimes, swap this
	// with: os.Exit(code) guarded by an env flag.)
	fmt.Fprintln(os.Stderr, "[fixture] failing on purpose from TestMain")
	_ = code // keep linter happy if you don't use 'code'
	os.Exit(1)
}

func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Fatalf("Add(2,3) = %d; want 5", got)
	}
}
