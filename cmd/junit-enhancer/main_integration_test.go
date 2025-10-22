package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Test helper to run go test and capture JUnit XML output using gotestsum
func runGoTestWithJUnit(t *testing.T, modulePath string) string {
	t.Helper()

	absModulePath, err := filepath.Abs(modulePath)
	if err != nil {
		t.Fatalf("Failed to get absolute path for %s: %v", modulePath, err)
	}

	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("junit_%d.xml", time.Now().UnixNano()))

	if _, err := exec.LookPath("gotestsum"); err != nil {
		t.Fatalf("gotestsum not found")
	}

	// #nosec G204 - tempFile path is controlled
	cmd := exec.Command(
		"gotestsum",
		"--junitfile", tempFile,
		"--format", "testname",
		"--",
		"-timeout", "5s",
		"--count", "1",
		"./...",
	)
	cmd.Dir = absModulePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			t.Logf("gotestsum had non-zero exit code (continuing): %v\nOutput: %s", exitError, output)
		} else {
			t.Fatalf("Failed to run gotestsum: %v\nOutput: %s", err, output)
		}
	}

	if _, err := os.Stat(tempFile); err != nil {
		t.Fatalf("JUnit file was not created: %v", err)
	}

	return tempFile
}

func getJunitFlags(t *testing.T, inputFile, outputFile, repoRoot string) []string {
	t.Helper()

	flags := []string{
		"run",
		".",
		"-input", inputFile,
		"-output", outputFile,
		"-repo-root", repoRoot,
	}

	if testing.Verbose() {
		flags = append(flags, "-verbose")
	}

	return flags
}

// verifyEnhancedOutput checks that the enhanced JUnit XML has file paths added correctly
// and that TestMain entries are filtered out.
func verifyEnhancedOutput(t *testing.T, junitEnhancedFile string, expectedPathPrefix string) {
	t.Helper()

	enhancedData, err := os.ReadFile(junitEnhancedFile) // #nosec G304 - controlled by test
	if err != nil {
		t.Fatalf("Failed to read enhanced XML: %v", err)
	}

	var enhancedSuites JUnitTestSuites
	if err := xml.Unmarshal(enhancedData, &enhancedSuites); err != nil {
		t.Fatalf("Failed to parse enhanced XML: %v", err)
	}

	filesAdded := 0
	totalTests := 0

	for _, suite := range enhancedSuites.Suites {
		for _, testCase := range suite.TestCases {
			totalTests++
			if testCase.File != "" {
				filesAdded++
				t.Logf("Test %s -> %s", testCase.Name, testCase.File)

				if !strings.HasPrefix(testCase.File, expectedPathPrefix) {
					t.Errorf("File path should start with %q, got: %s", expectedPathPrefix, testCase.File)
				}
				if !strings.HasSuffix(testCase.File, "_test.go") {
					t.Errorf("File path should point to a test file, got: %s", testCase.File)
				}
				if testCase.Classname == "" && testCase.Name == "TestMain" {
					t.Errorf("TestMain should have been filtered out. Found for suite: %s", suite.Name)
				}
			}
		}
	}

	t.Logf("Enhanced %d out of %d test cases with file paths", filesAdded, totalTests)

	if filesAdded == 0 {
		t.Error("Expected at least some test cases to have file paths added")
	}
}

func TestIntegration_JUnitEnhancer(t *testing.T) {
	t.Parallel()

	// Resolve repo root once for all cases (go up from cmd/junit-enhancer to repo root)
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get repo root: %v", err)
	}

	type tc struct {
		name           string
		modulePath     string
		expectedPrefix string
		expectExitCode int // expected exit code from `go run .`; 0 for success, 1 for intentional failures
	}

	tests := []tc{
		{
			name:           "MainModule",
			modulePath:     "./test-fixture",
			expectedPrefix: "cmd/junit-enhancer/test-fixture",
			expectExitCode: 0,
		},
		{
			name:           "ServiceModule",
			modulePath:     "./test-fixture/service",
			expectedPrefix: "cmd/junit-enhancer/test-fixture/service",
			expectExitCode: 0,
		},
		{
			name:           "BuildFailureModule",
			modulePath:     "./test-fixture/buildfailure",
			expectedPrefix: "cmd/junit-enhancer/test-fixture/buildfailure",
			expectExitCode: 1,
		},
		{
			name:           "TestMainFailure",
			modulePath:     "./test-fixture/testmainfailure",
			expectedPrefix: "cmd/junit-enhancer/test-fixture/testmainfailure",
			expectExitCode: 1,
		},
		{
			name:           "TimeoutFailure",
			modulePath:     "./test-fixture/timeout",
			expectedPrefix: "cmd/junit-enhancer/test-fixture/timeout",
			expectExitCode: 1,
		},
	}

	for _, tt := range tests {
		// capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Generate JUnit XML
			junitFile := runGoTestWithJUnit(t, tt.modulePath)
			junitEnhancedFile := junitFile + ".enhanced"

			// Build flags and run the enhancer
			flags := getJunitFlags(t, junitFile, junitEnhancedFile, repoRoot)

			// #nosec G204 - flags are controlled by test code
			cmd := exec.Command("go", flags...)
			output, err := cmd.CombinedOutput()

			switch tt.expectExitCode {
			case 0:
				if err != nil {
					t.Fatalf("junit-enhancer failed: %v\nOutput: %s", err, output)
				}
			default:
				exitErr, ok := err.(*exec.ExitError)
				if !ok {
					t.Fatalf(
						"Expected exit code %d but command did not return an ExitError.\nOutput: %s",
						tt.expectExitCode,
						output,
					)
				}
				if code := exitErr.ExitCode(); code != tt.expectExitCode {
					t.Fatalf("Expected exit code %d, got %d\nOutput: %s", tt.expectExitCode, code, output)
				}
			}

			t.Logf("junit-enhancer output (%s): %s", tt.name, output)

			// Verify enhanced XML
			verifyEnhancedOutput(t, junitEnhancedFile, tt.expectedPrefix)

			// Cleanup
			if err := os.Remove(junitFile); err != nil {
				t.Logf("Failed to remove %s: %v", junitFile, err)
			}
			if err := os.Remove(junitEnhancedFile); err != nil {
				t.Logf("Failed to remove %s: %v", junitEnhancedFile, err)
			}
		})
	}
}
