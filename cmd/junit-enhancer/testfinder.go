package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TestFinder helps locate test files and functions
type TestFinder struct {
	repoRoot string
	fileSet  *token.FileSet
	testMap  map[string]string // maps test name to file path
}

func NewTestFinder(repoRoot string) *TestFinder {
	return &TestFinder{
		repoRoot: repoRoot,
		fileSet:  token.NewFileSet(),
		testMap:  make(map[string]string),
	}
}

// BuildTestMap scans the repository for Go test files and builds a map of test names to file paths
func (tf *TestFinder) BuildTestMap() error {
	return filepath.WalkDir(tf.repoRoot, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and non-test files
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Skip vendor directories
		if strings.Contains(path, "/vendor/") {
			return nil
		}

		return tf.parseTestFile(path)
	})
}

// parseTestFile parses a Go test file and extracts test function names
func (tf *TestFinder) parseTestFile(filePath string) error {
	src, err := os.ReadFile(filePath) // #nosec G304 - filePath is controlled by filepath.WalkDir
	if err != nil {
		return err
	}

	// Parse the Go source file
	file, err := parser.ParseFile(tf.fileSet, filePath, src, parser.ParseComments)
	if err != nil {
		// Log parsing errors but don't fail completely
		fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", filePath, err)
		return nil
	}

	// Get relative path from repo root
	relPath, err := filepath.Rel(tf.repoRoot, filePath)
	if err != nil {
		return err
	}

	// Extract package name for classname matching
	packageName := file.Name.Name

	// Find test functions
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.IsExported() {
			funcName := fn.Name.Name
			if strings.HasPrefix(funcName, "Test") ||
				strings.HasPrefix(funcName, "Benchmark") ||
				strings.HasPrefix(funcName, "Example") ||
				strings.HasPrefix(funcName, "Fuzz") {
				// Create a key that matches the classname pattern
				key := fmt.Sprintf("%s.%s", packageName, funcName)
				tf.testMap[key] = relPath
			}
		}
	}

	return nil
}

// FindTestFile attempts to find the file for a given test case
func (tf *TestFinder) FindTestFile(className, testName string) string {
	// Extract package name from classname (remove module path prefix)
	parts := strings.Split(className, "/")
	var packageName string
	if len(parts) > 0 {
		packageName = parts[len(parts)-1]
	}

	// Try exact match with package name first
	if packageName != "" {
		key := packageName + "." + testName
		if file, exists := tf.testMap[key]; exists {
			return file
		}
	}

	// Try exact match without package name
	key := testName
	if file, exists := tf.testMap[key]; exists {
		return file
	}

	// Handle subtests and fuzz tests - try to match the parent test
	if strings.Contains(testName, "/") {
		parentTest := strings.Split(testName, "/")[0]

		// Try with package name
		if packageName != "" {
			key = packageName + "." + parentTest
			if file, exists := tf.testMap[key]; exists {
				return file
			}
		}

		// Try without package name
		key = parentTest
		if file, exists := tf.testMap[key]; exists {
			return file
		}
	}

	// Fallback: search for any file containing the test name or parent test name
	searchNames := []string{testName}
	if strings.Contains(testName, "/") {
		parentTest := strings.Split(testName, "/")[0]
		searchNames = append(searchNames, parentTest)
	}

	for _, searchName := range searchNames {
		for testKey, file := range tf.testMap {
			if strings.Contains(testKey, searchName) {
				return file
			}
		}
	}

	return ""
}
