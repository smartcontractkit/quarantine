# JUnit Enhancer CLI

A Go CLI tool that enhances JUnit XML test reports by adding file paths to test case entries based on the classname and test name.

## Features

- Reads JUnit XML test reports
- Matches test cases to their corresponding Go test files using classname and test name
- Adds relative file paths to test case entries
- Supports subtests
- Works with multiple Go modules in a repository

## Usage

```bash
# Build the CLI tool
go build -o junit-enhancer

# Enhance a JUnit XML file (modifies the file in place)
./junit-enhancer -input /path/to/junit.xml -repo-root /path/to/repository

# Enhance a JUnit XML file and save to a different location
./junit-enhancer -input /path/to/junit.xml -output /path/to/enhanced-junit.xml -repo-root /path/to/repository
```

## Command Line Options

- `-input`: Path to the input JUnit XML file (required)
- `-output`: Path to the output JUnit XML file (optional, defaults to input file)
- `-repo-root`: Path to the repository root (optional, defaults to current directory)

## How It Works

1. **Parse Test Files**: The tool scans the repository for `*_test.go` files and builds a map of test function names to file paths
2. **Match Test Cases**: For each test case in the JUnit XML, it attempts to match the classname and test name to a test file
3. **Handle Subtests**: For subtests (containing `/` in the name), it matches against the parent test function
4. **Add File Attribute**: Adds a `file` attribute to each test case with the relative path from the repository root

## Example

**Input JUnit XML:**
```xml
<testcase classname="github.com/smartcontractkit/chainlink/v2/core/bridges" name="TestBridgeTypeRequest" time="0.000000"></testcase>
```

**Enhanced Output:**
```xml
<testcase classname="github.com/smartcontractkit/chainlink/v2/core/bridges" name="TestBridgeTypeRequest" time="0.000000" file="core/bridges/bridge_test.go"></testcase>
```

## Integration with Multiple Go Modules

This tool is designed to be run once per Go module in a repository. For repositories with multiple Go modules:

1. Run the tool from each module's directory
2. Set the `-repo-root` flag to point to the repository root
3. The tool will find test files relative to the repository root, not the module root

## Building

```bash
cd cmd/junit-enhancer
go build -o junit-enhancer
```

## Dependencies

- `golang.org/x/tools` for Go source code parsing
