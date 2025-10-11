package helm

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

// Return: The cli to list the helm charts in a repo
func (chart HelmChart) List(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm search repo %s`, chart.Repo.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// display chart metadata - chart.yaml
// helm show chart $chartName/$RepoName

// display chart metadata - values.yaml
// helm show values $chartName/$RepoName

// var releaseValueShortDesc = "Display user defined values about a relase"
// 	Example: `
// 	xxx kbe-cilicium  kube-system
// cli := fmt.Sprintf(`helm get values %s -n %s`, args[0], args[1])
