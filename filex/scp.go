/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

// Name: Scp
//
// Description: Executes an scp command to copy files between local and remote machines.
//
// Parameters:
//
//	l: The logger to use for command output.
//	source: The source path (local or remote).
//	destination: The destination path (local or remote).
//
// Returns:
//
//	error: An error if the scp command fails.
func Scp(l logx.Logger, source string, destination string) error {
	l.Infof("Initiating SCP transfer from %s to %s", source, destination)

	// Construct the scp command string.
	// We use the -r flag for recursive copy.
	// Note: Be cautious with user-provided input to prevent command injection.
	command := fmt.Sprintf("scp -r %s %s", source, destination)

	// Execute the command using the helper function.
	output, err := run.RunOnLocal(command)
	if err != nil {
		// Log the captured output for debugging purposes.
		l.Error(strings.TrimSpace(output))
		return fmt.Errorf("scp command failed: %w", err)
	}

	l.Info("SCP transfer completed successfully.")
	return nil
}
