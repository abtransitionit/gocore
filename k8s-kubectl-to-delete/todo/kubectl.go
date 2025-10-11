/*
Copyright © 2025 Amar BELGACEM abtransitionit@hotmail.com
*/
package kubectl

import (
	"github.com/spf13/cobra"
)

type kubectlQueryProperty struct {
	flagVar     *bool
	name        string
	description string
	cli         string
}

// command description
var kubectlShortDesc = "Manage kubectl requests onto the cluster via the API Server"
var kubectlLongDesc = kubectlShortDesc + `. The CLI must be played from a cokpit VM.
	→ The cokpit VM 
		→ must not be a KBE Node, nor the VM hosted Kubectl.
		→ must have access to the VM that host the kubectl CLI.
	→ The VM that host kubectl CLI must also have the CLI LUC installed.`

// Parent command code
var KubectlCmd = &cobra.Command{
	Use:   "kctl",
	Short: kubectlShortDesc,
	Long:  kubectlLongDesc,
}

func init() {
	KubectlCmd.AddCommand(healthCmd)
	KubectlCmd.AddCommand(describeCmd)
	KubectlCmd.AddCommand(propertyCmd)
	// KubectlCmd.AddCommand(apiCmd)
	// KubectlCmd.AddCommand(customCmd)
}

func FormatKubectlResult(kubectlResult string) (result string, customErr string, err error) {
	return "", "", nil
}
