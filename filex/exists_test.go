/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"
	"path/filepath"
	"testing"
)

// Name: TestExistsFile
//
// Cases:
// - 0 - filepath is not provided
// - 1 - filepath is empty
// - 2 - filepath is not empty
// - 3 - filepath does not exist.
// - 4 - filepath exist.
// - 5 - filepath is a file that exists
// - 6 - filepath is a directory that exists
// The TestExistsFile function now uses a table-driven test approach.
// This makes the test cases more organized and easier to add to in the future.
func TestExistsFile(t *testing.T) {
	// Create a temporary directory for test files and directories.
	// This ensures a clean test environment.
	tempDir := t.TempDir()

	// Setup our test environment:
	// Create a temporary file.
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(tempFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("setup failed: could not create temp file: %v", err)
	}

	// Create a temporary directory.
	tempDirSub := filepath.Join(tempDir, "testdir")
	err = os.Mkdir(tempDirSub, 0755)
	if err != nil {
		t.Fatalf("setup failed: could not create temp directory: %v", err)
	}

	// Define the test cases in a slice of structs.
	// This is the core of the table-driven test pattern.
	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
		caseID  int
	}{
		// Case 0: filepath is not provided (should be handled as empty).
		{
			name:    "Case 0: filepath not provided (empty)",
			path:    "",
			want:    false,
			wantErr: true,
			caseID:  0,
		},
		// Case 1: filepath is empty.
		{
			name:    "Case 1: filepath is empty",
			path:    "",
			want:    false,
			wantErr: true,
			caseID:  1,
		},
		// Case 2: filepath is not empty (tested implicitly by other cases).
		// We'll add a specific case to show it's covered.
		{
			name:    "Case 2: filepath is not empty (valid)",
			path:    tempFile,
			want:    true,
			wantErr: false,
			caseID:  2,
		},
		// Case 3: filepath does not exist.
		{
			name:    "Case 3: filepath does not exist",
			path:    filepath.Join(tempDir, "nonexistent.file"),
			want:    false,
			wantErr: false,
			caseID:  3,
		},
		// Case 4: filepath exists (implicitly covered by 5 and 6, but we'll add a check).
		{
			name:    "Case 4: filepath exists (implicitly covered by case 5)",
			path:    tempFile,
			want:    true,
			wantErr: false,
			caseID:  4,
		},
		// Case 5: filepath is a file that exists.
		{
			name:    "Case 5: filepath is an existing file",
			path:    tempFile,
			want:    true,
			wantErr: false,
			caseID:  5,
		},
		// Case 6: filepath is a directory that exists.
		{
			name:    "Case 6: filepath is an existing directory",
			path:    tempDirSub,
			want:    false,
			wantErr: false,
			caseID:  6,
		},
	}

	// Iterate through the test cases and run a sub-test for each.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := ExistsFile(tc.path)

			// Check for the expected error.
			if (err != nil) != tc.wantErr {
				t.Errorf("ExistsFile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			// If no error was expected, check the return value.
			if err == nil && exists != tc.want {
				t.Errorf("ExistsFile() got = %v, want %v", exists, tc.want)
			}
		})
	}
}
