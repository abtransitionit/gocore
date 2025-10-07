package helm

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

func (repo HelmRepo) Add(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm repo add %s %s`, repo.Name, repo.Url),
		`helm repo update`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}
func (repo HelmRepo) Delete(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm repo remove %s`, repo.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}

// List installed repos in the config file
func (repo HelmRepo) List(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		`helm repo list`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// List charts for a given helm repo
func (repo HelmRepo) ListChart(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm search repo %s`, repo.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}
