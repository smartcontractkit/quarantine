package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTestFinder(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create a test file
	testFileContent := `package main

import "testing"

func TestExample(t *testing.T) {
	// Test implementation
}

func TestExampleSubtest(t *testing.T) {
	t.Run("subtest1", func(t *testing.T) {
		// Subtest implementation
	})
}

func FuzzMessageHasher(f *testing.F) {
	// Fuzz test implementation
}
`
	testFilePath := filepath.Join(tempDir, "example_test.go")
	err := os.WriteFile(testFilePath, []byte(testFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Initialize test finder
	finder := NewTestFinder(tempDir)
	err = finder.BuildTestMap()
	if err != nil {
		t.Fatalf("Failed to build test map: %v", err)
	}

	// Test finding the test file
	file := finder.FindTestFile("main", "TestExample")
	expected := "example_test.go"
	if file != expected {
		t.Errorf("Expected file %s, got %s", expected, file)
	}

	// Test finding subtest
	file = finder.FindTestFile("main", "TestExampleSubtest/subtest1")
	if file != expected {
		t.Errorf("Expected file %s for subtest, got %s", expected, file)
	}

	// Test finding fuzz test
	file = finder.FindTestFile("main", "FuzzMessageHasher")
	if file != expected {
		t.Errorf("Expected file %s for fuzz test, got %s", expected, file)
	}

	// Test finding fuzz test with generated input
	file = finder.FindTestFile("main", "FuzzMessageHasher/ff5b0490c60257f0")
	if file != expected {
		t.Errorf("Expected file %s for fuzz test with input, got %s", expected, file)
	}
}

func TestJUnitXMLProcessing(t *testing.T) {
	// Test XML parsing and marshaling
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="main" tests="1" failures="0" errors="0" time="0.001">
    <testcase classname="github.com/example/main" name="TestExample" time="0.000000"></testcase>
  </testsuite>
</testsuites>`

	var testSuites JUnitTestSuites
	err := xml.Unmarshal([]byte(xmlContent), &testSuites)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Verify parsing
	if len(testSuites.Suites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(testSuites.Suites))
	}

	if len(testSuites.Suites[0].TestCases) != 1 {
		t.Errorf("Expected 1 test case, got %d", len(testSuites.Suites[0].TestCases))
	}

	testCase := testSuites.Suites[0].TestCases[0]
	if testCase.Name != "TestExample" {
		t.Errorf("Expected test name 'TestExample', got '%s'", testCase.Name)
	}

	// Add file attribute and test marshaling
	testSuites.Suites[0].TestCases[0].File = "main_test.go"

	output, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	// Verify file attribute is included
	outputStr := string(output)
	if !strings.Contains(outputStr, `file="main_test.go"`) {
		t.Errorf("Expected file attribute in output, got: %s", outputStr)
	}
}
