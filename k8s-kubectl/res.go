package kubectl

import (
	"fmt"
	"strings"
)

// Return: The cli:Flag for the namespacetspecific to the resource
func (res Resource) nsFlag() string {
	if res.Ns != "" {
		return fmt.Sprintf("-n %s", res.Ns)
	}
	return ""
}

// Return: The cli to list that resource
func (res Resource) List() (string, error) {
	// check parameters
	if res.Type == "" {
		return "", fmt.Errorf("resource type cannot be empty")
	}

	// Define options dynamically
	resourceOptions := map[string]string{
		"pod":  "-A",
		"node": "-o wide | awk '{print $1,$8,$(NF-1),$6,$2,$4,$3}' | column -t",
		// "node": "-o wide ",
	}
	opts := resourceOptions[res.Type]

	// define cli
	var cmds = []string{
		strings.TrimSpace(fmt.Sprintf("kubectl get %s %s", res.Type, opts)),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Return: The cli to describe that resource
func (res Resource) Describe() (string, error) {
	// check parameters
	if res.Type == "" {
		return "", fmt.Errorf("resource type cannot be empty")
	}
	// check parameters
	if res.Name == "" {
		return "", fmt.Errorf("resource name cannot be empty")
	}

	var cmds = []string{
		strings.TrimSpace(fmt.Sprintf("kubectl describe %s %s %s", res.Type, res.Name, res.nsFlag())),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}

// Return: The cli to yaml that resource
func (res Resource) Yaml() (string, error) {
	// check parameters
	if res.Type == "" {
		return "", fmt.Errorf("resource type cannot be empty")
	}

	var cmds = []string{
		strings.TrimSpace(fmt.Sprintf("kubectl get %s %s %s -o yaml", res.Type, res.Name, res.nsFlag())),
	}
	cli := strings.Join(cmds, " && ")
	return cli, nil
}
