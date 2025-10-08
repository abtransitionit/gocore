// file golinux/run/run.go
package run

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/logx"
)

func RunCliSshLive(vmName, cli string) error {
	// step: check the VM is reachable
	isSshReachable, err := IsVmSshReachable(vmName)
	if err != nil {
		return errorx.Wrap(err, "failed to check VM SSH reachability")
	}
	if !isSshReachable {
		return errorx.New("vm '%s' is not SSH reachable", vmName)
	}

	// step: Base64 encode the input to handle complex quoting and special characters.
	cliEncoded := base64.StdEncoding.EncodeToString([]byte(cli))

	// step: Define the full SSH command to run (same as RunCliSsh).
	command := fmt.Sprintf(
		`ssh -o BatchMode=yes -o ConnectTimeout=5 %s "echo '%s' | base64 --decode | sh"`,
		vmName,
		cliEncoded,
	)

	// step: Use exec.Command so we can stream live output.
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout // stream remote stdout directly
	cmd.Stderr = os.Stderr // stream remote stderr directly

	// step: Run and stream live output
	if err := cmd.Run(); err != nil {
		return errorx.Wrap(err, "failed to run remote command on '%s'", vmName)
	}

	return nil
}

// Name: RunCliSsh
//
// Description: Executes a command on a remote VM via SSH. It first performs a SSH reachability check, then executes a command on the remote host, capturing the output.
//
// Inputs:
//
// - vmName:  string: The alias of the VM to connect to (e.g., "o1u").
// - command: string: The command string to be executed on the remote VM.
//
// Return:
//
// - string: The standard output and standard error from the remote command.
// - error: An error if the VM is not reachable, or if the SSH command fails.
//
// Notes:
//
// - The reachability check is performed using the IsVmSshReachable function.
// - This function is the remote equivalent of RunLocal. .
// - The `-o BatchMode=yes` flag prevents interactive prompts and long waits.
// - The `-o ConnectTimeout=5` flag sets a 5-second timeout, so the function never hang indefinitely. Without this, if the remote VM is powered off or unreachable, the SSH command could hang for several minutes before timing out
// - The remote command is Base64 encoded to avoid issues with complex quotes and special characters.
func RunCliSsh(vmName, cli string) (string, error) {

	// step: check the VM is reachable
	isSshReachable, err := IsVmSshReachable(vmName)
	if err != nil {
		// handle generic error explicitly: unexpected failure
		return "", errorx.Wrap(err, "failed to check VM SSH reachability")
	}
	if !isSshReachable {
		// handle specific error explicitly: expected outcome
		return "", errorx.New("vm '%s' is not SSH reachable", vmName)
	}

	// step: Base64 encode the input to handle complex quoting and special characters.
	cliEncoded := base64.StdEncoding.EncodeToString([]byte(cli))

	// step: Now that the VM is reachable, define the full SSH command to run.
	command := fmt.Sprintf(`ssh -o BatchMode=yes -o ConnectTimeout=5 %s "echo '%s' | base64 --decode | $SHELL -l"`, vmName, cliEncoded)

	// step: Run the command.
	output, err := RunCliLocal(command)

	// manage error
	if err != nil {
		// handle generic error explicitly: unexpected failure
		return output, errorx.Wrap(err, "failed to run remote command on '%s'", vmName)
	}

	// success
	return output, nil
}

// Name: RunCLILocal
// Description: Executes a local command or complex CLI pipeline.
// Inputs:
// - command: string: The complete command string to be executed (e.g., "ls -la" or "ssh -G myhost | grep hostname").
// Return:
//
// - string: The combined standard output and standard error from the command.
// - error: An error if the command fails to run or exits with a non-zero status.
//
// Notes:
//
// - Uses `sh -c` to ensure complex commands with pipes and redirects execute correctly.
// - Captures both standard output and standard error.
// - Trims leading/trailing whitespace from the final output.
func RunCliLocal(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command and wait for it to finish: This will return a non-nil error if the command exits with a non-zero status.
	err := cmd.Run()
	output := strings.TrimSpace(out.String())

	// manage error
	if err != nil {
		// handle generic error explicitly: unexpected failure
		return output, errorx.Wrap(err, "command failed: %s", output)
	}

	// success
	return output, nil
}

// RunOnVm executes a CLI command on a remote VM via SSH
func RunOnVm(vmName, cli string) (string, error) {
	cmd := exec.Command("ssh", vmName, cli)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// This error message now includes the output, which is useful for debugging.
		return "", fmt.Errorf("failed to run command on VM %s: %w, output: %s", vmName, err, string(output))
	}
	return string(output), nil
}

// func RunOnVm(vmName, cli string) error {
// 	cmd := exec.Command("ssh", vmName, cli)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to run command on VM %s: %v, output: %s", vmName, err, string(output))
// 	}
// 	return nil
// }

// RunOnLocal executes a CLI command on the local machine and returns the output.
// It is an analogy to RunOnVm, as it captures and returns all output for consistent error reporting.
func RunOnLocal(cli string) (string, error) {
	cmd := exec.Command("sh", "-c", cli)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to run command locally: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// Name: CustomErrorHandler
//
// Description: This is the signature for functions that check for a "soft" or "handled" error (e.g., specific Helm warnings).
//
// Inputs:
// - err: error: The error to check.
// - logger: *logx.Logger: A logger to use for logging messages.
//
// Return:
// - bool: true if the error was handled and execution should continue (output empty).
type CustomErrorHandler func(err error, logger logx.Logger) bool

// Name: NoOpErrorHandler
//
// Description: A NoOp (No Operation) Handler for commands that don't need special error checking.
//
// Notes:
// - This is critical for making the 'errorHandler' argument optional in practice.
func NoOpErrorHandler(err error, logger logx.Logger) bool {
	return false // Never handle the error; let the main function decide if it's fatal.
}

// Name: RunCliQuery
//
// Description: runs the provided command string (cli) either locally or remotely.
//
// Inputs:
// - cli: string: The command string to be executed (e.g., "ls -la", "kubectl get pods", "helm repo list", "goluc list").
// - logger: *logx.Logger: A logger to use for logging messages.
// - isLocal: bool: Whether to run the command locally or remotely.
// - remoteHost: string: The hostname or IP address of the remote host if running remotely.
// Notes:
//   - It uses isLocal to decide the execution method and remoteHost for SSH connection.
func ExecuteCliQuery(cli string, logger logx.Logger, isLocal bool, remoteHost string, errorHandler CustomErrorHandler) (string, error) {
	var output string
	var err error

	// 1. Determine execution environment and run the command
	if isLocal {
		logger.Debugf("running on local: %s", cli)
		output, err = RunOnLocal(cli)
	} else {
		// Ensure remoteHost is not empty if running remotely (good practice)
		if remoteHost == "" {
			return "", errors.New("remote host cannot be empty when running remotely")
		}
		logger.Debugf("running on remote: %s : %s", remoteHost, cli)
		output, err = RunCliSsh(remoteHost, cli)
	}

	// 2. Handle "errors" that are not true errors
	if errorHandler(err, logger) {
		return "", nil
	}
	// if handleHelmError(err, logger) {
	// 	return "", nil // Handled gracefully
	// }

	// 3. Handle true execution errors
	if err != nil {
		// Return a wrapped error that includes the command run
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}
