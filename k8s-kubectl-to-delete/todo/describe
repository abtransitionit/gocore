package kubectl

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/abelgacem/lucg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type describeQueryProperty struct {
	flagVar     *string
	name        string
	description string
	namespaced  bool
	cli         string
}

var nodeArg, podArg, nsArg string
var ListQueryDescribe = []describeQueryProperty{
	// {&nodeArg, "res",     "describe any resource",  false, `kubectl get all -A`},
	{&nsArg, "configmap", "describe a configmap", true, `kubectl get configmap -A`},
	{&nsArg, "crd", "describe a CustomResourceDefinition", false, `kubectl get crd -A`},
	{&nsArg, "deploy", "describe a deployment", true, `kubectl get deployment -A`},
	{&nsArg, "event", "describe events", true, `kubectl get events -A`},
	{&nsArg, "ingress", "describe an ingress", true, `kubectl get ingress -A`},
	{&nodeArg, "node", "describe a node", false, `kubectl get node`},
	{&nsArg, "ns", "describe a ns", false, `kubectl get ns`},
	{&podArg, "pod", "describe a pod", true, `kubectl get pod -A`},
	{&nodeArg, "resn", "list all namespaced resources", true, `kubectl get all -A`},
	{&podArg, "rs", "describe a replicaset", true, `kubectl get rs -A`},
	{&nsArg, "service", "describe a service", true, `kubectl get service -A`},
}

// Parent command
var describeShortDesc = "Describe a specific resource (ie. node, pod, ...)"
var describeCmd = &cobra.Command{
	Use:   "desc",
	Short: describeShortDesc,
	Example: `
		desc --service
		desc --pod
		desc --pod=2
		desc --pod=4
		desc --service=12
	`,
	Run: func(cmd *cobra.Command, args []string) {
		nbFlag := 0
		cmd.Flags().Visit(func(*pflag.Flag) { nbFlag++ })
		if nbFlag < 1 {
			fmt.Fprintln(os.Stderr, "❌ Error: you must specify a flag (ie. --node --pod, ...")
			return
		}
		// Initialize resource to manage in the List
		var resIdx int = -1
		var flagValue, flagName string
		var flagNsBool bool
		var resType, resName, resNamespace = "", "", ""

		// Printscreen
		fmt.Println("\n🟦", describeShortDesc)

		// define the request according to user flag
		for i, query := range ListQueryDescribe {
			if cmd.Flags().Changed(query.name) {
				resIdx = i
				flagValue = cmd.Flag(query.name).Value.String()
				flagName = cmd.Flag(query.name).Name
				flagNsBool = query.namespaced
			}
		}

		// Get the output of the request whether or not the user gives a value to the Flag
		cli := fmt.Sprintf(`%s`, ListQueryDescribe[resIdx].cli)
		queryResult, cerr, err := config.PlayQueryKubectl(cli, config.KubectlOptions{FormatOutput: true})
		if err != nil {
			fmt.Fprintln(os.Stderr, cerr)
			return
		}

		if strings.TrimSpace(flagValue) == "" {
			// No value is given - display the result
			fmt.Printf("%s", queryResult)
			return
		} else {
			// A value is given - describe a specific type of resource
			flagValueInt, err := strconv.Atoi(flagValue)
			if err != nil {
				fmt.Fprintln(os.Stderr, "❌ Error, cannot convert user choice to int: ", err)
				os.Exit(1)
			}

			// Get the line in the output corresponding to the flagValue
			lineArray := strings.Fields(strings.Split(queryResult, "\n")[flagValueInt])

			// get resource type
			resType = flagName

			// Determine if we're dealing with a namespaced resource or not based on the field in the line
			if flagNsBool {
				// For namespaced resources: lineArray[2] is the namespace name, lineArray[3] is the resource name
				resNamespace = lineArray[2]
				resName = lineArray[3]
				cli = fmt.Sprintf(`kubectl describe %s/%s -n %s`, resType, resName, resNamespace)
			} else {
				// For non-namespaced resources: lineArray[2] is the resource name
				resName = lineArray[2]
				cli = fmt.Sprintf(`kubectl describe %s/%s`, resType, resName)
			}

			// check if we are in detail mode
			if getFlag && flagNsBool {
				// if yes - get the yaml
				cli = fmt.Sprintf(`kubectl get %s/%s -o yaml -n %s`, resType, resName, resNamespace)
			} else if getFlag {
				cli = fmt.Sprintf(`kubectl get %s/%s -o yaml`, resType, resName)
			}

			if getFlag && !flagNsBool && resType == "ns" {
				fmt.Printf("resName: %s\n", resName)
				fmt.Printf("resNamespace: %s\n", resNamespace)
				fmt.Printf("resType: %s\n", resType)
				cli = fmt.Sprintf(`kubectl get all -n %s`, resName)
				fmt.Printf("cli: %s\n", cli)
			}

			// play the cli
			queryResult, cerr, err := config.PlayQueryKubectl(cli, config.KubectlOptions{FormatOutput: true})
			if err != nil {
				fmt.Fprintln(os.Stderr, cerr)
				return
			}
			fmt.Printf("%s", queryResult)
			return
		}
	},
}

