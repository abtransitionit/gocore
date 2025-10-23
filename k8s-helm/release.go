package helm

import (
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
func (release HelmRelease) versionFlag() string {
	if release.Chart.Version != "" {
		return fmt.Sprintf("--version %s", release.Chart.Version)
	}
	return ""
}

// Returns the cli to create a release from a chart into a k8s cluster
//
// Notes:
//   - it can create a release from chart in a repo or chart given as path in the FS
func (release HelmRelease) cliCreate() (string, error) {

	// 1 - create the release
	var cmds = []string{
		fmt.Sprintf(`
			helm install %s %s --atomic --wait %s --namespace %s %s
			`,
			release.Name,
			release.Chart.FullName,
			release.versionFlag(),
			release.Namespace,
			release.valueFlag()),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

func (release HelmRelease) cliDryCreate() (string, error) {

	// 1 - create the release
	var cmds = []string{
		fmt.Sprintf(`
			helm install %s %s --debug --dry-run=server %s --namespace %s %s
			`,
			release.Name,
			release.Chart.FullName,
			release.versionFlag(),
			release.Namespace,
			release.valueFlag()),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Returns the cli to delete a release in a k8s cluster
func (release HelmRelease) cliDelete() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm uninstall %s --namespace %s`, release.Name, release.Namespace),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Returns the cli to describe a release in a k8s cluster - ie. display all prints out all the Kubernetes resources that were uploaded to the server
func (release HelmRelease) cliDescribe() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm get manifest %s --namespace %s`, release.Name, release.Namespace),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Returns the cli to list the releases installed in a k8s cluster
func (release HelmRelease) cliList() (string, error) {
	var cmds = []string{
		"helm list -A", //  list releases in namespace dd
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// create a helm release into a kubernetes cluster
func (release HelmRelease) List(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// Check parameters

	// define cli
	cli, err := HelmRelease{}.cliList()
	if err != nil {
		return "", fmt.Errorf("failed to build helm list command: %w", err)
	}

	// // play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil

}

func (release HelmRelease) Create(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// Check parameters

	// define cli
	cli, err := release.cliCreate()
	if err != nil {
		return "", fmt.Errorf("failed to create helm add release command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil

	// return cli, nil
}
func (release HelmRelease) DryCreate(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// Check parameters

	// define cli
	cli, err := release.cliDryCreate()
	if err != nil {
		return "", fmt.Errorf("failed to create helm add release command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil

	// return cli, nil
}

func (release HelmRelease) Delete(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// Check parameters

	// define cli
	cli, err := release.cliDelete()
	if err != nil {
		return "", fmt.Errorf("failed to create helm add release command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil
}

func (release HelmRelease) Describe(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// Check parameters

	// define cli
	cli, err := release.cliDescribe()
	if err != nil {
		return "", fmt.Errorf("failed to create helm add release command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil
}

// fmt.Sprintf(`helm upgrade --install %s %s/%s --version %s --namespace %s -f %s`, chartReleaseFlag, repoNameFlag, chartNameFlag, chartVersionFlag, k8sNamespaceFlag, fileConfFlag)
// fmt.Sprintf(`helm upgrade --install %s %s/%s --version %s --namespace %s`, chartReleaseFlag, repoNameFlag, chartNameFlag, chartVersionFlag, k8sNamespaceFlag)
