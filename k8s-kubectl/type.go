package kubectl

// Generic resource struct
type Resource struct {
	Name string
	Type string
}

// type Node struct {
// 	Name string
// 	Type string
// }
// type Pod struct {
// 	Name string
// 	Type string
// }

// type Ns struct {
// 	Name string
// 	Type string
// }

// kubectl get pods -n db         # see Kubernetes pods

// kubectl get all -n db
