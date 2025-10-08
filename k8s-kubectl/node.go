package kubectl

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

// Name: listNodes
//
// Description: Performs the logic to list nodes
//
// Returns:
// - string: the raw output string
// - error: any error that occurred during the execution.
func ListNode(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// define cli
	cli, err := Resource{Type: "node"}.List()
	if err != nil {
		return "", fmt.Errorf("failed to build kubectl list command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}

func ListNs(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// define cli
	cli, err := Resource{Type: "ns"}.List()
	if err != nil {
		return "", fmt.Errorf("failed to build kubectl list command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}

func ListPod(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// define cli
	cli, err := Resource{Type: "pod"}.List()
	if err != nil {
		return "", fmt.Errorf("failed to build kubectl list command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}
