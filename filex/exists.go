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
//   - bool:  indicating if the file exists.
//   - error: if a problem occurred.
//
// ExistsFile checks if a path exists and is a regular file
func ExistsFile(filePath string) (bool, error) {

	// === Input validation and path existence check ===
	exists, err := ExistsPath(filePath)
	if err != nil || !exists {
		return false, err
	}

	// === Path exists — get file info ===
	info, err := os.Stat(filePath)
	if err != nil {
		// Specific error: permission denied — handle explicitly
		if os.IsPermission(err) {
			return false, errorx.Wrap(err, "permission denied accessing %s", filePath)
		}
		// Specific error: any other (unexpected) error — handle explicitly
		return false, errorx.Wrap(err, "failed to get info of path at %s", filePath)
	}

	// === Path exists — check that it is a regular file ===
	isFile := info.Mode().IsRegular() // true if regular file, false if other type (e.g., directory, symlink)
	return isFile, nil
}

// Name: ExistsFolder
//
// Description: checks if a folder exists at the given path.
//
// Parameters:
//
//	folderPath: The path to the folder.
//
// Returns:
//
//   - bool:  indicating if the folder exists.
//   - error: if a problem occurred.
func ExistsFolder(folderPath string) (bool, error) {

	// === Input validation and path existence check ===
	exists, err := ExistsPath(folderPath)
	if err != nil || !exists {
		return false, err
	}

	// === Path exists — get file info ===
	info, err := os.Stat(folderPath)
	if err != nil {
		// Specific error: permission denied — handle explicitly
		if os.IsPermission(err) {
			return false, errorx.Wrap(err, "permission denied accessing %s", folderPath)
		}
		// Specific error: any other (unexpected) error — handle explicitly
		return false, errorx.Wrap(err, "failed to get info of path at %s", folderPath)
	}

	// === Path exists — check that it is a directory ===
	isDir := info.IsDir() // true if directory, false if other type (e.g., file, symlink)
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
//	A boolean indicating if the path exists and an error if a problem occurred.
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
