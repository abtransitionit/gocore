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
		fmt.Sprintf(`helm create %s`, chart.FullName),
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

// Return: The cli to list the kind of all templates a helm charts it will create
// Todo: it is linux specific. should be in library:golinux not gocore. or change the code to be OS agnostic
func (chart HelmChart) ListNbKind() (string, error) {

	var cmds = []string{
		fmt.Sprintf(`echo -e "Kind\tCount" && helm template %s | grep '^kind:' | sort | uniq -c | sed 's/kind: *//' | awk '{printf "%%s\t%%s\n", $2, $1}'`, chart.FullName),
		// fmt.Sprintf(`echo "COUNT  KIND" && helm template %s | grep '^kind:' | sort | uniq -c | sed 's/kind: *//' | awk '{print $1, $2":", $3}'`, chart.FullName),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}
func (chart HelmChart) ListKind() (string, error) {

	var cmds = []string{
		fmt.Sprintf(`echo -e "Kind\tName" && helm template %s | grep -A2 '^kind:' | grep -A2 '^kind:' | egrep 'kind|name' | paste - - | sed 's/kind: *//' | sed 's/name: *//'`, chart.FullName),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

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

// Returns the list of all kind the chart will create
func (chart HelmChart) ListNbChartKind(local bool, remoteHost string, logger logx.Logger) (string, error) {

	// define cli
	cli, err := chart.ListNbKind()
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

// Returns the list of all kind the chart will create
func (chart HelmChart) ListChartKind(local bool, remoteHost string, logger logx.Logger) (string, error) {

	// define cli
	cli, err := chart.ListKind()
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

// func ListKind(local bool, remoteHost string, chart HelmChart, logger logx.Logger) (string, error) {
// 	// define cli
// 	cli, err := chart.ListKind()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to build cli: %w", err)
// 	}

// 	// play cli
// 	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
// 	}

// 	// return response
// 	return output, nil

// }

// display chart metadata - chart.yaml
// helm show chart $chartName/$RepoName

// display chart metadata - values.yaml
// helm show values $chartName/$RepoName

// var releaseValueShortDesc = "Display user defined values about a relase"
// 	Example: `
// 	xxx kbe-cilicium  kube-system
// cli := fmt.Sprintf(`helm get values %s -n %s`, args[0], args[1])
