/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"

	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/logx"
)

// Name: Touch
//
// Description: creates a file at the given path if it does not exist.
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
func Touch(filePath string) (bool, error) {
	// First, check if the file already exists.
	// os.Stat returns information about the file.
	_, err := os.Stat(filePath)
	if err != nil {
		// If the file does not exist, we need to create it.
		if os.IsNotExist(err) {
			file, createErr := os.Create(filePath)
			if createErr != nil {
				// Log a warning if the file cannot be created.
				// This is a critical failure, so we'll use our ErrorWithStack method.
				logx.ErrorWithStack(createErr, "failed to create file at %s", filePath)

				// Return false for the bool and wrap the creation error.
				return false, errorx.Wrap(createErr, "failed to create file at %s", filePath)
			}
			// It is crucial to close the file to release system resources.
			defer file.Close()
			return true, nil
		}
		// If the error is not 'file not exist', it's a different, unexpected error.
		return false, errorx.Wrap(err, "failed to check status of file at %s", filePath)
	}

	// The file already exists, so we do nothing as requested.
	return false, nil
}
