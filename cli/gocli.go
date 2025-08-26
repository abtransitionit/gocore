// file: gocore/cli/gocli.go
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abtransitionit/gocore/logx"
)

// Name: BuildProject
//
// Description: builds a Go project for the current platform.
//
// Parameters:
//
//	projectPath: The path to the Go project to build.
//	outputDir: The directory to save the built artifact.
//
// Returns:
//
//	error: An error if the build fails.
func BuildGoProject(l logx.Logger, projectPath, outputDir string) error {

	// check parameters
	if projectPath == "" {
		return fmt.Errorf("project path is empty")
	}
	if outputDir == "" {
		return fmt.Errorf("output directory is empty")
	}

	// Get the name of the project from the path.
	projectName := filepath.Base(projectPath)
	outputFile := filepath.Join(outputDir, projectName)

	// check folder exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", outputDir)
	}
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	l.Infof("Building project '%s' from path: %s into output: %s", projectName, projectPath, outputFile)
	// TODO:
	// success
	l.Infof("Successfully built project to: %s", outputFile)
	return nil
}
