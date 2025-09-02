package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTestFinder(t *testing.T) {
	t.Parallel()
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
	err := os.WriteFile(testFilePath, []byte(testFileContent), 0600)
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
