/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"

	"github.com/abtransitionit/gocore/errorx"
)

// Name: Touch
//
// Description: creates a file at the given path.
//
// Parameters:
//
//	filePath: The file path.
//
// Returns:
//
//   - bool:  if a new file was created.
//   - error: if an error occurred.
//
// Notes:
//
//   - If the path already exists, the function returns a specific error without changing anything.
//   - It does not differentiate between existing files, directories, or symlinks.
func Touch(filePath string) (bool, error) {
	// === Input validation ===
	if filePath == "" {
		return false, errorx.New("path cannot be empty")
	}

	// Attempt to create the file.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		// If the file already exists
		if os.IsExist(err) {
			// handle specific error explicitly: expected outcome
			return false, errorx.New("path already exists: %s", filePath)
		}

		// If there is a permission error
		if os.IsPermission(err) {
			// handle specific error explicitly: unexpected failure
			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
		}

		// handle generic errors explicitly: unexpected  errors
		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
	}

	// Ensure the file is closed.
	defer file.Close()

	// === File created ===
	return true, nil
}
