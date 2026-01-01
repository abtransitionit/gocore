package filex

import (
	"fmt"
	"os"
)

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
func CreateFileFromString(filePath, content string, failIfExists bool) (bool, error) {
	// check - Input validation
	if filePath == "" {
		return false, fmt.Errorf("path cannot be empty")
	}

	// Check - if the path is a directory
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		return false, fmt.Errorf("path is a directory not a file: %s", filePath)
	}

	// Handle failIfExists input
	createFlags := os.O_CREATE | os.O_WRONLY // fail if it already exists
	if failIfExists {
		createFlags |= os.O_EXCL // do not fail if it already exists
	}
	// Attempt to create the file (fail if it already exists)
	file, err := os.OpenFile(filePath, createFlags, 0666)
	if err != nil {
		if os.IsExist(err) {
			if failIfExists {
				return false, fmt.Errorf("file already exists: %s", filePath)
			}
			// failIfExists == false → existing file is OK → not an error → just return false
			return false, nil
		}
		if os.IsPermission(err) {
			return false, fmt.Errorf("creating the file: permission denied to create file at %s > %w", filePath, err)
		}
		return false, fmt.Errorf("creatingf file: %s >  %w", filePath, err)
	}
	defer file.Close()

	// Write content only if not empty
	if content != "" {
		if _, err := file.WriteString(content); err != nil {
			return false, fmt.Errorf("writing content to file %s: %w", filePath, err)
		}
	}
	return true, nil
}
