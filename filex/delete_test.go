/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"
	"path/filepath"
	"testing"
)

// Name: TestDeleteFile
func TestDeleteFile(t *testing.T) {
	// Setup isolated test environment
	tempDir := t.TempDir()

	// create inputs for the test : a regular file.
	tempFile := filepath.Join(tempDir, "existing.txt")
	err := os.WriteFile(tempFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	// create inputs for the test : an empty folder.
	tempSubDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(tempSubDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err) // exit immediately if test setup fails
	}

	// // create inputs for the test : a folder with a file.
	tempSubDirWithFile := filepath.Join(tempDir, "subdir-with-file")
	err = os.Mkdir(tempSubDirWithFile, 0755)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempSubDirWithFile, "file.txt"), []byte("data"), 0644)
	if err != nil {
		t.Fatalf("failed to create file inside temporary directory: %v", err)
	}

	// Define test cases.
	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{
			name:    "Case 1: nominal: deleting an existing file",
			path:    tempFile,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 2: deleting a non-existent file",
			path:    filepath.Join(tempDir, "non-existent.txt"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 3: deleting an empty path string",
			path:    "",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Case 4: deleting an empty directory",
			path:    tempSubDir,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Case 4: deleting a non-empty directory",
			path:    tempSubDirWithFile,
			want:    false,
			wantErr: true,
		},
	}

	// Iterate through the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obtainedResult, err := DeleteFile(tc.path)
			expectedResult := tc.want

			// Assertion for expected error state
			if tc.wantErr {
				// We expected an error, so we assert that one was obtained.
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				// We didn't expect an error, so we assert that none was obtained.
				if err != nil {
					t.Errorf("Obtained an unexpected error: %v", err)
				}

				// compare Obtained Vs expected
				if obtainedResult != expectedResult {
					t.Errorf("Obtained result (%v) did not match expected result (%v)", obtainedResult, expectedResult)
				}
			}
		})
	}
}
