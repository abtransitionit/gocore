/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"

	"github.com/abtransitionit/gocore/errorx"
)

// Name: ExistsFile
//
// Description: Checks if the path exists and is a regular file.
//
// Parameters:
//
//	filePath: The file path.
//
// Returns:
//
//   - bool:  True if the path exists and is a regular file.
//   - error: If an error occurred other than "file not found."
//
// Notes:
//
//   - This function follows symbolic links to their target.
func ExistsFile(filePath string) (bool, error) {
	// === Input validation ===
	if filePath == "" {
		return false, errorx.New("path cannot be empty")
	}

	// Use os.Stat to get file info and follow symlinks. This is the most efficient - way to check for both existence and file type with a single system call.
	info, err := os.Stat(filePath)
	if err != nil {
		// If the file already exists
		if os.IsNotExist(err) {
			// handle specific error explicitly: expected outcome : file not found
			return false, nil
		}

		// If there is a permission error
		if os.IsPermission(err) {
			// handle specific error explicitly: unexpected failure
			return false, errorx.Wrap(err, "permission denied accessing %s", filePath)
		}

		// handle generic errors explicitly: unexpected  errors
		return false, errorx.Wrap(err, "failed to get file info for %s", filePath)
	}

	// === Path exists — check that it is a regular file ===
	isFile := info.Mode().IsRegular() // true if regular file, false if other type (e.g., directory, symlinkdir, symlinkfile)
	return isFile, nil

}

// Name: ExistsFolder
//
// Description: Checks if a folder exists at the given path.
//
// Parameters:
//
//	folderPath: The path to the folder.
//
// Returns:
//
//   - bool:  True if the path exists and is a directory.
//   - error: If an error occurred other than "file not found."
//
// Notes:
//
//   - This function follows symbolic links to their target.
func ExistsFolder(folderPath string) (bool, error) {
	// === Input validation ===
	if folderPath == "" {
		return false, errorx.New("path cannot be empty")
	}

	// Use os.Stat to get file info and follow symlinks. This is the most efficient - way to check for both existence and file type with a single system call.
	info, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			// handle specific error explicitly: expected outcome : file not found
			return false, nil
		}

		// If there is a permission error
		if os.IsPermission(err) {
			// handle specific error explicitly: unexpected failure
			return false, errorx.Wrap(err, "permission denied accessing %s", folderPath)
		}

		// handle generic errors explicitly: unexpected  errors
		return false, errorx.Wrap(err, "failed to get file info for %s", folderPath)
	}

	// === Path exists — check that it is a directory ===
	isDir := info.IsDir() // true if directory, false if other type (e.g., directory, symlinkdir, symlinkfile)
	return isDir, nil
}

// Name: ExistsPath
//
// Description: checks if an OS given path exists.
//
// Parameters:
//
//	filePath: The path to the file or folder.
//
// Returns:
//
//	bool: indicating if the path exists and an error if a problem occurred.
//	error: if an error occurred.
//
// Notes::
// - The check doesn't distinguish between a file and a folder.
func ExistsPath(path string) (bool, error) {

	// === Input validation ===
	if path == "" {
		// Expected error: user provided an empty path
		return false, errorx.New("path cannot be empty")
	}

	// === Check if the path exists ===
	_, err := os.Stat(path)
	if err != nil {
		// Specific error: path does not exist — expected behavior
		if os.IsNotExist(err) {
			return false, nil
		}
		// Specific error: permission denied — handle explicitly
		if os.IsPermission(err) {
			return false, errorx.Wrap(err, "permission denied accessing %s", path)
		}
		// Specific error: any other (unexpected) error — handle explicitly
		return false, errorx.Wrap(err, "failed to check status of path at %s", path)
	}

	// === Path exists ===
	return true, nil
}

// Name: IsFilePresent
//
// Description: checks if a file exists at the given path.
//
// Parameters:
//
//	filePath: The path to the file.
//
// Returns:
//
//	bool: true if the file exists, false otherwise.
//	error: if an error occurred.
func IsFilePresent(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		// handle specific error explicitly: expected outcome: The file exists
		return true, nil
	} else if os.IsNotExist(err) {
		// handle specific error explicitly: expected outcome: The file does not exists
		return false, nil
	} else {
		// handle generic errors explicitly: unexpected  errors
		return false, errorx.Wrap(err, "failed to get file info for %s", filePath)
	}
}
