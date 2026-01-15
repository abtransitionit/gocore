package kubectl

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

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
