// Package quarantine provides a way to mark tests as flaky.
package quarantine

import (
	"os"
	"strings"
	"testing"
	"unicode"
)

const (
	// RunQuarantinedTestsEnvVar is the environment variable that controls whether to run quarantined tests.
	RunQuarantinedTestsEnvVar = "RUN_QUARANTINED_TESTS"
	// RunTimeoutTestsEnvVar is the environment variable that controls whether to run timeout tests.
	RunTimeoutTestsEnvVar = "RUN_TIMEOUT_TESTS"
)

// Flaky marks a test as flaky.
// To run tests marked as flaky, set the RUN_FLAKY_TESTS environment variable to true.
// To skip tests marked as flaky, set the RUN_FLAKY_TESTS environment variable to false (or don't set it at all).
//
// Example:
//
//	func TestFlaky(t *testing.T) {
//		quarantine.Flaky(t, "TEST-123")
//	}
func Flaky(tb testing.TB, ticket string) {
	tb.Helper()

	skipTest(tb, RunQuarantinedTestsEnvVar, "flaky", ticket)
}

// Timeout marks a test that is expected to timeout.
// To run tests marked as timeout, set the RUN_TIMEOUT_TESTS environment variable to true.
// To skip tests marked as timeout, set the RUN_TIMEOUT_TESTS environment variable to false (or don't set it at all).
//
// Example:
//
//	func TestTimeout(t *testing.T) {
//		quarantine.Timeout(t, "TEST-123")
//	}
func Timeout(tb testing.TB, ticket string) {
	tb.Helper()

	skipTest(tb, RunTimeoutTestsEnvVar, "timeout", ticket)
}

func skipTest(tb testing.TB, envVar, classification, ticket string) {
	tb.Helper()

	attr(tb, classification, ticket)
	classifiedStr := "Classified by branch-out (https://github.com/smartcontractkit/branch-out)"
	if os.Getenv(envVar) != "true" {
		tb.Skipf(
			"To run '%s' tests, set %s='true'.\n%s",
			classification,
			envVar,
			classifiedStr,
		)
	} else {
		tb.Logf("Running test marked as '%s'.", classification)
		tb.Cleanup(func() {
			tb.Logf(
				"Test is marked as %s, but still ran. To skip %s tests, set %s='false'.\n%s",
				classification, classification, envVar, classifiedStr,
			)
		})
	}
}

// attr replicates the functionality of testing.TB.Attr() for compatibility with older Go versions.
// It emits a test attribute in the same format as the native Attr method.
func attr(tb testing.TB, key, value string) {
	if strings.ContainsFunc(key, unicode.IsSpace) {
		tb.Errorf("disallowed whitespace in attribute key %q", key)
		return
	}
	if strings.ContainsAny(value, "\r\n") {
		tb.Errorf("disallowed newline in attribute value %q", value)
		return
	}
	// Emit the attribute in the same format as testing.TB.Attr()
	tb.Logf("=== ATTR  %s %s %s", tb.Name(), key, value)
}
