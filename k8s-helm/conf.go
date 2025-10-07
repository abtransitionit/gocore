package helm

import (
	"strings"
)

// list helm envars
func GetEnv() (string, error) {
	var cmds = []string{
		`helm env`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}
