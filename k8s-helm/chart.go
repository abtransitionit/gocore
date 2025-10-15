package helm

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

// Return: The cli to list the helm charts in a repo
func (chart HelmChart) Create() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm create %s`, chart.FullPath),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Return: The cli to list the helm charts in a repo
func (chart HelmChart) List() (string, error) {
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

// Returns the list of helm charts in a helm repo
func ListChart(local bool, remoteHost string, repo HelmRepo, logger logx.Logger) (string, error) {

	// define cli
	// cli, err := helm.HelmRepo{Name: repoName}.ListChart()
	cli, err := repo.ListChart()
	if err != nil {
		return "", fmt.Errorf("failed to build cli: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil
}
func CreateChart(local bool, remoteHost string, chart HelmChart, logger logx.Logger) (string, error) {

	// define cli
	cli, err := chart.Create()
	if err != nil {
		return "", fmt.Errorf("failed to build cli: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil
}
