package main

import (
	"io"
	"os"
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockExit stores the exit code from os.Exit
var mockExitCode int

// mockOsExit replaces os.Exit for testing, capturing the exit code.
func mockOsExit(code int) {
	mockExitCode = code
	panic("os.Exit called") // Panic to stop execution like os.Exit would, but allow recovery in test
}

// runMainWithArgs executes the main function with given arguments and captures its output and exit code.
func runMainWithArgs(args []string) (stdoutStr, stderrStr string, exitCode int) {
	// Reset mockExitCode for each run
	mockExitCode = 0

	oldOsArgs := os.Args
	oldOsStdout := os.Stdout
	oldOsStderr := os.Stderr
	// Store the original osExit and defer its restoration
	originalOsExit := osExit
	osExit = mockOsExit // Assign the mock function
	defer func() {
		os.Args = oldOsArgs
		os.Stdout = oldOsStdout
		os.Stderr = oldOsStderr
		osExit = originalOsExit // Restore the original os.Exit
	}()

	os.Args = append([]string{"cmd"}, args...) // Prepend command name

	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr

	// Recover from panic caused by mockOsExit
	defer func() {
		if r := recover(); r != nil {
			if r != "os.Exit called" {
				panic(r) // Re-panic if it's not the one from mockOsExit
			}
		}
		exitCode = mockExitCode
	}()

	main() // Execute the application's main function

	wOut.Close()
	wErr.Close()

	stdoutBytes, _ := io.ReadAll(rOut)
	stderrBytes, _ := io.ReadAll(rErr)

	return string(stdoutBytes), string(stderrBytes), mockExitCode
}

// It's necessary to allow os.Exit to be replaced for testing.
// This is a common pattern.
var osExit = os.Exit

func TestMainCLI(t *testing.T) {
	tests := []struct {
		name                 string
		args                 []string
		expectedExitCode     int
		expectedStdoutEmpty  bool // Check if stdout should be empty
		expectedStderrContains string
	}{
		{
			name:                "no args (default syllables)",
			args:                []string{}, // Default behavior uses 2 syllables
			expectedExitCode:    0,
			expectedStdoutEmpty: false,
		},
		{
			name:                "syllables 3",
			args:                []string{"-s3"},
			expectedExitCode:    0,
			expectedStdoutEmpty: false,
		},
		{
			name:                "gender male",
			args:                []string{"--gender", "male"},
			expectedExitCode:    0,
			expectedStdoutEmpty: false,
		},
		{
			name:                "gender female short flag",
			args:                []string{"-g", "female"},
			expectedExitCode:    0,
			expectedStdoutEmpty: false,
		},
		{
			name:                "gender male syllables 3",
			args:                []string{"--gender", "male", "--syllables", "3"},
			expectedExitCode:    0,
			expectedStdoutEmpty: false,
		},
		{
			name:                "gender female syllables 2",
			args:                []string{"--gender", "female", "-s", "2"},
			expectedExitCode:    0,
			expectedStdoutEmpty: false,
		},
		{
			name:                 "invalid gender",
			args:                 []string{"--gender", "xyz"},
			expectedExitCode:     1,
			expectedStdoutEmpty:  true,
			expectedStderrContains: "invalid gender: xyz",
		},
		{
			name:                 "gender male syllables 0",
			args:                 []string{"--gender", "male", "--syllables", "0"},
			expectedExitCode:     1,
			expectedStdoutEmpty:  true,
			expectedStderrContains: "number of syllables must be one or greater",
		},
		{
			name:                 "syllables 0 (no gender)",
			args:                 []string{"-s0"},
			expectedExitCode:     1,
			expectedStdoutEmpty:  true,
			expectedStderrContains: "number of syllables must be one or greater",
		},
		{
			name:                 "help flag",
			args:                 []string{"--help"},
			expectedExitCode:     0, // go-flags exits with 0 on --help
			expectedStdoutEmpty:  false, // Help message goes to stdout
			expectedStderrContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runMainWithArgs(tt.args)

			assert.Equal(t, tt.expectedExitCode, exitCode, "unexpected exit code")

			if tt.expectedStdoutEmpty {
				assert.Empty(t, stdout, "expected stdout to be empty")
			} else {
				assert.NotEmpty(t, stdout, "expected stdout to be non-empty")
			}

			if tt.expectedStderrContains != "" {
				require.Contains(t, stderr, tt.expectedStderrContains, "stderr does not contain expected message")
			}
			// The empty 'else' block that caused a lint error was here. It has been removed.
			// If tt.expectedStderrContains is empty, no assertion is made on stderr,
			// allowing for cases where stderr might have non-critical output or be empty.
		})
	}
}
