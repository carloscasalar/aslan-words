package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_cmd_should_generate_a_word_when_called_with_n3(t *testing.T) {
	// Redirect stdout to capture the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Given
	// the command line arguments are "-n3"
	os.Args = []string{"path-to-cmd", "-s3"}

	// When
	// Call the main function
	main()

	// Restore stdout and read the output
	w.Close()
	os.Stdout = old
	output, err := io.ReadAll(r)

	// Check the output
	require.NoError(t, err, "command execution failed with error: %v, output: %s", err, output)
	assert.Greater(t, len(output), 0, "expected non-empty output, got: %s", output)
}
