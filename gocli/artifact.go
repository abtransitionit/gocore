package gocli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/filex"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/property"
	"github.com/abtransitionit/gocore/run"
)

// Name: BuildProject
//
// Description: builds a Go project for the current platform.
//
// Parameters:
//
//	projectPath: The full path to the Go project to build (e.g., "/path/to/projectFolder").
//	outputDir: The directory to save the built artifact (e.g., "/tmp").
//
// Returns:
//
//	error: An error if the build fails.
//
// Todos:
//
//   - rewrite that function that builds a Go project for the current and linux platforms.
//   - spit it in 3 functions: BuildProject, BuildProjectForCurrentPlatform, BuildProjectForPlatform
//   - a platfrom is an os type and an os architecture
func BuildGoProject(logger logx.Logger, projectPath, outputDir string) error {

	// check parameters
	if projectPath == "" {
		return fmt.Errorf("project path is empty")
	}
	if outputDir == "" {
		return fmt.Errorf("output directory is empty")
	}

	// check folders exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", outputDir)
	}
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// set the artifact name and the full artifact path from parameters
	projectName := filepath.Base(projectPath)
	outputFile := filepath.Join(outputDir, projectName)

	// set the targeted OS:type and OS:arch dynamically
	goos, err := property.GetProperty("ostype")
	if err != nil {
		return fmt.Errorf("failed to get OS type: %w", err)
	}
	goarch, err := property.GetProperty("osarch")
	if err != nil {
		return fmt.Errorf("failed to get OS architecture: %w", err)
	}

	// build the artifact for the current platform
	logger.Infof("Building artifact '%s' from project: %s for platform:%s/%s: ", outputFile, projectPath, goos, goarch)
	command := fmt.Sprintf("GOOS=%s GOARCH=%s go build -o %s %s", goos, goarch, outputFile, projectPath)
	output, err := run.RunOnLocal(command)
	if err != nil {
		logger.Errorf("go build command failed: %v with output: %s", err, output)
		return err
	}

	// build the artifact for linux/amd64 platform
	goos = "linux"
	goarch = "amd64"
	outputFile = filepath.Join(outputDir, projectName+"-"+goos)
	logger.Infof("Building artifact '%s' from project: %s for platform:%s/%s: ", outputFile, projectPath, goos, goarch)
	command = fmt.Sprintf("GOOS=%s GOARCH=%s go build -o %s %s", goos, goarch, outputFile, projectPath)
	output, err = run.RunOnLocal(command)
	if err != nil {
		logger.Errorf("go build command failed: %v with output: %s", err, output)
		return err
	}

	// success
	logger.Infof("Successfully built project to file: %s", outputFile)
	return nil
}

// Name: DeployGoArtifact
//
// Description: Deploys a Go artifact to a rlocal or remote destination without sudo privileges.
//
// Parameters:
//
//	l: The logger to use.
//	artifactPath: The full path to the local executable to be deployed (e.g., "/tmp/goluc-linux").
//	remoteDestination: The remote destination path, including host and file (e.g., "o1u:/usr/local/bin/goluc").
//
// Returns:
//
//	bool: true if the deployment was successful.
//	error: An error if the deployment failed.
func DeployGoArtifact(logger logx.Logger, artifactPath, remoteDestination string) (bool, error) {
	// check parameters
	if artifactPath == "" {
		return false, fmt.Errorf("artifact path is empty")
	}
	if remoteDestination == "" {
		return false, fmt.Errorf("remote destination is empty")
	}

	// check if the artifact exists.
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return false, errorx.Wrap(err, "local artifact not found at path: %s", artifactPath)
	}

	// scp the artifact to a non root location on the remote machine - if dst path is root use filex.ScpAsSudo works
	err := filex.Scp(logger, artifactPath, remoteDestination)
	if err != nil {
		logger.Errorf("%v", err)
		return false, err
	}

	// sucess
	// logger.Debugf("Successfully deployed artifact to %s", remoteDestination)
	return true, nil
}
