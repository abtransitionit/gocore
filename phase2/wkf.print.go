// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/list"
)

func (wf *Workflow) Print() {
	wf.printInternal(false)
}

func (wf *Workflow) PrintWithParams() {
	wf.printInternal(true)
}

func (wf *Workflow) printInternal(showParams bool) {
	var b strings.Builder

	if showParams {
		b.WriteString("Phase\tNode\tDescription\tDependencies\tParams\n")
	} else {
		b.WriteString("Phase\tNode\tDescription\tDependencies\n")
	}

	// Topologically sort phases
	sorted, err := wf.TopoSorted()
	if err != nil {
		fmt.Println("Error sorting workflow:", err)
		return
	}

	// Iterate over sorted phases
	for _, p := range sorted {
		deps := "none"
		if len(p.Dependencies) > 0 {
			deps = strings.Join(p.Dependencies, ", ")
		}

		node := p.Node
		if node == "" {
			node = "none"
		}

		if showParams {
			params := "none"
			if len(p.Params) > 0 {
				var kv []string
				for k, v := range p.Params {
					kv = append(kv, fmt.Sprintf("%s=%s", k, v))
				}
				params = strings.Join(kv, ", ")
			}
			b.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
				p.Name, node, p.Description, deps, params))
		} else {
			b.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\n",
				p.Name, node, p.Description, deps))
		}
	}

	list.PrettyPrintTable(b.String())
}
