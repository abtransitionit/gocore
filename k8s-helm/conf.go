package helm

import (
	"strings"
)

// Return: The cli to list the helm client envars
func GetEnv() (string, error) {
	var cmds = []string{
		`helm env`,
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil

}
