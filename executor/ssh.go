package executor

import (
	"fmt"

	"github.com/abtransitionit/gocore/errorx"
)

// Name: IsSSHConfigured
// Description: Checks if a given VM alias is configured in the user's SSH configuration files.
// Purpose: validate that a VM alias exists before attempting an SSH connection.
// Inputs:
// - vm: string: The VM alias (e.g., "o1u").
// Return:
// - bool: Returns true if the VM is configured, false otherwise.
// - error: An error if there is an issue reading or parsing the SSH config files.
// Notes:
// - `ssh -G VmName` prints the resolved configuration (as a set of kvpairs) that the SSH client will use for a given host.
// - comparing the value of the key "hostname" to the value of the key "host" (aka. the VM alias).

func IsVmSshConfigured(vmName string) (bool, error) {
	// define The CLI
	command := fmt.Sprintf("ssh -G %s | grep 'hostname ' | cut -d' ' -f2", vmName)

	// run the CLI
	resolvedHostname, err := RunLocal(command)

	if err != nil {
		// An error here could indicate a problem with ssh itself, or that the `grep` command
		// found no match (in which case the output will be empty, which is fine).
		// We can return false and the error to the caller for proper handling.
		// return false, errorx.WithStack(fmt.Errorf("failed to run ssh command: %w", err))
		return false, errorx.WithStack(fmt.Errorf("failed to get resolved hostname: %w", err))
	}

	// The logic : compare the resolved hostname to the alias.
	// If they are not the same, the VM is configured. We also check for an empty
	// resolved hostname, which would indicate no match was found.
	isConfigured := resolvedHostname != vmName && resolvedHostname != ""

	return isConfigured, nil
}

// Name: IsVmSshReachable
// Description: Checks if a VM is currently reachable via SSH.
// Purpose: This function attempts to establish a non-interactive SSH connection to a VM.
// A successful connection indicates the VM is online and the SSH daemon is running.
// A failure indicates the VM is either not running, or there is a network issue
// preventing a connection.
// Inputs:
// - vmName: string: The alias of the VM to check (e.g., "o1u").
// Return:
// - bool: `true` if the VM is reachable, `false` otherwise.
// - error: An error if the underlying command fails for an unexpected reason.
// Notes:
// - The `-o BatchMode=yes` flag prevents interactive prompts and long waits.
// - The `-o ConnectTimeout=5` flag sets a 5-second timeout, so the function never hang indefinitely. Without this, if the remote VM is powered off or unreachable, the SSH command could hang for several minutes before timing out
func IsVmSshReachable(vmName string) (bool, error) {

	// Step 1: Check if the VM is configured.
	isConfigured, err := IsVmSshConfigured(vmName)
	if err != nil {
		return false, errorx.WithStack(fmt.Errorf("failed to check vm configuration: %w", err))
	}
	if !isConfigured {
		// If not configured, we can return immediately without waiting for a timeout.
		return false, nil
	}

	// Step 2: Now that we know it's configured, check if it's reachable.
	// command := fmt.Sprintf("ssh %s true", vmName)
	command := fmt.Sprintf("ssh -o BatchMode=yes -o ConnectTimeout=5 %s 'exit'", vmName)

	_, err = RunCliLocal(command)
	if err != nil {
		return false, nil
	}

	return true, nil
}
