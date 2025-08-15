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
// Description: checks if a file exists at the given path.
//
// Parameters:
//
//	filePath: The path to the file.
//
// Returns:
//
//	A boolean indicating if a file exists and an error if a problem occurred.
func ExistsFile(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errorx.Wrap(err, "failed to check status of file at %s", filePath)
	}

	// The path exists, now check if it's a regular file.
	return info.Mode().IsRegular(), nil
}

// Name: ExistsFolder
// Description: checks if a folder exists at the given path.
//
// Parameters:
//
//	folderPath: The path to the folder.
//
// Returns:
//
//	A boolean indicating if a folder exists and an error if a problem occurred.
func ExistsFolder(folderPath string) (bool, error) {
	info, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errorx.Wrap(err, "failed to check status of folder at %s", folderPath)
	}

	// The path exists, now check if it's a directory.
	return info.IsDir(), nil
}

// Name: ExistsPath
// Description: checks if an OS given path exists.
//
// Parameters:
//
//	filePath: The path to the file or folder.
//
// Returns:
//
//	A boolean indicating if the path exists and an error if a problem occurred.
//
// Notes::
// - The check doesn't distinguish between a file and a folder.
func ExistsPath(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errorx.Wrap(err, "failed to check status of path at %s", filePath)
	}

	return true, nil
}

// Name: FileOrFolderExists
//
// Description: checks for the existence of a file or directory at the given path.
//
// Parameters:
//
//	filePath: The path to the file or directory.
//
// Returns:
// - bool:  the file exists
// - error: a problem occurred
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, which is not an error for this function.
			return false, nil
		}
		// A different error occurred (e.g., permissions).
		// We wrap it to add a stack trace.
		return false, errorx.Wrap(err, "failed to check status of file at %s", filePath)
	}

	// The file or directory exists.
	return true, nil
}
