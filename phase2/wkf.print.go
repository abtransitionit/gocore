// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"strings"
)

func (wf *Workflow) GetTablePhase() (string, error) {
	return wf.getTablePhaseInternal()
}

func (wf *Workflow) getTablePhaseInternal() (string, error) {
	var b strings.Builder

	// Header
	b.WriteString("Phase\tNode\tFn\tParam\n")

	// Topologically sort phases
	sorted, err := wf.TopoPhaseSorted()
	if err != nil {
		fmt.Println("Error sorting workflow:", err)
		return "", err
	}

	// Iterate over sorted phases
	for _, p := range sorted {
		// node
		node := p.Node
		if node == "" {
			node = "none"
		}

		// fn
		fn := p.Fn
		if fn == "" {
			fn = "none"
		}

		// param
		params := "none"
		if len(p.Param) > 0 {
			params = strings.Join(p.Param, ", ")
		}

		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n", p.Name, node, fn, params)
	}

	return b.String(), nil
}

func (wf *Workflow) GetTableTier() (string, error) {
	return wf.getTableTierInternal(false)
}

func (wf *Workflow) GetTableTierWithParams() (string, error) {
	return wf.getTableTierInternal(true)
}

func (wf *Workflow) getTableTierInternal(showParams bool) (string, error) {
	sortedTiers, err := wf.TopoTierSorted()
	if err != nil {
		fmt.Println("Error sorting workflow by tiers:", err)
		return "", err
	}

	var b strings.Builder

	// Table header
	if showParams {
		b.WriteString("Tier\tIdP\tPhase\tNode\tDescription\tDependencies\tParams\n")
	} else {
		b.WriteString("Tier\tIdP\tPhase\tNode\tDescription\tDependencies\n")
	}

	// Iterate through tiers
	for tierIndex, tier := range sortedTiers {
		tierID := tierIndex + 1
		for phaseIndex, p := range tier {
			idp := phaseIndex + 1

			deps := "none"
			if len(p.Dependency) > 0 {
				deps = strings.Join(p.Dependency, ", ")
			}

			node := p.Node
			if node == "" {
				node = "none"
			}

			if showParams {
				params := "none"
				if len(p.Param) > 0 {
					var kv []string
					for k, v := range p.Param {
						kv = append(kv, fmt.Sprintf("%s=%s", k, v))
					}
					params = strings.Join(kv, ", ")
				}
				b.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%s\t%s\t%s\n",
					tierID, idp, p.Name, node, p.Description, deps, params))
			} else {
				b.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%s\t%s\n",
					tierID, idp, p.Name, node, p.Description, deps))
			}
		}
	}

	return b.String(), nil
}
