package helm

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

func (release HelmRelease) valueFlag() string {
	if release.ValueFile != "" {
		return fmt.Sprintf("-f %s", release.ValueFile)
	}
	return ""
}

// Returns the cli to create a release from a chart into a k8s cluster
func (release HelmRelease) create() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`
			helm install %s %s --atomic --wait --version %s --namespace %s %s
			`,
			release.Name,
			release.Chart.FullName,
			release.Chart.Version,
			release.Namespace,
			release.valueFlag()),
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

// create a helm release into a kubernetes cluster
func ListRelease(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// Check parameters

	// define cli
	var cmds = []string{
		`helm list`,
	}
	cli := strings.Join(cmds, " && ")

	// // play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil

}

func CreateRelease(local bool, remoteHost string, release HelmRelease, logger logx.Logger) (string, error) {
	// Check parameters

	// define cli
	cli, err := release.create()
	if err != nil {
		return "", fmt.Errorf("failed to create helm add release command: %w", err)
	}

	// // play cli
	// output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	// }

	// // return response
	// return output, nil

	return cli, nil
}
