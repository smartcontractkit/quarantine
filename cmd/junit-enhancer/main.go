// Package main implements a JUnit XML enhancer that adds file path information to test cases.
package main

import (
	"encoding/xml"
	"flag"
	"os"
)

// JUnit XML structures based on gotestsum's internal/junitxml package
// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName  xml.Name         `xml:"testsuites"`
	Name     string           `xml:"name,attr,omitempty"`
	Tests    int              `xml:"tests,attr"`
	Failures int              `xml:"failures,attr"`
	Errors   int              `xml:"errors,attr"`
	Time     string           `xml:"time,attr"`
	Suites   []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite is a single JUnit test suite which may contain many testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Skipped    int             `xml:"skipped,attr,omitempty"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase `xml:"testcase"`
	Timestamp  string          `xml:"timestamp,attr"`
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	Time        string            `xml:"time,attr"`
	File        string            `xml:"file,attr,omitempty"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

// processTestCaseFilePath attempts to add file path information to a test case
// Returns true if the test case was matched with a file path (or already had one)
func processTestCaseFilePath(tCase *JUnitTestCase, finder *TestFinder, logger *Logger) bool {
	// Skip if file is already set
	if tCase.File != "" {
		logger.Warning("File already populated for %s. Skipping.", tCase.Name)
		return true
	}

	// Find the test file
	file := finder.FindTestFile(tCase.Classname, tCase.Name)
	if file != "" {
		tCase.File = file
		logger.Debug("Matched: %s -> %s", tCase.Name, file)
		return true
	}

	logger.Warning("Could not find file for test %s in class %s", tCase.Name, tCase.Classname)
	return false
}

func main() {
	var (
		inputFile  = flag.String("input", "", "Path to JUnit XML file")
		outputFile = flag.String("output", "", "Path to output JUnit XML file (defaults to input file)")
		repoRoot   = flag.String("repo-root", ".", "Path to repository root")
		verbose    = flag.Bool("verbose", false, "Enable verbose output for debugging")
	)
	flag.Parse()

	// Initialize logger
	logger := NewLogger(*verbose)

	if *inputFile == "" {
		logger.Fatal("Input file is required")
	}

	if *outputFile == "" {
		*outputFile = *inputFile
	}

	// Read and parse the JUnit XML file
	xmlData, err := os.ReadFile(*inputFile)
	if err != nil {
		logger.Fatal("Failed to read input file: %v", err)
	}

	var testSuites JUnitTestSuites
	if err := xml.Unmarshal(xmlData, &testSuites); err != nil {
		logger.Fatal("Failed to parse XML: %v", err)
	}

	// Initialize test finder and build test map
	finder := NewTestFinder(*repoRoot)
	if err := finder.BuildTestMap(); err != nil {
		logger.Fatal("Failed to build test map: %v", err)
	}

	logger.Debug("Built test map with %d entries", len(finder.testMap))
	logger.Debug("Sample entries:")
	for key, file := range finder.testMap {
		logger.Debug("  %s -> %s", key, file)
	}

	// Process test suites: filter and add file information in one pass
	matched := 0
	total := 0
	var filteredSuites []JUnitTestSuite

	for i, suite := range testSuites.Suites {
		if (suite.Tests == 0) && (suite.Name == "") {
			logger.Warning("Skipping unnamed and empty test suite at index %d", i)
			continue
		}

		var filteredTestCases []JUnitTestCase
		for _, tCase := range suite.TestCases {
			if tCase.Classname == "" && tCase.Name == "TestMain" {
				// Skip TestMain test cases which are typically due to
				// package-level setup/teardown issues, like compilation errors.
				logger.Warning("Skipping TestMain test case in suite %s", suite.Name)
				if tCase.Failure != nil {
					logger.Error("%s TestMain failure: %s", suite.Name, tCase.Failure.Contents)
				}
				continue
			}

			// Process file information for valid test cases
			total++
			if processTestCaseFilePath(&tCase, finder, logger) {
				matched++
			}

			filteredTestCases = append(filteredTestCases, tCase)
		}
		suite.TestCases = filteredTestCases
		filteredSuites = append(filteredSuites, suite)
	}

	testSuites.Suites = filteredSuites

	// Marshal back to XML
	output, err := xml.MarshalIndent(testSuites, "", "\t")
	if err != nil {
		logger.Fatal("Failed to marshal XML: %v", err)
	}

	// Add XML header
	xmlOutput := []byte(xml.Header + string(output))

	// Write to output file
	if err := os.WriteFile(*outputFile, xmlOutput, 0600); err != nil {
		logger.Fatal("Failed to write output file: %v", err)
	}

	logger.Info("Successfully enhanced JUnit XML file: %s (%d/%d test cases matched)", *outputFile, matched, total)
}
