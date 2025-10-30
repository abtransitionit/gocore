// File in gocore/phase/adapter.go
package phase

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/abtransitionit/gocore/yamlx"
)

var Funcs = map[string]func([]string){
	"checkVmAccess": func(nodes []string) { fmt.Println("mock: checkVmAccess") },
	"copyAgent":     func(nodes []string) { fmt.Println("mock: copyAgent") },
	"osUpgrade":     func(nodes []string) { fmt.Println("mock: osUpgrade") },
}

func (wf *Workflow2) Execute() {
	// toposort the phases of the workflow
	phases, _ := wf.TopoSorted2()
	// Loop over the phases
	for _, phase := range phases {
		// get the nodes of a THE phase
		nodes := ResolveNodeSet(phase.Node)
		fmt.Printf("👉 manager: Executing phase %s on nodes %v\n", phase.Name, nodes)
		// get the function THE phase
		if fn, ok := Funcs[phase.Fn]; ok {
			// call the function
			fn(nodes) // call mocked function
		} else {
			fmt.Println("mock: no function for", phase.Fn)
		}
		fmt.Printf("👉 manager: Executed phase %s on nodes %v\n", phase.Name, nodes)
	}
}

func GetWorkflow() (*Workflow2, error) {
	// 1. Define YAML workflow file path
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("could not get caller information")
	}
	workflowPath := filepath.Join(path.Dir(file), "wkf.yaml")

	fmt.Println("workflowPath:", workflowPath)

	// 2. Load the yaml using the generic function from lib.go
	workflow, err := yamlx.LoadFile[Workflow2](workflowPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load workflow from %s: %w", workflowPath, err)
	}

	return workflow, nil

}
