// File in gocore/phase/adapter.go
package phase

import (
	"fmt"

	"github.com/abtransitionit/gocore/viperx"
)

var ResolveFn = map[string]func([]string){
	"checkVmAccess": func(nodes []string) { fmt.Println("mock: checkVmAccess") },
	"copyAgent":     func(nodes []string) { fmt.Println("mock: copyAgent") },
	"osUpgrade":     func(nodes []string) { fmt.Println("mock: osUpgrade") },
}

func (wf *Workflow2) Execute(cfg *viperx.Config) {
	// toposort the phases of the workflow
	phases, _ := wf.TopoSorted2()
	// Loop over the phases
	for _, phase := range phases {
		// get the nodes for this phase
		nodes := ResolveNode(cfg, phase.Node)
		fmt.Printf("👉 Executing phase %s on nodes %v\n", phase.Name, nodes)
		// get the function of THE phase
		if fn, ok := ResolveFn[phase.Fn]; ok {
			// call the function
			fn(nodes) // call mocked function
		} else {
			fmt.Println("mock: no function for", phase.Fn)
		}
		fmt.Printf("🔹 Executed phase %s on nodes %v\n", phase.Name, nodes)
	}
}
