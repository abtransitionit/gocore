// file: gocore/cli/gocli.go
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/properties"
	"github.com/abtransitionit/gocore/run"
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

	// Get the targeted OS:type and OS:arch
	goos, err := properties.GetPropertyLocal("ostype")
	if err != nil {
		return fmt.Errorf("failed to get OS type: %w", err)
	}
	goarch, err := properties.GetPropertyLocal("osarch")
	if err != nil {
		return fmt.Errorf("failed to get OS architecture: %w", err)
	}

	// check folders exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", outputDir)
	}
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// build the artifact for the current platform
	l.Infof("Building project '%s' from path: %s into output: %s for %s:%s: ", projectName, projectPath, outputFile, goos, goarch)
	command := fmt.Sprintf("GOOS=%s GOARCH=%s go build -o %s %s", goos, goarch, outputFile, projectPath)
	output, err := run.RunOnLocal(command)
	if err != nil {
		l.Errorf("go build command failed: %v with output: %s", err, output)
		return err
	}

	// build the artifact for linux platform
	goos = "linux"
	goarch = "amd64"
	outputFile = filepath.Join(outputDir, projectName+"-"+goos)
	l.Infof("Building project '%s' from path: %s into output: %s for %s:%s: ", projectName, projectPath, outputFile, goos, goarch)
	command = fmt.Sprintf("GOOS=%s GOARCH=%s go build -o %s %s", goos, goarch, outputFile, projectPath)
	output, err = run.RunOnLocal(command)
	if err != nil {
		l.Errorf("go build command failed: %v with output: %s", err, output)
		return err
	}

	// TODO:
	// success
	l.Infof("Successfully built project to file: %s", outputFile)
	return nil
}
