package executor

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/logx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Name: RunCliSsh
// Description: Executes a command on a remote VM via SSH.
// Purpose: This function is the remote equivalent of RunLocal. It first performs a
// reachability check, then executes a command on the remote host, capturing the output.
// Inputs:
// - vmName:  string: The alias of the VM to connect to (e.g., "o1u").
// - command: string: The command string to be executed on the remote VM.
// Return:
// - string: The standard output and standard error from the remote command.
// - error: An error if the VM is not reachable, or if the SSH command fails.
// Notes:
// - The `-o BatchMode=yes` flag prevents interactive prompts and long waits.
// - The `-o ConnectTimeout=5` flag sets a 5-second timeout, so the function never hang indefinitely. Without this, if the remote VM is powered off or unreachable, the SSH command could hang for several minutes before timing out
func RunCliSsh(vmName, cli string) (string, error) {

	// Step: check the VM is reachable
	isReachable, err := IsVmSshReachable(vmName)
	if err != nil {
		return "", errors.Wrap(err, "failed to check VM reachability")
	}
	if !isReachable {
		return "", fmt.Errorf("vm '%s' is not reachable", vmName)
	}

	// Step: Now the VM it's reachable, define the CLI to run
	command := fmt.Sprintf("ssh -o BatchMode=yes -o ConnectTimeout=5 %s '%s'", vmName, cli)

	// Step: Run the CLI
	output, err := RunCliLocal(command)

	// manage error
	if err != nil {
		return output, errors.Wrap(err, fmt.Sprintf("failed to run remote command on '%s'", vmName))
	}

	// success
	return output, nil
}

// Name: RunCliSsh
// Description: Executes a command on a remote VM via SSH.
// Purpose: It is the remote equivalent of RunCliLocal. It captures the output of the remote
// command and handles potential errors, such as a lost connection or a command
// that fails on the remote host.
// Inputs:
// - vmName: string: The alias of the VM to connect to (e.g., "o1u").
// - command: string: The command string to be executed on the remote VM.
// Return:
//   - string: The standard output and standard error from the remote command.
//   - error: An error if the SSH command fails to run or the remote command
//     exits with a non-zero status.
//
// Notes:
// - The `-o BatchMode=yes` flag prevents interactive prompts and long waits.
// - The `-o ConnectTimeout=5` flag sets a 5-second timeout, so the function never hang indefinitely. Without this, if the remote VM is powered off or unreachable, the SSH command could hang for several minutes before timing out
func RunCliSsh2(vmName, command string) (string, error) {
	// This includes safety flags and the command to be run on the remote machine.
	fullCommand := fmt.Sprintf("ssh -o BatchMode=yes -o ConnectTimeout=5 %s '%s'", vmName, command)

	// This ensures that all the logic for running local commands, capturing
	// output, and error handling is reused.
	output, err := RunLocal(fullCommand)
	if err != nil {
		return output, errorx.WithStack(fmt.Errorf("failed to run remote command: %w", err))
	}

	return output, nil
}

// Name: RunSSH
// Description: Executes a command on a remote host via SSH.
// Purpose: It connects to a remote host and executes a single command. It relies
// on the user's SSH configuration files for host details and an active
// SSH agent for authentication.
// Inputs:
// - host: string: The VM alias (e.g., "o1u") as defined in the user's SSH config.
// - command: string: The command string to be executed on the remote host.
// Return:
// - string: The combined standard output and standard error from the command.
// - error: An error if the connection fails or the remote command exits with a non-zero status.
// Notes:
// This function first checks if the host is configured in the SSH config files.
// It also uses `ssh.InsecureIgnoreHostKey()` for simplicity. In a production
// environment, you should use `knownhosts.New()` to verify the host's identity
// and protect against man-in-the-middle attacks.
func RunCliSsh3(host string, command string) (string, error) {
	logx.Init()
	logx.Info("Executing command on %s: '%s'", host, command)

	if configured, err := IsSSHConfigured(host); err != nil {
		return "", errorx.WithStack(fmt.Errorf("failed to check SSH config for host '%s': %w", host, err))
	} else if !configured {
		return "", errorx.WithStack(fmt.Errorf("host '%s' is not configured in your SSH config files", host))
	}

	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		return "", errorx.WithStack(fmt.Errorf("SSH_AUTH_SOCK environment variable not set. Is ssh-agent running?"))
	}

	sshAgent, err := net.Dial("unix", authSock)
	if err != nil {
		return "", errorx.WithStack(fmt.Errorf("failed to connect to SSH agent: %w", err))
	}
	defer sshAgent.Close()

	config := &ssh.ClientConfig{
		User:            "user_from_config_file",
		Auth:            []ssh.AuthMethod{ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	session, err := client.NewSession()
	if err != nil {
		return "", errorx.WithStack(fmt.Errorf("failed to create session: %w", err))
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b

	if err := session.Run(command); err != nil {
		return b.String(), errorx.WithStack(fmt.Errorf("command failed with error: %w, output: %s", err, b.String()))
	}

	logx.Info("Command on %s executed successfully. Output: %s", host, b.String())
	return b.String(), nil
}

// Name: RunCLILocal
// Description: Executes a local command or complex CLI pipeline.
// Purpose: Provides a portable and safe way to execute local commands. It captures
// both standard output and standard error, returning them as a single string.
// This function is the local counterpart to RunSSH.
// Inputs:
// - command: string: The complete command string to be executed (e.g., "ls -la" or "ssh -G myhost | grep hostname").
// Return:
// - string: The combined standard output and standard error from the command, with leading/trailing whitespace removed.
// - error: An error if the command fails to run or exits with a non-zero status.
// ...
// Notes:
// - Uses `sh -c` to ensure complex commands with pipes and redirects execute correctly.
// - Captures both standard output and standard error.
// - Trims leading/trailing whitespace from the final output.
func RunCliLocal(command string) (string, error) {
	// cmd := exec.Command("bash", "-c", command)
	cmd := exec.Command("sh", "-c", command)

	// config: capture both standard output and standard error into a single buffer.
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command and wait for it to finish: This will return a non-nil error if the command exits with a non-zero status.
	err := cmd.Run()
	output := strings.TrimSpace(out.String())

	// manage error
	if err != nil {
		return output, errorx.WithStack(fmt.Errorf("command failed: %w, output: %s", err, out.String()))
	}

	// success
	return output, nil
}
