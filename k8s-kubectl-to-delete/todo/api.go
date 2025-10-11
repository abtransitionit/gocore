/*
Copyright © 2025 Amar BELGACEM abtransitionit@hotmail.com
*/
package kubectl

// list of all resource API server know
// kubectl api-resources --verbs=list   --sort-by=name -o wide | less -S
// kubectl api-resources --verbs=create --sort-by=name -o wide | grep -v create, | less -S

var ListQueryApi = []kubectlQueryProperty{
	{new(bool), "res", "List all resources sort by name", `kubectl api-resources -o wide --sort-by=name`},
	{new(bool), "res-ns", "List namespaced resources", `kubectl api-resources -o wide --sort-by=name --namespaced=false`},
	{new(bool), "res-nons", "List not namespaced resources", `kubectl api-resources -o wide --sort-by=name --namespaced=true`},
	{new(bool), "apiserverip", "List API server IP", `kubectl get endpoints kubernetes -o jsonpath='{.subsets[0].addresses[0].ip}'`},
	{new(bool), "res", "List (all) resources (ie. namespaced or not), group by kind", `kubectl get \$(kubectl api-resources --verbs=list -o name --sort-by=name | paste -sd ,) -A`},
	{new(bool), "resn", "List (only) namespaded resources: group by kind", `kubectl get all -A`},
	{new(bool), "cnipod", "List CNI pods agents only status: 1 pod  per line", `kubectl -n kube-system get pods -l k8s-app=cilium -o wide | awk \"{print \\\$1,\\\$2,\\\$3,\\\$4,\\\$5,\\\$6,\\\$7}\" | column -t`},
	{new(bool), "label", "List node labels (which allow node categorization)", `kubectl get nodes -o custom-columns=\"NODE:.metadata.name, LABELS:.metadata.labels\"`},
	{new(bool), "node", "List nodes info:  1 node per line", `kubectl get nodes -o wide | awk \"{print \\\$1, \\\$2, \\\$3, \\\$4, \\\$6,\\\$8}\" | column -t`},
	{new(bool), "ns", "List namespace info: 1 namespace per line", `kubectl get ns`},
	{new(bool), "pod", "List pods  info:  1 pod  per line", `kubectl get pods -A -o wide | awk \"{print \\\$1,\\\$2,\\\$3,\\\$4,\\\$7,\\\$8,\\\$6}\" | column -t`},
	{new(bool), "taint", "List node taints (which allow or not the scheduling of workload/user pods)", `kubectl get nodes -o custom-columns=\"NODE:.metadata.name, TAINTS:.spec.taints\"`},
	{new(bool), "test", "List API server IP", `kubectl get endpoints kubernetes -o jsonpath='{.subsets[0].addresses[0].ip}'`},
	// {new(bool), "apires",    "List resource",  `export TERM=xterm-256color; tput rmam; kubectl api-resources -o wide --sort-by=name; tput sman 2>/dev/null`},
	// {new(bool), "res-nons",  "List not namespaced resources", `kubectl api-resources -o wide --sort-by=name |  awk \"NR == 1 || /true/\"`},
}

// do  | less -S | cat -n | column -t

// for resource in $(kubectl api-resources --verbs=list --namespaced -o name | sort); do   kubectl get "$resource" --all-namespaces --ignore-not-found --sort-by=.metadata.namespace; done
//kubectl api-resources --verbs=list -o wide
// kubectl get pods -n db         # see Kubernetes pods

// kubectl get all -n db
// kubectl get $(kubectl api-resources --verbs=list -o name | paste -sd ",") -A  | less -S
// kubectl get $(kubectl api-resources --verbs=list -o name --sort-by=name | paste -sd ",") -A | less -S
// kubectl get nodes,pods, all -A | less -S
// kubectl get events --all-namespaces --sort-by=.metadata.creationTimestamp | grep -i warning
// kubectl get events --all-namespaces --sort-by=.metadata.creationTimestamp | grep -i warning
// kubectl get all -n kdashb

// kubectl logs nginx-pod-12345 -n kube-system
// kubectl logs -f nginx-pod-12345             # Follow logs (like tail -f)
// kubectl logs --previous nginx-pod-12345     # View logs from a previously crashed container
// kubectl logs <pod-name> -c <container-name> # pod with multiple container

// kubectl describe pod <pod-name>
// kubectl describe node <node-name>
// kubectl describe svc <service-name>

// kubectl top nodes
// kubectl top pods --all-namespaces
// kubectl get service -A --watch # live
// kubectl exec -it <pod-name> -- /bin/sh # enter a pod
// kubectl run -it --rm debug --image=busybox --restart=Never -- sh
//   - nslookup kube-dns.kube-system.svc.cluster.local
