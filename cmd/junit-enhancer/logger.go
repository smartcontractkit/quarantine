// Package main implements a JUnit XML enhancer that adds file path information to test cases.
package main

import (
	"fmt"
	"os"
)

// Logger provides GitHub Actions-aware logging functionality
type Logger struct {
	isGitHubActions bool
	verbose         bool
}

// NewLogger creates a new logger instance
func NewLogger(verbose bool) *Logger {
	isGHA := isRunningInGitHubActions()

	// Automatically enable verbose mode if we're in GitHub Actions and RUNNER_DEBUG is enabled
	if isGHA && os.Getenv("RUNNER_DEBUG") == "1" {
		verbose = true
	}

	return &Logger{
		isGitHubActions: isGHA,
		verbose:         verbose,
	}
}

// isRunningInGitHubActions detects if the code is running in GitHub Actions
func isRunningInGitHubActions() bool {
	// GitHub Actions sets GITHUB_ACTIONS=true
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return true
	}
	// Alternative check: CI is commonly set in GitHub Actions
	if os.Getenv("CI") == "true" && os.Getenv("GITHUB_WORKFLOW") != "" {
		return true
	}
	return false
}

// Logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.verbose {
		return
	}

	message := fmt.Sprintf(format, args...)
	if l.isGitHubActions {
		fmt.Fprintf(os.Stderr, "::debug::%s\n", message)
	} else {
		fmt.Fprintf(os.Stderr, "Debug: %s\n", message)
	}
}

// Logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s\n", message)
}

// Logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if l.isGitHubActions {
		fmt.Fprintf(os.Stderr, "::warning::%s\n", message)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", message)
	}
}

// Logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if l.isGitHubActions {
		fmt.Fprintf(os.Stderr, "::error::%s\n", message)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	}
}

// Logs a fatal error message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Error(format, args...)
	os.Exit(1)
}
