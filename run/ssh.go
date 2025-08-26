package run

import (
	"fmt"

	"github.com/abtransitionit/gocore/errorx"
)

// Name: IsVmSshConfigured
//
// Description: Checks if a given VM alias is configured in the user's SSH configuration files.
//
// Inputs:
//
// - vmName: string: The VM alias (e.g., "o1u").
//
// Return:
//
// - bool: Returns true if the VM is configured, false otherwise.
// - error: An error if there is an issue reading or parsing the SSH config files.
//
// Notes:
//
// - `ssh -G VmName` prints the resolved configuration (as a set of kvpairs) that the SSH client will use for a given host.
// - comparing the value of the key "hostname" to the value of the key "host" (aka. the VM alias).
// Todo:
//   - Must be adapted to stay in gocore
//   - Do not use this linux speciifc code: grep 'hostname ' | cut -d' ' -f2
func IsVmSshConfigured(vmName string) (bool, error) {
	// define The CLI
	command := fmt.Sprintf("ssh -G %s | grep 'hostname ' | cut -d' ' -f2", vmName)

	// run the CLI
	resolvedHostname, err := RunCliLocal(command)

	if err != nil {
		// handle generic error explicitly: unexpected failure
		return false, errorx.Wrap(err, "failed to get resolved hostname")
	}

	// The logic : compare the resolved hostname to the alias.
	// If they are not the same, the VM is configured. We also check for an empty
	// resolved hostname, which would indicate no match was found.
	isConfigured := resolvedHostname != vmName && resolvedHostname != ""

	return isConfigured, nil
}

// Name: IsVmSshReachable
//
// Description: Checks if a given VM alias is reachable via SSH.
//
// Purpose: This function attempts to establish a non-interactive SSH connection to a VM.
// Inputs:
// - vmName: string: The alias of the VM to check (e.g., "o1u").
// Return:
// - bool: `true` if the VM is reachable, `false` otherwise.
// - error: An error if the underlying command fails for an unexpected reason.
// Notes:
// - The `-o BatchMode=yes` flag prevents interactive prompts and long waits.
// - The `-o ConnectTimeout=5` flag sets a 5-second timeout, so the function never hang indefinitely. Without this, if the remote VM is powered off or unreachable, the SSH command could hang for several minutes before timing out
func IsVmSshReachable(vmName string) (bool, error) {

	// 1: Check if the VM is configured.
	isConfigured, err := IsVmSshConfigured(vmName)
	if err != nil {
		// handle generic error explicitly: unexpected failure
		return false, errorx.Wrap(err, "failed to check vm configuration")
	}
	if !isConfigured {
		// If not configured, we return false without waiting for a timeout.
		return false, nil
	}

	// 2: Now that we know it's configured, check if it's reachable.
	// command := fmt.Sprintf("ssh %s true", vmName)
	command := fmt.Sprintf("ssh -o BatchMode=yes -o ConnectTimeout=5 %s 'exit'", vmName)

	_, err = RunCliLocal(command)
	if err != nil {
		// If RunCliLocal returns an error, it means the SSH connection failed.
		// This is the expected behavior for a non-reachable host.
		return false, nil
	}

	return true, nil
}
