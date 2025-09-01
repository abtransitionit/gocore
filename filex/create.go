/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"os"

	"github.com/abtransitionit/gocore/errorx"
)

// Name: CreateFileFromString
//
// Description: creates a file at the given path with the specified content.
//
// Parameters:
//
//	filePath: The file path.
//	content: The content to write to the file.
//
// Returns:
//
//   - bool:  if a new file was created.
//   - error: if an error occurred.
//
// Example Usage:
//
//	created, err := CreateFileFromString("/tmp/hello.txt", "Hello world!")
//	if err != nil {
//	  log.Fatal(err)
//	}
//	if created {
//	  fmt.Println("File created")
//	} else {
//	  fmt.Println("File already existed")
//	}
//
// Notes:
//
//   - If the path already exists, the function returns a specific error without changing anything.
//   - It does not differentiate between existing files, directories, or symlinks.
func CreateFileFromString(filePath, content string) (bool, error) {
	// Input validation
	if filePath == "" {
		return false, errorx.New("path cannot be empty")
	}

	// Check if the path is a directory
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		return false, errorx.New("path is a directory: %s", filePath)
	}

	// Attempt to create the file (fail if it already exists)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			// File already exists → not an error, just return false
			return false, nil
		}
		if os.IsPermission(err) {
			return false, errorx.Wrap(err, "permission denied to create file at %s", filePath)
		}
		return false, errorx.Wrap(err, "failed to create file at %s", filePath)
	}
	defer file.Close()

	// Write content
	if _, err := file.WriteString(content); err != nil {
		return false, errorx.Wrap(err, "failed to write content to file at %s", filePath)
	}

	return true, nil
}
