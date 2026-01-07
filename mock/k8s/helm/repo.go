package helm

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/golinux/mock/run"
)

func (i *Repo) Add(hostName string, logger logx.Logger) error {
	// 1 - get and play cli
	if _, err := run.RunCli(hostName, i.cliToAdd(), logger); err != nil {
		return err
	}
	// handle success
	logger.Debugf("%s:%s > started OS service for the current session", hostName, i.Name)
	return nil
}

func (i *Repo) List(hostName string, logger logx.Logger) error {
	// 1 - get and play cli
	if _, err := run.RunCli(hostName, i.cliToList(), logger); err != nil {
		return err
	}
	// handle success
	logger.Debugf("%s:%s > started OS service for the current session", hostName, i.Name)
	return nil
}

func (i *Repo) ListChart(hostName string, logger logx.Logger) error {
	// 1 - get and play cli
	if _, err := run.RunCli(hostName, i.cliToListChart(), logger); err != nil {
		return err
	}
	// handle success
	logger.Debugf("%s:%s > started OS service for the current session", hostName, i.Name)
	return nil
}

func (i *Repo) cliToAdd() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm repo add %s %s`, i.Name, i.Url),
		`helm repo update`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}

func (i *Repo) cliToDelete() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm repo remove %s`, i.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}

// Returns the cli to list all repositories
func (i *Repo) cliToList() (string, error) {
	var cmds = []string{
		`helm repo list`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Returns the cli to list the chart in a repo
func (i *Repo) cliToListChart() (string, error) {
	var cmds = []string{
		fmt.Sprintf(`helm search repo %s`, i.Name),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}
