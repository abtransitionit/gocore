/*
Copyright © 2025 Amar BELGACEM abtransitionit@hotmail.com
*/
package kubectl

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/abelgacem/lucg/config"

)


var ListQueryProperty = []kubectlQueryProperty{
	{new(bool), "apiserverip",    "List API server IP",   `kubectl get endpoints kubernetes -o jsonpath='{.subsets[0].addresses[0].ip}'`},
}

// Parent command
var propertyShortDesc = "Display specific informations about a KBE cluster"
var propertyCmd = &cobra.Command{
	Use:   "property",
	Short: propertyShortDesc,
	Run: func(cmd *cobra.Command, args []string) {
		// printscreen
		fmt.Println("\n🟦",propertyShortDesc)

		// choose query to play accoding to user flags
		for _, query := range ListQueryProperty {
			if query.flagVar != nil && *query.flagVar {
				kubectlCli := query.cli
				cli := fmt.Sprintf(`%s`, kubectlCli)
				// Play the request
				output,cerr,err := config.PlayQueryKubectl(cli,config.KubectlOptions{})
				if err != nil { fmt.Fprintln(os.Stderr, cerr); return}
						// display the output
				fmt.Println(output)
			}
		}
	},
}

	
func init() {
	for _, query := range ListQueryProperty {
		propertyCmd.Flags().BoolVar(query.flagVar, query.name, false, query.description)
	}
}


