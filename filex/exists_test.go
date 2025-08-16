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

// Name: TestExistsFile
//
// Cases:
//
//  1. Path is an empty string (invalid input).
//  2. Path is non-empty but does not exist as any OS object.
//  3. Path is non-empty and exists as a regular file.
//  4. Path is non-empty and exists as an OS object that is not a regular file (e.g., directory, symlink).
func TestExistsFile(t *testing.T) {
	// Setup isolated test environment
	tempDir := t.TempDir()

	// - create inputs for the test : a regular file.
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(tempFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err) // exit immediately if test setup fails
	}

	// - create inputs for the test : a folder.
	tempSubDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err) // exit immediately if test setup fails
	}

	// - create inputs for the test : a symlink pointing to a regular file.
	tempSymlink := filepath.Join(tempDir, "symlink")
	err = os.Symlink(tempFile, tempSymlink)
	if err != nil {
		t.Fatalf("failed to create temporary symlink: %v", err) // exit immediately if test setup fails
	}

	// Define test cases.
	tests := []struct {
		name    string // test case name
		path    string // the input
		want    bool   // expected result
		wantErr bool   // expected error
	}{
		{
			name:    "Case 1: nominal: existing regular file",
			path:    tempFile,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 2: empty path string",
			path:    "",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Case 3: non-empty string pointing to a non-existent object",
			path:    filepath.Join(tempDir, "nonexistent.file"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "Case 4a: existing non-regular file (directory)",
			path:    tempSubDir,
			want:    false,
			wantErr: false,
		},
		// {
		// 	name:    "Case 4b: existing non-regular file (symlink)",
		// 	path:    tempSymlink,
		// 	want:    false,
		// 	wantErr: false,
		// },
		{
			name:    "Case 4c: existing non-regular file (symlink)",
			path:    tempSymlink,
			want:    true, // os.Stat will follow the symlink and see it as a file.
			wantErr: false,
		},
	}

	// Iterate through the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// -------------------------------
			// 1. run the fucntion under test with the provided inputs for this test case
			// -------------------------------
			got, err := ExistsFile(tc.path)

			// -------------------------------
			// 2. Assertions
			// -------------------------------

			// 2a. Check if an error is expected
			if tc.wantErr {
				assert.Error(t, err, "expected an error but got none") // An error should occur
			} else {
				assert.NoError(t, err, "unexpected error") // No error should occur
				// 2b. Compare the actual return value with the expected one
				// The returned message should be empty if no error occurred
				assert.Equal(t, tc.want, got, "ExistsFile(%q) returned %v, expected %v", tc.path, got, tc.want)
				// assert.Equal(t, tc.want, exists, "Incorrect boolean result")

			}

		})
	}
}

// Name: TestExistsFolder
func TestExistsFolder(t *testing.T) {
	// Setup isolated test environment.
	tempDir := t.TempDir()

	// Create a temporary regular file.
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(tempFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	// Create a temporary directory.
	tempSubDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}

	// Create a symlink to the directory.
	tempSymlink := filepath.Join(tempDir, "symlink-dir")
	err = os.Symlink(tempSubDir, tempSymlink)
	if err != nil {
		t.Fatalf("failed to create temporary symlink: %v", err)
	}

	// Define test cases in a slice of structs.
	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{
			name:    "case 1: nominal: existing directory",
			path:    tempSubDir,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 2: empty path string",
			path:    "",
			want:    false,
			wantErr: true,
		},
		{
			name:    "non-existent path",
			path:    filepath.Join(tempDir, "nonexistent.dir"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "existing regular file",
			path:    tempFile,
			want:    false,
			wantErr: false,
		},
		{
			name:    "symlink to directory",
			path:    tempSymlink,
			want:    true, // os.Stat will follow the symlink and see it as a directory.
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := ExistsFolder(tc.path)

			if tc.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.Equal(t, tc.want, exists, "Incorrect boolean result")
			}
		})
	}
}

// Name: TestExistsPath
func TestExistsPath(t *testing.T) {
	// Setup isolated test environment
	tempDir := t.TempDir()

	// Create inputs for the test - a regular file
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(tempFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file for test setup: %v", err)
	}

	// Create inputs for the test - a folder
	tempSubDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory for test setup: %v", err)
	}

	// Create inputs for the test - a symlink
	tempSymlink := filepath.Join(tempDir, "symlink-file")
	err = os.Symlink(tempFile, tempSymlink)
	if err != nil {
		t.Fatalf("failed to create temporary symlink for test setup: %v", err)
	}

	// Define test cases.
	tests := []struct {
		name    string // test case name
		path    string // the input
		want    bool   // expected result
		wantErr bool   // expected error
	}{
		{
			name:    "Case 1a: nominal: existing regular file",
			path:    tempFile,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 1b: nominal: existing folder",
			path:    tempSubDir,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 1c: nominal: existing symlink",
			path:    tempSymlink,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 2: empty path string",
			path:    "",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Case 3: non-existent path",
			path:    filepath.Join(tempDir, "nonexistent.file"),
			want:    false,
			wantErr: false,
		},
	}

	// Iterate through the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. run the function under test with the provided inputs for this test case
			exists, err := ExistsPath(tc.path)

			// 2. Assertions
			if tc.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.Equal(t, tc.want, exists, "Incorrect boolean result")
			}
		})
	}
}

// Assertions
// assert.NoError(t, err, "failed to create temporary file") // No error should occur while creating it
// assert.NoError(t, err, "failed to create temporary directory") // No error should occur while creating it
// assert.NoError(t, err, "failed to create temporary symlink") // No error should occur while creating it