var getFlag bool

func init() {
	describeCmd.Flags().BoolVar(&getFlag, "get", false, "kubectl get the resource instead kubectl describe it")
	for _, query := range ListQueryDescribe {
		describeCmd.Flags().StringVar(query.flagVar, query.name, "", query.description)
		describeCmd.Flag(query.name).NoOptDefVal = " "
	}
}

// Todo
// Flag -h should display better

// Todo
// now the code provide             > kubectl describe kind/<resName> [-n namespace]
// add a flag --detail that provide > kubectl get      kind/<resName> [-n namespace] -o yaml

/*
# kubectl describe

| Resource Type          | Required Information          | Example Command                     |
|------------------------|-------------------------------|-------------------------------------|
| pod/<name>             | name, namespace               | `kubectl describe pod/nginx -n web` |
| node/<name>            | name                          | `kubectl describe node/node01`      |
| ns/<name>              | name                          | `kubectl describe ns/default`       |
| deployment/<name>      | name, namespace               | `kubectl describe deploy/nginx -n web` |
| daemonset/<name>       | name, namespace               | `kubectl describe ds/fluentd -n logging` |
| statefulset/<name>     | name, namespace               | `kubectl describe sts/mysql -n db`  |
| replicaset/<name>      | name, namespace               | `kubectl describe rs/nginx-xyz -n web` |
| service/<name>         | name, namespace               | `kubectl describe svc/nginx -n web` |
| ingress/<name>         | name, namespace               | `kubectl describe ingress/app -n web` |
| configmap/<name>       | name, namespace               | `kubectl describe cm/config -n app` |
| secret/<name>          | name, namespace               | `kubectl describe secret/db-creds -n app` |
| persistentvolume/<name> | name                        | `kubectl describe pv/pv001`         |
| persistentvolumeclaim/<name> | name, namespace        | `kubectl describe pvc/data -n app`  |
| storageclass/<name>   | name                          | `kubectl describe sc/fast`          |
| serviceaccount/<name> | name, namespace               | `kubectl describe sa/default -n app` |
| role/<name>          | name, namespace               | `kubectl describe role/admin -n app` |
| rolebinding/<name>   | name, namespace               | `kubectl describe rolebinding/admin -n app` |
| clusterrole/<name>   | name                          | `kubectl describe clusterrole/admin` |
| clusterrolebinding/<name> | name                     | `kubectl describe clusterrolebinding/admin` |
| cronjob/<name>       | name, namespace               | `kubectl describe cj/backup -n ops` |
| job/<name>           | name, namespace               | `kubectl describe job/migration -n db` |
| endpoint/<name>      | name, namespace               | `kubectl describe ep/nginx -n web`  |
| event/<name>         | name, namespace               | `kubectl describe event/xyz -n web` |
| limitrange/<name>    | name, namespace               | `kubectl describe limits/default -n app` |
| resourcequota/<name> | name, namespace               | `kubectl describe quota/default -n app` |
| hpa/<name>           | name, namespace               | `kubectl describe hpa/frontend -n web` |
| networkpolicy/<name> | name, namespace               | `kubectl describe netpol/default -n app` |
| poddisruptionbudget/<name> | name, namespace       | `kubectl describe pdb/zk -n db`     |
| priorityclass/<name> | name                          | `kubectl describe pc/high`          |
| runtimeclass/<name>  | name                          | `kubectl describe rc/gvisor`        |
| csidriver/<name>     | name                          | `kubectl describe csidriver/ebs`    |
| csinode/<name>       | name                          | `kubectl describe csinode/node01`   |
| volumeattachment/<name> | name                      | `kubectl describe va/va-xyz`        |

*/

// // Limit the number of flag passsed to 1
// nbFlag:=0;cmd.Flags().Visit(func(*pflag.Flag) {nbFlag++})
// if nbFlag != 1 {
// 	fmt.Fprintln(os.Stderr, "❌ Error: you must specify one uniq flag  (ie. --node --pod, ...")
// 	return
// }
