/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Name: Touch
func TestTouch(t *testing.T) {
	// Setup isolated test environment
	tempDir := t.TempDir()

	// create inputs for the test : a regular file.
	tempFile := filepath.Join(tempDir, "existing.txt")
	err := os.WriteFile(tempFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	// create inputs for the test : a folder.
	tempSubDir := filepath.Join(tempDir, "existing-dir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}

	// Define test cases.
	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{
			name:    "Case 1: nominal: touching a non-existent file",
			path:    filepath.Join(tempDir, "new-file.txt"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 2: touching an existing file",
			path:    tempFile,
			want:    false,
			wantErr: false,
		},
		{
			name:    "Case 3: touching an empty path string",
			path:    "",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Case 4: path is a directory",
			path:    tempSubDir,
			want:    false,
			wantErr: true,
		},
	}

	// Iterate through the test cases
	for _, tc := range tests {
		// Run the function under test with the current test case data
		t.Run(tc.name, func(t *testing.T) {
			obtainedResult, err := ExistsFile(tc.path)
			expectedResult := tc.want

			// Assertion for expected error state
			if tc.wantErr {
				// We expected an error, so we assert that one was obtained.
				assert.Error(t, err, "Expected to obtain an error, but got nil")
			} else {
				// We didn't expect an error, so we assert that none was obtained.
				assert.NoError(t, err, "Obtained an unexpected error: %v", err)

				// compare Obtained Vs expected - the test fails if the obtainedResult does not match the expected one.
				assert.Equal(t, expectedResult, obtainedResult, "Obtained result (%v) did not match expected result (%v)", obtainedResult, expectedResult)
			}
		})
	}
}
