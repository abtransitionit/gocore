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
// Description: creates an empty file at the given path if it does not exist.
//
// Parameters:
//
//	filePath: The file path.
//
// Returns:
//
//   - bool:  if a new file was created.
//   - error: if an error occured.
//
// Notes:
// - If the file already exists, this function does nothing and returns successfully.
// - If the path is a directory, this function returns an error.
func Touch(filePath string) (bool, error) {

	// === Input validation ===
	if filePath == "" {
		return false, errorx.New("path cannot be empty")
	}

	// Check if the path already exists.
	exists, err := ExistsPath(filePath)
	if err != nil {
		// Specific error: any other (unexpected) error — handle explicitly
		return false, errorx.Wrap(err, "failed to check if file exists at %s", filePath)
	}

	// The path exists.
	if exists {
		// check if it is a directory.
		isDir, err := ExistsFolder(filePath)

		if err != nil {
			// Specific error: any other (unexpected) error — handle explicitly
			return false, errorx.Wrap(err, "failed to check if path is a folder at %s", filePath)
		}

		if isDir {
			// Specific error: path is a directory — handle explicitly
			return false, errorx.New("path is a directory, not a file at %s", filePath)
		}

		// === File already exists ===
		return false, nil
	}

	// The path does not exist - create it.
	file, err := os.Create(filePath)
	if err != nil {
		// Specific error: permission denied — handle explicitly
		if os.IsPermission(err) {
			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
		}
		// Specific error: any other (unexpected) error — handle explicitly
		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
	}
	defer file.Close()

	// === File created ===
	return true, nil
}
