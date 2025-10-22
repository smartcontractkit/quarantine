package broken

import (
	"testing"
)

func TestBuildFailure(t *testing.T) {
	// pass extra param to purposefully break the function call
	result, err := Break(1, 2, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5, got %f", result)
	}
}
