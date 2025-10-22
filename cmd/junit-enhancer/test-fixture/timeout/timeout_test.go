package timeout

import (
	"testing"
	"time"
)

func TestCurrentTime_ShouldTimeout(t *testing.T) {
	t.Parallel()
	time.Sleep(15 * time.Second)
}

func TestCurrentTime_ShouldNotTimeout(t *testing.T) {
	t.Parallel()
	ct := CurrentTime()
	if ct != 0 {
		t.Fatalf("Expected CurrentTime to return zero.")
	}
}
