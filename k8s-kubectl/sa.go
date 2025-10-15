package kubectl

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/run"
)

// Returns the list of node as a string
func ListSa(local bool, remoteHost string, logger logx.Logger) (string, error) {
	// define cli
	cli, err := Resource{Type: "sa"}.List()
	if err != nil {
		return "", fmt.Errorf("failed to build kubectl list command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}

func DescribeSa(local bool, remoteHost string, sa Resource, logger logx.Logger) (string, error) {

	// define cli
	cli, err := sa.Describe()
	if err != nil {
		return "", fmt.Errorf("failed to build kubectl list command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}

func YamlSa(local bool, remoteHost string, sa Resource, logger logx.Logger) (string, error) {

	// define cli
	cli, err := sa.Yaml()
	if err != nil {
		return "", fmt.Errorf("failed to build kubectl list command: %w", err)
	}

	// play cli
	output, err := run.ExecuteCliQuery(cli, logger, local, remoteHost, run.NoOpErrorHandler)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %s: %w", cli, err)
	}

	return output, nil
}
