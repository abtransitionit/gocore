// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"strings"
)

func (wf *Workflow) GetTablePhase() (string, error) {
	return wf.getTablePhaseInternal(false)
}

func (wf *Workflow) GetTablePhaseWithParams() (string, error) {
	return wf.getTablePhaseInternal(true)
}

func (wf *Workflow) getTablePhaseInternal(showParams bool) (string, error) {
	var b strings.Builder

	if showParams {
		b.WriteString("Phase\tNode\tDescription\tDependencies\tParams\n")
	} else {
		b.WriteString("Phase\tNode\tDescription\tDependencies\n")
	}

	// Topologically sort phases
	sorted, err := wf.TopoPhaseSorted()
	if err != nil {
		fmt.Println("Error sorting workflow:", err)
		return "", err
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
