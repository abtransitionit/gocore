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
// Description: unit test file
//
// Case:
// - 1 - filepath does not exist.
// - 2 - filepath exist.
// - 3 - filepath is a file that exists
// - 4 - filepath is a directory that exists
func TestExistsFile(t *testing.T) {
	// Create a temporary directory for our test files and folders.
	// This ensures our test is clean and doesn't affect the user's filesystem.
	tempDir := t.TempDir()

	// --- Case 1: file exists ---
	// Create a temporary file inside our temporary directory.
	filePath := filepath.Join(tempDir, "test-file.txt")
	err := os.WriteFile(filePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	exists, err := ExistsFile(filePath)
	if err != nil {
		t.Fatalf("ExistsFile returned an unexpected error for an existing file: %v", err)
	}
	if !exists {
		t.Errorf("ExistsFile should have returned true for an existing file, but got false")
	}

	// --- Case 2: path is a directory ---
	// This is a crucial check for your function's logic.
	dirPath := filepath.Join(tempDir, "test-dir")
	err = os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	exists, err = ExistsFile(dirPath)
	if err != nil {
		t.Fatalf("ExistsFile returned an unexpected error for a directory: %v", err)
	}
	if exists {
		t.Errorf("ExistsFile should have returned false for a directory, but got true")
	}

	// --- Case 3: path does not exist ---
	nonExistentPath := filepath.Join(tempDir, "non-existent-file.txt")
	exists, err = ExistsFile(nonExistentPath)
	if err != nil {
		t.Fatalf("ExistsFile returned an unexpected error for a non-existent path: %v", err)
	}
	if exists {
		t.Errorf("ExistsFile should have returned false for a non-existent path, but got true")
	}

	// --- Case 4: path is empty ---
	// Your code has a specific check for this. We need to test it.
	exists, err = ExistsFile("")
	if err == nil {
		t.Errorf("ExistsFile should have returned an error for an empty path, but it did not")
	}
	if exists {
		t.Errorf("ExistsFile should have returned false for an empty path, but got true")
	}
}
