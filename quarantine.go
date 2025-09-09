// Package quarantine provides a way to mark tests as flaky.
package quarantine

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode"
)

// RunQuarantinedTestsEnvVar is the environment variable that controls whether to run quarantined tests.
const RunQuarantinedTestsEnvVar = "RUN_QUARANTINED_TESTS"

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

	explanationStr := fmt.Sprintf(
		"known flaky test. Ticket %s",
		ticket,
	)
	classifiedStr := "Classified by branch-out (https://github.com/smartcontractkit/branch-out)"
	// tb.Attr("flaky_test", ticket) - only compatible with Go 1.25+
	attr(tb, "flaky_test", ticket)
	//nolint:forbidigo // Config doesn't make sense here
	if os.Getenv(RunQuarantinedTestsEnvVar) != "true" {
		tb.Skipf(
			"Skipping %s. To run quarantined tests, set the %s environment variable to true.\n%s",
			explanationStr,
			RunQuarantinedTestsEnvVar,
			classifiedStr,
		)
	} else {
		tb.Logf("Running %s", explanationStr)
		tb.Cleanup(func() {
			tb.Logf(
				"Test is marked as quarantined, but still ran. %s. To skip quarantined tests, ensure the %s environment variable is set to false.\n%s",
				explanationStr,
				RunQuarantinedTestsEnvVar,
				classifiedStr,
			)
		})
	}
}
