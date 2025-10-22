package compiles

import (
	"testing"
)

func TestNoBuildFailure(t *testing.T) {
	// pass extra param to purposefully break the function call
	result, err := Fix(1, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 0.5 {
		t.Errorf("Expected 0.5, got %f", result)
	}
}
