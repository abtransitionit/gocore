package helm

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

func (release HelmRelease) Create(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm install %s %s/%s --atomic --wait --version %s --namespace kube-system -f %s`, release.Name, release.ChartName, release.Repo.Name, release.Version, release.ValueFile),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}
func (release HelmRelease) Delete(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm uninstall %s --namespace kube-system`, release.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}

// List installed repos in the config file
func (release HelmRelease) List(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		"helm list -A", //  list releases in namespace dd
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// helm status mymaria -n db      # release status, notes, and resources
