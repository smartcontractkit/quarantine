package timeout

import (
	"testing"
	"time"
)

func TestCurrentTime_ShouldTimeout(t *testing.T) {
	t.Parallel()

	deadline, ok := t.Deadline()
	if !ok {
		t.Fatalf("Deadline is not set, needed for test to timeout properly")
	}

	time.Sleep(time.Until(deadline) + 1*time.Second)
}

func TestCurrentTime_ShouldNotTimeout(t *testing.T) {
	t.Parallel()
	ct := CurrentTime()
	if ct != 0 {
		t.Fatalf("Expected CurrentTime to return zero.")
	}
}
