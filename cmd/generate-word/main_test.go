package main

import (
	"io"
	// "log" // No longer needed
	"os"
	// "runtime" // No longer using runtime.Goexit()
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupCmdTest prepares the environment for running the main() function as a test.
// It redirects stdout and stderr, and sets os.Args.
// Now returns write ends of pipes as well.
func setupCmdTest(t *testing.T, args ...string) (rOut, wOut, rErr, wErr *os.File, oldStdout, oldStderr *os.File) {
	t.Helper()

	oldStdout = os.Stdout
	var err error
	rOut, wOut, err = os.Pipe()
	require.NoError(t, err)
	os.Stdout = wOut

	oldStderr = os.Stderr
	rErr, wErr, err = os.Pipe()
	require.NoError(t, err)
	os.Stderr = wErr

	os.Args = append([]string{"cmd"}, args...)

	return
}

// cleanupCmdTest restores stdout and stderr, closes pipes, and reads their content.
func cleanupCmdTest(t *testing.T, rOut, wOut, oldStdout, rErr, wErr, oldStderr *os.File) (stdout, stderr string) {
	t.Helper()

	wOut.Close()
	os.Stdout = oldStdout
	outBytes, err := io.ReadAll(rOut)
	require.NoError(t, err)
	stdout = string(outBytes)

	wErr.Close()
	os.Stderr = oldStderr
	errBytes, err := io.ReadAll(rErr)
	require.NoError(t, err)
	stderr = string(errBytes)
	return
}

// executeMain runs main and captures its output.
func executeMain(t *testing.T, args ...string) (stdout, stderr string) {
	t.Helper()
	// Correctly receive all 6 return values from setupCmdTest
	rOut, wOut, rErr, wErr, oldStdout, oldStderr := setupCmdTest(t, args...)

	// We need to capture panics if log.Fatalf is called, as it calls os.Exit
	// However, the output to stderr should be captured before os.Exit.
	// For this test structure, we'll check stderr. A more robust way for os.Exit
	// would involve running as a separate process.
	main() // If main calls log.Fatalf, execution stops here for this goroutine.

	return cleanupCmdTest(t, rOut, wOut, oldStdout, rErr, wErr, oldStderr)
}

func Test_cmd_default_syllables(t *testing.T) {
	stdout, stderr := executeMain(t, "-s3")
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output, got: %s", stdout)
}

func TestCmd_TypeMale_GeneratesOutput(t *testing.T) {
	stdout, stderr := executeMain(t, "--type", "male")
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output for --type male, got: %s", stdout)
}

func TestCmd_TypeFemale_GeneratesOutput(t *testing.T) {
	stdout, stderr := executeMain(t, "--type", "female")
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output for --type female, got: %s", stdout)
}

func TestCmd_TypeWord_GeneratesOutput(t *testing.T) {
	stdout, stderr := executeMain(t, "--type", "word")
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output for --type word, got: %s", stdout)
}

func TestCmd_DefaultType_IsWord_GeneratesOutput(t *testing.T) {
	// No --type flag, should default to "word"
	stdout, stderr := executeMain(t)
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output for default type, got: %s", stdout)
}

func TestCmd_TypeMale_WithSyllables_GeneratesOutput(t *testing.T) {
	stdout, stderr := executeMain(t, "--type", "male", "-s", "2")
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output for --type male -s 2, got: %s", stdout)
}

func TestCmd_TypeFemale_WithSyllables_GeneratesOutput(t *testing.T) {
	stdout, stderr := executeMain(t, "--type", "female", "-s", "5")
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output for --type female -s 5, got: %s", stdout)
}

func TestCmd_Syllables_WithTypeWord_GeneratesOutput(t *testing.T) {
	stdout, stderr := executeMain(t, "--type", "word", "-s", "4")
	assert.Empty(t, stderr, "stderr should be empty for successful execution")
	assert.Greater(t, len(strings.TrimSpace(stdout)), 0, "expected non-empty output for --type word -s 4, got: %s", stdout)
}

func TestCmd_TypeInvalid_PrintsError(t *testing.T) {
	// main calls log.Fatalf which prints to stderr and exits.
	// We can't easily check the exit code without `exec` and a helper process.
	// However, we can check if stderr contains the expected error message.
	// The `executeMain` helper will capture stderr.
	// Note: log.Fatalf output includes date and time. We check for a substring.

	// log.Fatalf in main() writes to os.Stderr, which is redirected by setupCmdTest.
	// So, we don't need to (and shouldn't) set log.SetOutput(io.Discard) here if we want to capture the log output.

	// Note: log.Fatalf in main calls os.Exit(1).
	// This test structure might not cleanly capture that exit for assertion,
	// but it should capture stderr output before exit.
	// A more robust way to test os.Exit is to exec the command.

	oldOsExit := osExit // main.osExit (package variable)
	var exitCode int
	var exited bool // Flag to check if our osExit was called
	osExit = func(code int) { // redirect main.osExit to this func
		exitCode = code
		exited = true
		// Do not call runtime.Goexit() or os.Exit() here to allow test to complete assertions
	}
	defer func() { osExit = oldOsExit }() // restore main.osExit

	stdout, stderr := executeMain(t, "--type", "invalidtype")

	// Assertions
	require.True(t, exited, "osExit mock should have been called")
	assert.Equal(t, 1, exitCode, "Expected os.Exit(1) to be called by main.osExit")

	t.Logf("TestCmd_TypeInvalid_PrintsError captured stdout: %q", stdout)
	t.Logf("TestCmd_TypeInvalid_PrintsError captured stderr: %q", stderr)

	// Even though main() tried to exit, the output should still be captured by the pipes
	// before the mock osExit returns.
	assert.Empty(t, stdout, "stdout should be empty on error")
	assert.Contains(t, stderr, "Error processing command-line options: invalid word type 'invalidtype'. Valid types are: word, male, female", "stderr should contain the specific invalid type error message")
}

// Allow mocking os.Exit for tests
// var osExit = os.Exit // This is defined in main.go as a package variable, so tests can modify it.
// It is defined in main.go as a package variable, so tests can modify it.

// To make the TestCmd_TypeInvalid_PrintsError more robust, we need to ensure that
// log.Fatalf in main.go is what we're testing.
// The previous test for invalid type will capture stderr.
// If main calls log.Fatalf, the test will proceed to check stderr.
// The `executeMain` structure is such that `main()` is called directly.
// If `log.Fatalf` occurs, it calls `os.Exit(1)`.
// The test `TestCmd_TypeInvalid_PrintsError` needs to handle this.
// One way is to replace `log.Fatalf` with a mock in tests, or check for `os.Exit`
// by running the command as a separate process.
// The current `TestCmd_TypeInvalid_PrintsError` will work because `log.Fatalf` prints to Stderr *before* exiting.
// The `cleanupCmdTest` will read Stderr.
// However, the test runner itself might be affected by `os.Exit(1)`.
// The `osExit` variable trick is a common way to handle this without subprocessing.
// Let's make sure the `executeMain` function is used in `TestCmd_TypeInvalid_PrintsError`.
// It is.
// The `log.SetOutput(io.Discard)` is to prevent the log.Fatalf from printing to the actual test runner's stderr.
// The `defer log.SetOutput(originalLogOutput)` restores it.
// This setup for `TestCmd_TypeInvalid_PrintsError` should be reasonably robust for the current `main.go` structure.
// The original test `Test_cmd_should_generate_a_word_when_called_with_n3`
// needs to be renamed or updated to fit the new helper structure.
// I will rename it to `Test_cmd_default_syllables` and use the helpers.
// The `os.Args` manipulation is now handled by `setupCmdTest`.
// The `bytes.Buffer` for capturing output is also handled by helpers.
// The `io.ReadAll` and `w.Close()` are also handled.
// The old test needs to be fully replaced.
// The previous code block had an incomplete thought process for the invalid type test.
// The `log.SetOutput(io.Discard)` and `osExit` mocking in `TestCmd_TypeInvalid_PrintsError`
// are the correct patterns to make this testable without `exec.Command`.
// The `executeMain` function will call `main()`. If `log.Fatalf` is hit in `main()` due to the invalid type,
// our mocked `osExit` will be called, `exitCode` will be set, and the function will not actually exit.
// Then, `cleanupCmdTest` will read `stderr`. This is a good approach.

// Need to import "log" for the invalid type test.
// Need to import "bytes" for capturing stdout/stderr if that was the plan (os.Pipe is used, so bytes.Buffer not strictly needed here).
// Need to import "strings" for TrimSpace and Contains.
// All necessary imports seem to be covered by the diff.
