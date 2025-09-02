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

	// Change to the module directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(modulePath); err != nil {
		t.Fatalf("Failed to change to module directory %s: %v", modulePath, err)
	}

	// Create a temporary file for JUnit output
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("junit_%d.xml", time.Now().UnixNano()))

	// Check if gotestsum is available, if not use fallback
	if _, err := exec.LookPath("gotestsum"); err != nil {
		t.Fatalf("gotestsum not found")
	}

	// Run gotestsum with JUnit output
	cmd := exec.Command("gotestsum", "--junitfile", tempFile, "--format", "testname", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// It's ok if tests fail, we still want to process the output
		if exitError, ok := err.(*exec.ExitError); ok {
			t.Logf("gotestsum had non-zero exit code: %v\nOutput: %s", exitError, output)
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

func TestIntegration_MainModule(t *testing.T) {
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

	flags := []string{
		"run",
		"main.go",
		"-input", junitFile,
		"-output", junitEnhancedFile,
		"-repo-root", repoRoot,
		"-verbose",
	}

	// Run the junit-enhancer tool
	cmd := exec.Command("go", flags...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("junit-enhancer failed: %v\nOutput: %s", err, output)
	}

	t.Logf("junit-enhancer output: %s", output)

	verifyEnhancedOutput(t, junitEnhancedFile, "cmd/junit-enhancer/test-fixture")

	// Clean up
	os.Remove(junitFile)
	os.Remove(junitFile + ".enhanced")
}

func TestIntegration_ServiceModule(t *testing.T) {
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

	flags := []string{
		"run",
		"main.go",
		"-input", junitFile,
		"-output", junitEnhancedFile,
		"-repo-root", repoRoot,
		"-verbose",
	}

	// Run the junit-enhancer tool
	cmd := exec.Command("go", flags...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("junit-enhancer failed: %v\nOutput: %s", err, output)
	}

	t.Logf("junit-enhancer output: %s", output)

	verifyEnhancedOutput(t, junitEnhancedFile, "cmd/junit-enhancer/test-fixture/service")

	os.Remove(junitFile)
	os.Remove(junitEnhancedFile)
}

func verifyEnhancedOutput(t *testing.T, junitEnhancedFile string, expectedPathPrefix string) {
	// Read and verify the enhanced XML
	enhancedData, err := os.ReadFile(junitEnhancedFile)
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
			}
		}
	}

	t.Logf("Enhanced %d out of %d test cases with file paths", filesAdded, totalTests)

	if filesAdded == 0 {
		t.Error("Expected at least some test cases to have file paths added")
	}
}
