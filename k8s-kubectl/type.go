package kubectl

// Name: Resource
//
// Description: A Generic K8s resource
type Resource struct {
	Name string
	Type string
	Ns   string
}
