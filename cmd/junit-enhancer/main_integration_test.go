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

	// Get absolute path to avoid directory change conflicts
	absModulePath, err := filepath.Abs(modulePath)
	if err != nil {
		t.Fatalf("Failed to get absolute path for %s: %v", modulePath, err)
	}

	// Create a temporary file for JUnit output
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("junit_%d.xml", time.Now().UnixNano()))

	// Check if gotestsum is available, if not use fallback
	if _, err := exec.LookPath("gotestsum"); err != nil {
		t.Fatalf("gotestsum not found")
	}

	// Run gotestsum with JUnit output in the target directory
	// #nosec G204 - tempFile path is controlled
	cmd := exec.Command(
		"gotestsum",
		"--junitfile",
		tempFile,
		"--format",
		"testname",
		"--",
		"-timeout",
		"5s",
		"--count",
		"1",
		"./...",
	)
	cmd.Dir = absModulePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		// It's ok if tests fail, we still want to process the output
		if exitError, ok := err.(*exec.ExitError); ok {
			t.Logf("gotestsum had non-zero exit code (continuing): %v\nOutput: %s", exitError, output)
		} else {
			t.Fatalf("Failed to run gotestsum: %v\nOutput: %s", err, output)
		}
	}

	// Verify the JUnit file was created
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

func verifyEnhancedOutput(t *testing.T, junitEnhancedFile string, expectedPathPrefix string) {
	t.Helper()

	// Read and verify the enhanced XML
	enhancedData, err := os.ReadFile(junitEnhancedFile) // #nosec G304 - file path is controlled by test
	if err != nil {
		t.Fatalf("Failed to read enhanced XML: %v", err)
	}

	var enhancedSuites JUnitTestSuites
	if err := xml.Unmarshal(enhancedData, &enhancedSuites); err != nil {
		t.Fatalf("Failed to parse enhanced XML: %v", err)
	}

	// Verify that test cases have file attributes
	filesAdded := 0
	totalTests := 0

	for _, suite := range enhancedSuites.Suites {
		for _, testCase := range suite.TestCases {
			totalTests++
			if testCase.File != "" {
				filesAdded++
				t.Logf("Test %s -> %s", testCase.Name, testCase.File)

				// Verify the file path is relative and reasonable
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

// Happy-case for testing mapping functionality
func TestIntegration_MainModule(t *testing.T) {
	t.Parallel()
	// Test the main test-fixture module
	testFixturePath := "./test-fixture"

	// Generate JUnit XML
	junitFile := runGoTestWithJUnit(t, testFixturePath)
	junitEnhancedFile := junitFile + ".enhanced"

	// Get the repository root (go up from cmd/junit-enhancer to repo root)
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get repo root: %v", err)
	}

	flags := getJunitFlags(t, junitFile, junitEnhancedFile, repoRoot)

	// Run the junit-enhancer tool
	cmd := exec.Command("go", flags...) // #nosec G204 - flags are controlled by test code
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("junit-enhancer failed: %v\nOutput: %s", err, output)
	}

	verifyEnhancedOutput(t, junitEnhancedFile, "cmd/junit-enhancer/test-fixture")

	// Clean up
	if err := os.Remove(junitFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitFile, err)
	}
	if err := os.Remove(junitEnhancedFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitEnhancedFile, err)
	}
}

// Happy-case for testing mapping functionality in a separate module
func TestIntegration_ServiceModule(t *testing.T) {
	t.Parallel()

	// Test the service module (separate Go module)
	servicePath := "./test-fixture/service"

	// Generate JUnit XML
	junitFile := runGoTestWithJUnit(t, servicePath)
	junitEnhancedFile := junitFile + ".enhanced"

	// Get the repository root
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get repo root: %v", err)
	}

	flags := getJunitFlags(t, junitFile, junitEnhancedFile, repoRoot)

	// Run the junit-enhancer tool
	cmd := exec.Command("go", flags...) // #nosec G204 - flags are controlled by test code

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("junit-enhancer failed: %v\nOutput: %s", err, output)
	}

	t.Logf("junit-enhancer output: %s", output)

	verifyEnhancedOutput(t, junitEnhancedFile, "cmd/junit-enhancer/test-fixture/service")

	if err := os.Remove(junitFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitFile, err)
	}
	if err := os.Remove(junitEnhancedFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitEnhancedFile, err)
	}
}

// Tests a module with intentional build failures
func TestIntegration_BuildFailureModule(t *testing.T) {
	t.Parallel()

	buildFailurePath := "./test-fixture/buildfailure"

	// Generate JUnit XML
	junitFile := runGoTestWithJUnit(t, buildFailurePath)
	junitEnhancedFile := junitFile + ".enhanced"

	// Get the repository root
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get repo root: %v", err)
	}

	flags := getJunitFlags(t, junitFile, junitEnhancedFile, repoRoot)

	// Run the junit-enhancer tool
	cmd := exec.Command("go", flags...) // #nosec G204 - flags are controlled by test code
	output, err := cmd.CombinedOutput()

	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode := exitError.ExitCode()
		if exitCode != 1 {
			t.Fatalf("junit-enhancer failed with unexpected exit code %d\nOutput: %s", exitCode, output)
		}
	} else {
		t.Fatalf(" Expected an exit code of 1 from junit-enhancer but got none.\nOutput: %s", output)
	}

	t.Logf("junit-enhancer output: %s", output)

	if err := os.Remove(junitFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitFile, err)
	}
	if err := os.Remove(junitEnhancedFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitEnhancedFile, err)
	}
}

// Tests a module with intentional TestMain failure
func TestIntegration_TestMainFailure(t *testing.T) {
	t.Parallel()

	// Test the broken module (separate Go module)
	testMainFailurePath := "./test-fixture/testmainfailure"

	// Generate JUnit XML
	junitFile := runGoTestWithJUnit(t, testMainFailurePath)
	junitEnhancedFile := junitFile + ".enhanced"

	// Get the repository root
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get repo root: %v", err)
	}

	flags := getJunitFlags(t, junitFile, junitEnhancedFile, repoRoot)

	// Run the junit-enhancer tool
	cmd := exec.Command("go", flags...) // #nosec G204 - flags are controlled by test code
	output, err := cmd.CombinedOutput()

	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode := exitError.ExitCode()
		if exitCode != 1 {
			t.Fatalf("junit-enhancer failed with unexpected exit code %d\nOutput: %s", exitCode, output)
		}
	} else {
		t.Fatalf(" Expected an exit code of 1 from junit-enhancer but got none.\nOutput: %s", output)
	}

	t.Logf("junit-enhancer output: %s", output)

	if err := os.Remove(junitFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitFile, err)
	}
	if err := os.Remove(junitEnhancedFile); err != nil {
		t.Logf("Failed to remove %s: %v", junitEnhancedFile, err)
	}
}
