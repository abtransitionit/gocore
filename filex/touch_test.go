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

// Name: TestTouch
func TestTouch(t *testing.T) {
	// Setup isolated test environment
	tempDir := t.TempDir()

	// Create inputs for the test: a regular file
	tempFile := filepath.Join(tempDir, "existing.txt")
	err := os.WriteFile(tempFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	// Create inputs for the test: a folder.
	tempSubDir := filepath.Join(tempDir, "existing-dir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}

	// Create inputs for the test: a symlink pointing to a regular file.
	tempSymlinkFile := filepath.Join(tempDir, "symlink-to-file")
	err = os.Symlink(tempFile, tempSymlinkFile)
	if err != nil {
		t.Fatalf("failed to create temporary symlink to file: %v", err)
	}

	// Create inputs for the test: a symlink pointing to a folder.
	tempSymlinkDir := filepath.Join(tempDir, "symlink-to-dir")
	err = os.Symlink(tempSubDir, tempSymlinkDir)
	if err != nil {
		t.Fatalf("failed to create temporary symlink to dir: %v", err)
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
			wantErr: true,
		},
		{
			name:    "Case 3: touching an empty path string",
			path:    "",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Case 4a: path is a directory",
			path:    tempSubDir,
			want:    false,
			wantErr: true,
		},
		{
			name:    "Case 4b: path is a symlink to a file",
			path:    tempSymlinkFile,
			want:    false,
			wantErr: true,
		},
		{
			name:    "Case 4c: path is a symlink to a directory",
			path:    tempSymlinkDir,
			want:    false,
			wantErr: true,
		},
	}

	// Iterate through the test cases
	for _, tc := range tests {
		// Run the function under test with the current test case data
		t.Run(tc.name, func(t *testing.T) {
			obtainedResult, err := Touch(tc.path)
			expectedResult := tc.want

			// Assertions
			if tc.wantErr {
				if assert.Error(t, err, "expected an error but got none") {
					t.Logf("Expected Error, Obtained Error: %v", err)
				}
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.Equal(t, expectedResult, obtainedResult, "Incorrect boolean result")
				t.Logf("Expected: %v, Obtained: %v", expectedResult, obtainedResult)
			}
		})
	}
}
