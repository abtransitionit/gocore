package cilium

import "strings"

func Status() (string, error) {
	var cmds = []string{
		`cilium status`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}
