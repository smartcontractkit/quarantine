package quarantine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/quarantine"
)

func TestFlaky(t *testing.T) {
	t.Run("skip flaky tests", func(t *testing.T) {
		t.Setenv(quarantine.RunQuarantinedTestsEnvVar, "false")
		quarantine.Flaky(t, "TEST-123")

		t.Cleanup(func() {
			require.True(t, t.Skipped(), "quarantined test should be skipped when RUN_FLAKY_TESTS is false")
		})
	})

	t.Run("run flaky tests", func(t *testing.T) {
		t.Setenv(quarantine.RunQuarantinedTestsEnvVar, "true")
		quarantine.Flaky(t, "TEST-123")

		t.Cleanup(func() {
			require.False(t, t.Skipped(), "quarantined test should not be skipped when RUN_FLAKY_TESTS is true")
			t.Log("This test is intentionally skipped! Skipping = Passing!")
		})
	})
}

func TestTimeout(t *testing.T) {
	t.Run("skip timeout tests", func(t *testing.T) {
		t.Setenv(quarantine.RunTimeoutTestsEnvVar, "false")
		quarantine.Timeout(t, "TEST-123")

		t.Cleanup(func() {
			require.True(t, t.Skipped(), "timeout test should be skipped when RUN_TIMEOUT_TESTS is false")
		})
	})

	t.Run("run timeout tests", func(t *testing.T) {
		t.Setenv(quarantine.RunTimeoutTestsEnvVar, "true")
		quarantine.Timeout(t, "TEST-123")

		t.Cleanup(func() {
			require.False(t, t.Skipped(), "timeout test should not be skipped when RUN_TIMEOUT_TESTS is true")
			t.Log("This test is intentionally skipped! Skipping = Passing!")
		})
	})
}
