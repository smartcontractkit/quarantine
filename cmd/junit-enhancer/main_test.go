package main

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestJUnitXMLProcessing(t *testing.T) {
	t.Parallel()

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
