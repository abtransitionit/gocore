package kubectl

import (
	"fmt"
	"strings"
)

func (res Resource) List() (string, error) {
	// check parameters
	if res.Type == "" {
		return "", fmt.Errorf("resource type cannot be empty")
	}
	// Define options dynamically
	opts := ""
	switch res.Type {
	case "pod":
		opts = "-A"
	case "node":
		opts = "-o wide | awk '{print $1,$8,$6,$2,$4,$3}' | column -t"
		// opts = "-o wide "
	}

	var cmds = []string{
		strings.TrimSpace(fmt.Sprintf("kubectl get %s %s", res.Type, opts)),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

func (res Resource) Describe() (string, error) {
	// check parameters
	if res.Type == "" {
		return "", fmt.Errorf("resource type cannot be empty")
	}

	var cmds = []string{
		strings.TrimSpace(fmt.Sprintf("kubectl describe %s %s", res.Type, res.Name)),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}
