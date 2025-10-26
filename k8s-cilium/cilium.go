package cilium

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

func (Cilium) DisplayStatus(local bool, remoteHost string, logger logx.Logger) (string, error) {

	// define cli
	var cmds = []string{
		`cilium status`,
	}
	cli := strings.Join(cmds, " && ")

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, HandleCiliumError)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}
