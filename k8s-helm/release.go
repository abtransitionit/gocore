package helm

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

// Returns the cli to create a release from a chart into a k8s cluster
func (release HelmRelease) Create(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`
			helm install %s %s/%s --atomic --wait --version %s --namespace kube-system -f %s
			`, release.Name, release.ChartName, release.Repo.Name, release.Version, release.ValueFile),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// fmt.Sprintf(`helm upgrade --install %s %s/%s --version %s --namespace %s -f %s`, chartReleaseFlag, repoNameFlag, chartNameFlag, chartVersionFlag, k8sNamespaceFlag, fileConfFlag)
// fmt.Sprintf(`helm upgrade --install %s %s/%s --version %s --namespace %s`, chartReleaseFlag, repoNameFlag, chartNameFlag, chartVersionFlag, k8sNamespaceFlag)

// Returns the cli to delete a release in a k8s cluster
func (release HelmRelease) Delete(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm uninstall %s --namespace kube-system`, release.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}

// Returns the cli to list the releases installed in a k8s cluster
func (release HelmRelease) List(ctx context.Context, logger logx.Logger) (string, error) {
	var cmds = []string{
		"helm list -A", //  list releases in namespace dd
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// helm status mymaria -n db      # release status, notes, and resources
