package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// JUnit XML structures
type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	Tests     int        `xml:"tests,attr"`
	Failures  int        `xml:"failures,attr"`
	Errors    int        `xml:"errors,attr"`
	Time      string     `xml:"time,attr"`
	TestCases []TestCase `xml:"testcase"`
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	ClassName string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      string   `xml:"time,attr"`
	File      string   `xml:"file,attr,omitempty"`
	Failure   *Failure `xml:"failure,omitempty"`
	Error     *Error   `xml:"error,omitempty"`
	Skipped   *Skipped `xml:"skipped,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

type Error struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

type Skipped struct {
	Message string `xml:"message,attr"`
}

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
	return filepath.WalkDir(tf.repoRoot, func(path string, d fs.DirEntry, err error) error {
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
	src, err := os.ReadFile(filePath)
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
			if strings.HasPrefix(funcName, "Test") || strings.HasPrefix(funcName, "Benchmark") || strings.HasPrefix(funcName, "Example") || strings.HasPrefix(funcName, "Fuzz") {
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
		key := fmt.Sprintf("%s.%s", packageName, testName)
		if file, exists := tf.testMap[key]; exists {
			return file
		}
	}

	// Try exact match without package name
	key := fmt.Sprintf("%s", testName)
	if file, exists := tf.testMap[key]; exists {
		return file
	}

	// Handle subtests and fuzz tests - try to match the parent test
	if strings.Contains(testName, "/") {
		parentTest := strings.Split(testName, "/")[0]
		
		// Try with package name
		if packageName != "" {
			key = fmt.Sprintf("%s.%s", packageName, parentTest)
			if file, exists := tf.testMap[key]; exists {
				return file
			}
		}
		
		// Try without package name
		key = fmt.Sprintf("%s", parentTest)
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

func main() {
	var (
		inputFile  = flag.String("input", "", "Path to JUnit XML file")
		outputFile = flag.String("output", "", "Path to output JUnit XML file (defaults to input file)")
		repoRoot   = flag.String("repo-root", ".", "Path to repository root")
		verbose    = flag.Bool("verbose", false, "Enable verbose output for debugging")
	)
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Input file is required")
	}

	if *outputFile == "" {
		*outputFile = *inputFile
	}

	// Read and parse the JUnit XML file
	xmlData, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	var testSuites TestSuites
	if err := xml.Unmarshal(xmlData, &testSuites); err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}

	// Initialize test finder and build test map
	finder := NewTestFinder(*repoRoot)
	if err := finder.BuildTestMap(); err != nil {
		log.Fatalf("Failed to build test map: %v", err)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Built test map with %d entries\n", len(finder.testMap))
		fmt.Fprintf(os.Stderr, "Sample entries:\n")
		for key, file := range finder.testMap {
			fmt.Fprintf(os.Stderr, "  %s -> %s\n", key, file)
		}
	}

	// Process each test case and add file information
	matched := 0
	total := 0
	for i := range testSuites.TestSuites {
		for j := range testSuites.TestSuites[i].TestCases {
			testCase := &testSuites.TestSuites[i].TestCases[j]
			total++
			
			// Skip if file is already set
			if testCase.File != "" {
				matched++
				continue
			}

			// Find the test file
			file := finder.FindTestFile(testCase.ClassName, testCase.Name)
			if file != "" {
				testCase.File = file
				matched++
				if *verbose {
					fmt.Fprintf(os.Stderr, "Matched: %s -> %s\n", testCase.Name, file)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Could not find file for test %s in class %s\n", 
					testCase.Name, testCase.ClassName)
			}
		}
	}

	// Marshal back to XML
	output, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal XML: %v", err)
	}

	// Add XML header
	xmlOutput := []byte(xml.Header + string(output))

	// Write to output file
	if err := os.WriteFile(*outputFile, xmlOutput, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Successfully enhanced JUnit XML file: %s (%d/%d test cases matched)\n", *outputFile, matched, total)
}
