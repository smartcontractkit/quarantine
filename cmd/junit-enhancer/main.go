// Package main implements a JUnit XML enhancer that adds file path information to test cases.
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
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

	var testSuites JUnitTestSuites
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

	// Filter out all TestMain testcases
	for i := range testSuites.Suites {
		var filtered []JUnitTestCase
		for j := range testSuites.Suites[i].TestCases {
			testCase := &testSuites.Suites[i].TestCases[j]
			if testCase.Classname == "" && testCase.Name == "TestMain" {
				continue
			}
			filtered = append(filtered, *testCase)
		}
		testSuites.Suites[i].TestCases = filtered
	}

	// Process each test case and add file information
	matched := 0
	total := 0
	for i := range testSuites.Suites {
		for j := range testSuites.Suites[i].TestCases {
			testCase := &testSuites.Suites[i].TestCases[j]
			total++

			// Skip if file is already set
			if testCase.File != "" {
				matched++
				continue
			}

			// Find the test file
			file := finder.FindTestFile(testCase.Classname, testCase.Name)
			if file != "" {
				testCase.File = file
				matched++
				if *verbose {
					fmt.Fprintf(os.Stderr, "Matched: %s -> %s\n", testCase.Name, file)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Could not find file for test %s in class %s\n",
					testCase.Name, testCase.Classname)
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
	if err := os.WriteFile(*outputFile, xmlOutput, 0600); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Successfully enhanced JUnit XML file: %s (%d/%d test cases matched)\n", *outputFile, matched, total)
}
