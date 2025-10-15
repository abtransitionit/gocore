package helm

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

// Returns the cli to add a repo
func (repo HelmRepo) Add() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm repo add %s %s`, repo.Name, repo.Url),
		`helm repo update`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}

// Returns the cli to delete a repo
func (repo HelmRepo) Delete() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm repo remove %s`, repo.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}

// Returns the cli to list all repositories
func (repo HelmRepo) List() (string, error) {
	var cmds = []string{
		`helm repo list`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Returns the cli to list the chart in a repo
func (repo HelmRepo) ListChart() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm search repo %s`, repo.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Returns the list of helm repositories
func ListRepo(local bool, remoteHost string, logger logx.Logger) (string, error) {

	// define cli
	cli, err := HelmRepo{}.List()
	if err != nil {
		return "", fmt.Errorf("failed to build helm list command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, HandleHelmError)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil
}

// Returns the list of helm reposotories referenced in the whitelist
func ListRepoReferenced(local bool, remoteHost string, logger logx.Logger) (string, error) {
	var output strings.Builder

	// Header line — tab-separated
	output.WriteString("NAME\tURL\tDESCRIPTION\n")

	for _, r := range MapHelmRepoReference {
		line := fmt.Sprintf("%s\t%s\t%s\n", r.Name, r.Url, r.Desc)
		output.WriteString(line)
	}

	return output.String(), nil
}

// Add a helm repository to the client configuration file
func AddRepo(local bool, remoteHost string, repo HelmRepo, logger logx.Logger) (string, error) {

	// Check parameters

	// define cli
	cli, err := repo.Add()
	if err != nil {
		return "", fmt.Errorf("failed to add helm repository command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	// return response
	return output, nil

}
