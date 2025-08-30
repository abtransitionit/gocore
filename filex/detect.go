/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/abtransitionit/gocore/run"
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
//   - Currently reads the magic bytes of the file to determine its type.
//   - Currently supports ".tar.gz", ".zip", and defaults to "binary".
func DetectBinaryType(ctx context.Context, vmName string, filePath string) (string, error) {

	// Remote detect
	if vmName != "" {
		cmd := fmt.Sprintf("goluc do detect %s", filePath)
		fileName, err := run.RunCliSsh(vmName, cmd)
		if err != nil {
			return "", fmt.Errorf("failed to remote detect file typr of '%s' from '%s': %w", filePath, vmName, err)
		}
		return strings.TrimSpace(fileName), nil
	}

	// Local detect
	if filePath == "" {
		return "", fmt.Errorf("empty file path")
	}

	// Open file
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file %s: %w", filePath, err)
	}
	defer f.Close()

	// Read first 4 bytes
	magic := make([]byte, 4)
	n, err := f.Read(magic)
	if err != nil {
		return "", fmt.Errorf("cannot read magic bytes from %s: %w", filePath, err)
	}
	if n < 2 {
		return "binary", nil // too short to determine, assume binary
	}

	// Check for gzip (tar.gz)
	if magic[0] == 0x1F && magic[1] == 0x8B {
		return "tgz", nil
	}

	// Check for zip
	if n >= 4 && magic[0] == 0x50 && magic[1] == 0x4B && magic[2] == 0x03 && magic[3] == 0x04 {
		return "zip", nil
	}

	// Default: binary
	return "exe", nil
}
