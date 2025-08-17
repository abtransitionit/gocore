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
	// Input validation
	if filePath == "" {
		return false, errorx.New("path cannot be empty")
	}

	// First, check if the path is a directory.
	// This is a common and important check for file creation functions.
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		return false, errorx.New("path is a directory: %s", filePath)
	}

	// Attempt to create the file.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		// If the file already exists, it is not an error
		if os.IsExist(err) {
			return false, nil
		}

		// If there is a permission error
		if os.IsPermission(err) {
			// handle specific error explicitly: unexpected failure
			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
		}

		// A generic catch-all for all other unexpected errors.
		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
	}

	// Ensure the file is closed.
	defer file.Close()

	// File was successfully created.
	return true, nil
}

// func Touch(filePath string) (bool, error) {
// 	// Input validation
// 	if filePath == "" {
// 		return false, errorx.New("path cannot be empty")
// 	}

// 	// Attempt to create the file.
// 	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0666)
// 	if err != nil {
// 		// If the file already exists, which is not an error in this context.
// 		if os.IsExist(err) {
// 			return false, nil
// 		}

// 		// Explicitly handle the permission denied error, as it's important.
// 		if os.IsPermission(err) {
// 			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
// 		}

// 		// A generic catch-all for all other unexpected errors, including
// 		// "path is a directory" or other file system issues.
// 		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
// 	}

// 	// Ensure the file is closed.
// 	defer file.Close()

// 	// File was successfully created.
// 	return true, nil
// }

// func Touch(filePath string) (bool, error) {
// 	// Input validation
// 	if filePath == "" {
// 		return false, errorx.New("path cannot be empty")
// 	}

// 	// Attempt to create the file.
// 	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0666)
// 	if err != nil {
// 		// If the file already exists
// 		if os.IsExist(err) {
// 			// it is not an error
// 			return false, nil
// 		}

// 		// if the path is a directory.
// 		if info, statErr := os.Stat(filePath); statErr == nil && info.IsDir() {
// 			// handle specific error explicitly: unexpected failure
// 			return false, errorx.New("path is a directory: %s", filePath)
// 		}

// 		// If there is a permission error
// 		if os.IsPermission(err) {
// 			// handle specific error explicitly: unexpected failure
// 			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
// 		}

// 		// handle specific error explicitly: unexpected failure
// 		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
// 	}

// 	// Ensure the file is closed.
// 	defer file.Close()

// 	// === File created ===
// 	return true, nil
// }

// func Touch(filePath string) (bool, error) {
// 	// Input validation
// 	if filePath == "" {
// 		return false, errorx.New("path cannot be empty")
// 	}

// 	// Attempt to create the file.
// 	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0666)
// 	if err != nil {
// 		// If the file already exists
// 		if os.IsExist(err) {
// 			// it is not an error
// 			return false, nil
// 		}

// 		// If there is a permission error
// 		if os.IsPermission(err) {
// 			// handle specific error explicitly: unexpected failure
// 			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
// 		}

// 		// handle specific error explicitly: unexpected failure
// 		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
// 	}

// 	// Ensure the file is closed.
// 	defer file.Close()

// 	// === File created ===
// 	return true, nil
// }

// func Touch(filePath string) (bool, error) {
// 	// === Input validation ===
// 	if filePath == "" {
// 		return false, errorx.New("path cannot be empty")
// 	}

// 	// Attempt to create the file.
// 	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0666)
// 	if err != nil {
// 		// If the file already exists
// 		if os.IsExist(err) {
// 			// it is not an error
// 			return false, nil
// 		}

// 		// if the path is a directory.
// 		if info, statErr := os.Stat(filePath); statErr == nil && info.IsDir() {
// 			// handle specific error explicitly: unexpected failure
// 			return false, errorx.Wrap(err, "path is a directory: %s", filePath)
// 		}

// 		// If there is a permission error
// 		if os.IsPermission(err) {
// 			// handle specific error explicitly: unexpected failure
// 			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
// 		}

// 		// handle specific error explicitly: unexpected failure
// 		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
// 	}

// 	// Ensure the file is closed.
// 	defer file.Close()

// 	// === File created ===
// 	return true, nil
// }
