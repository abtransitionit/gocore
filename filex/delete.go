/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"

	"github.com/abtransitionit/gocore/errorx"
)

// Name: DeleteFile
//
// Description: Deletes the file at the specified path.
//
// Returns:
//
//   - bool: `true` if the file was successfully deleted or did not exist.
//   - error: Returns an error if an unexpected issue occurs during deletion.
func DeleteFile(filePath string) (bool, error) {
	// Input validation.
	if filePath == "" {
		return false, errorx.New("path cannot be empty")
	}

	// First, check if the file exists at all.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// handle specific error explicitly: expected outcome: The file is already gone
		return true, nil
	}

	// Now that we know it exists, remove it.
	err := os.Remove(filePath)
	if err != nil {
		// If there is a permission error.
		if os.IsPermission(err) {
			// handle specific error explicitly: unexpected failure
			return false, errorx.Wrap(err, "permission denied to delete file at %s", filePath)
		}

		// A generic catch-all for all other unexpected errors.
		return false, errorx.Wrap(err, "failed to delete file at %s", filePath)
	}

	// success
	return true, nil
}
