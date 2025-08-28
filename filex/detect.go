/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Name: detect
//
// Description: detcts the type of a file (eg. .tar.gz, .zip, txt, binary)
//
// Parameters:
//
//	filePath: The path to the file.
//
// Returns:
//   - string: The type of the file.
//   - error: If an error occurred.
//
// Notes:
//
//   - Currently inspects the file path and returns the type of archive or "binary".
//   - Currently supports ".tar.gz", ".zip", and defaults to "binary".
func Detect(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("empty file path")
	}

	// get the file extension
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".gz":
		// check if the file ends with .tar.gz
		if strings.HasSuffix(strings.ToLower(filePath), ".tar.gz") {
			return "tar.gz", nil
		}
		return "gz", nil
	case ".zip":
		return "zip", nil
	default:
		return "binary", nil
	}
}
