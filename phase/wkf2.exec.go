// File in gocore/phase/adapter.go
package phase

import (
	"fmt"

	"github.com/abtransitionit/gocore/viperx"
)

func ResolveFn2(fnRegistry map[string]any, name string) (any, error) {
	fn, ok := fnRegistry[name]
	if !ok {
		return nil, fmt.Errorf("no function registered for %q", name)
	}
	return fn, nil
}
func ResolveFn3(fnRegistry map[string]PhaseFunc3, name string) (PhaseFunc3, error) {
	fn, ok := fnRegistry[name]
	if !ok {
		return nil, fmt.Errorf("no function registered for %q", name)
	}
	return fn, nil
}

// var ResolveFn = map[string]func([]string){
// 	"checkVmAccess": func(nodes []string) { fmt.Println("mock: checkVmAccess") },
// 	"copyAgent":     func(nodes []string) { fmt.Println("mock: copyAgent") },
// 	"osUpgrade":     func(nodes []string) { fmt.Println("mock: osUpgrade") },
// }

// func (wf *Workflow2) Execute2(cfg *viperx.Config, fnRegistry map[string]any)
// 	// toposort the phases of the workflow
// 	phases, _ := wf.TopoSorted2()
// 	// Loop over the phases
// 	for _, phase := range phases {

// 		// get the function and node for this phase
// 		nodes := ResolveNode(cfg, phase.Node)
// 		fn, err := ResolveFn2(fnRegistry, phase.Fn)

// 		// manage errors
// 		if err != nil {
// 			fmt.Println("pbs with function", phase.Fn, err)
// 			continue
// 		}
// 		// log
// 		fmt.Printf("👉 Executing phase (%s) > function (%s) > node (%v)\n", phase.Name, phase.Fn, nodes)

// 		// call the function
// 		fn
// 	}
// }

func (wf *Workflow2) Execute3(cfg *viperx.Config, fnRegistry map[string]PhaseFunc3, fnCaller func(PhaseFunc3, []string)) {
	// toposort the phases of the workflow
	phases, _ := wf.TopoSorted2()
	// Loop over the phases
	for _, phase := range phases {

		// get the function and node for this phase
		nodes := ResolveNode(cfg, phase.Node)
		fn, err := ResolveFn3(fnRegistry, phase.Fn)

		// manage errors
		if err != nil {
			fmt.Println("pbs with function", phase.Fn, err)
			continue
		}
		// log
		fmt.Printf("👉 Executing phase (%s) > function (%s) > node (%v)\n", phase.Name, phase.Fn, nodes)

		// call the function
		fnCaller(fn, nodes)
	}
}
