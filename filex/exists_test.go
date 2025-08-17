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
func TestExistsFile(t *testing.T) {
	// Setup isolated test environment
	tempDir := t.TempDir()

	// create inputs for the test : a regular file.
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(tempFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err) // exit immediately if test setup fails
	}

	// create inputs for the test : a folder.
	tempSubDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err) // exit immediately if test setup fails
	}

	// create inputs for the test : a symlink pointing to a regular file.
	tempSymlink := filepath.Join(tempDir, "symlink")
	err = os.Symlink(tempFile, tempSymlink)
	if err != nil {
		t.Fatalf("failed to create temporary symlink: %v", err) // exit immediately if test setup fails
	}

	// create inputs for the test : a symlink pointing to a folder.
	tempSymlinkDir := filepath.Join(tempDir, "symlink-to-dir")
	err = os.Symlink(tempSubDir, tempSymlinkDir)
	if err != nil {
		t.Fatalf("failed to create temporary symlink to dir: %v", err) // exit immediately if test setup fails
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
			name:    "Case 3: non-existent file",
			path:    filepath.Join(tempDir, "nonexistent.file"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "Case 4: existing non-regular file (directory)",
			path:    tempSubDir,
			want:    false,
			wantErr: false,
		},
		{
			name:    "Case 5: symlink to a file",
			path:    tempSymlink,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 6: symlink to a folder",
			path:    tempSymlinkDir,
			want:    true,
			wantErr: false,
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

func TestExistsFolder(t *testing.T) {
	// Setup isolated test environment
	tempDir := t.TempDir()

	// create inputs for the test : a regular file.
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(tempFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err) // exit immediately if test setup fails
	}

	// create inputs for the test : a folder.
	tempSubDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err) // exit immediately if test setup fails
	}

	// create inputs for the test : a symlink pointing to a regular file.
	tempSymlinkFile := filepath.Join(tempDir, "symlink-to-file")
	err = os.Symlink(tempFile, tempSymlinkFile)
	if err != nil {
		t.Fatalf("failed to create temporary symlink to file: %v", err) // exit immediately if test setup fails
	}

	// create inputs for the test : a symlink pointing to a folder.
	tempSymlinkDir := filepath.Join(tempDir, "symlink-to-dir")
	err = os.Symlink(tempSubDir, tempSymlinkDir)
	if err != nil {
		t.Fatalf("failed to create temporary symlink to dir: %v", err) // exit immediately if test setup fails
	}

	// Define test cases.
	tests := []struct {
		name    string // test case name
		path    string // the input
		want    bool   // expected result
		wantErr bool   // expected error
	}{
		{
			name:    "Case 1: nominal: existing folder",
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
			name:    "Case 3: non-existent path",
			path:    filepath.Join(tempDir, "nonexistent.file"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "Case 4: existing regular file",
			path:    tempFile,
			want:    false,
			wantErr: false,
		},
		{
			name:    "Case 5: symlink to a file",
			path:    tempSymlinkFile,
			want:    false,
			wantErr: false,
		},
		{
			name:    "Case 6: symlink to a folder",
			path:    tempSymlinkDir,
			want:    true,
			wantErr: false,
		},
	}

	// Iterate through the test cases
	for _, tc := range tests {
		// Run the function under test with the current test case data
		t.Run(tc.name, func(t *testing.T) {
			obtainedResult, err := ExistsFolder(tc.path)

			// Assertion for expected error state
			if tc.wantErr {
				// We expected an error, so we assert that one was obtained.
				assert.Error(t, err, "Expected to obtain an error, but got nil")
			} else {
				// We didn't expect an error, so we assert that none was obtained.
				assert.NoError(t, err, "Obtained an unexpected error: %v", err)

				// compare Obtained Vs expected - the test fails if the obtainedResult does not match the expected one.
				assert.Equal(t, tc.want, obtainedResult, "Obtained result (%v) did not match expected result (%v)", obtainedResult, tc.want)
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
