// file golinux/run/run.go
package run

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"

	"github.com/abtransitionit/gocore/errorx"
)

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
	command := fmt.Sprintf(`ssh -o BatchMode=yes -o ConnectTimeout=5 %s "echo '%s' | base64 --decode | sh"`, vmName, cliEncoded)

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
func RunOnVm(vmName, cli string) error {
	cmd := exec.Command("ssh", vmName, cli)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command on VM %s: %v, output: %s", vmName, err, string(output))
	}
	return nil
}

// RunOnLocal executes a CLI command on the local machine and returns the output.
// It is an analogy to RunOnVm, as it captures and returns all output for consistent error reporting.
func RunOnLocal(cli string) (string, error) {
	cmd := exec.Command("sh", "-c", cli)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run command locally: %v, output: %s", err, string(output))
	}

	return string(output), nil
}
