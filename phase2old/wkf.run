// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"

	"github.com/abtransitionit/gocore/viperx"
)

func (wf *Workflow) Execute(cfg *viperx.Config) {
	// toposort the phases of the workflow
	phases, _ := wf.TopoSorted()
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

// func (wf *Workflow) Execute(cfg *viperx.Config, fnRegistry map[string]PhaseFunc3, fnCaller func(PhaseFunc3, []string)) {
// 	// toposort the phases of the workflow
// 	phases, _ := wf.TopoSorted2()
// 	// Loop over the phases
// 	for _, phase := range phases {

// 		// get the function and node for this phase
// 		nodes := ResolveNode(cfg, phase.Node)
// 		fn, err := ResolveFn3(fnRegistry, phase.Fn)

// 		// manage errors
// 		if err != nil {
// 			fmt.Println("pbs with function", phase.Fn, err)
// 			continue
// 		}
// 		// log
// 		fmt.Printf("👉 Executing phase (%s) > function (%s) > node (%v)\n", phase.Name, phase.Fn, nodes)

// 		// call the function
// 		fnCaller(fn, nodes)
// 	}
// }

// func (wf *Workflow) Execute(cfg *viperx.Config) {
// 	// toposort the phases of the workflow
// 	phases, _ := wf.TopoSorted2()
// 	// Loop over the phases
// 	for _, phase := range phases {
// 		// get the nodes for this phase
// 		nodes := ResolveNode(cfg, phase.Node)
// 		fmt.Printf("👉 Executing phase %s on nodes %v\n", phase.Name, nodes)
// 		// get the function of THE phase
// 		if fn, ok := ResolveFn[phase.Fn]; ok {
// 			// call the function
// 			fn(nodes) // call mocked function
// 		} else {
// 			fmt.Println("mock: no function for", phase.Fn)
// 		}
// 		fmt.Printf("🔹 Executed phase %s on nodes %v\n", phase.Name, nodes)
// 	}
// }
