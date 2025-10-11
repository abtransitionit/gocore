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

var ListQueryHealth = []kubectlQueryProperty{
	{new(bool), "res",     "List (all) resources (ie. namespaced or not), group by kind",   `kubectl get \$(kubectl api-resources --verbs=list -o name --sort-by=name | paste -sd ,) -A`},
	{new(bool), "resn",    "List (only) namespaded resources: group by kind",   	`kubectl get all -A`},
	{new(bool), "cnipod",  "List CNI pods agents only status: 1 pod  per line",   `kubectl -n kube-system get pods -l k8s-app=cilium -o wide | awk \"{print \\\$1,\\\$2,\\\$3,\\\$4,\\\$5,\\\$6,\\\$7}\" | column -t`},
	{new(bool), "label",   "List node labels (which allow node categorization)", 	`kubectl get nodes -o custom-columns=\"NODE:.metadata.name, LABELS:.metadata.labels\"`},
	{new(bool), "node",    "List nodes info:  1 node per line",       						`kubectl get nodes -o wide | awk \"{print \\\$1, \\\$2, \\\$3, \\\$4, \\\$6,\\\$8}\" | column -t`},
	{new(bool), "ns",      "List namespace info: 1 namespace per line",   				`kubectl get ns`},
	{new(bool), "pod",     "List pods  info:  1 pod  per line",                   `kubectl get pods -A -o wide | awk \"{print \\\$1,\\\$2,\\\$3,\\\$4,\\\$7,\\\$8,\\\$6}\" | column -t`},
	{new(bool), "taint",   "List node taints (which allow or not the scheduling of workload/user pods)",   `kubectl get nodes -o custom-columns=\"NODE:.metadata.name, TAINTS:.spec.taints\"`},
	{new(bool), "test",    "List API server IP",   `kubectl get endpoints kubernetes -o jsonpath='{.subsets[0].addresses[0].ip}'`},
}

// Parent command
var healthShortDesc = "List informations about a KBE cluster"
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: healthShortDesc,
	Run: func(cmd *cobra.Command, args []string) {
		// printscreen
		// fmt.Println("\n🟦",healthShortDesc)

		// choose query to play accoding to user flags
		for _, query := range ListQueryHealth {
			if query.flagVar != nil && *query.flagVar {
				kubectlCli := query.cli
				cli := fmt.Sprintf(`%s`, kubectlCli)
				// Play the request
				output,cerr,err := config.PlayQueryKubectl(cli,config.KubectlOptions{FormatOutput: true})
				if err != nil { fmt.Fprintln(os.Stderr, cerr); return}
						// display the output
				fmt.Println(output)
			}
		}
	},
}

	
func init() {
	for _, query := range ListQueryHealth {
		healthCmd.Flags().BoolVar(query.flagVar, query.name, false, query.description)
	}
}
