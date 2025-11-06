package phase2

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

// Description: returns a view of the phases in the workflow
//
// Notes:
//
// - Allow to view the workflow phases

func (wf *Workflow) GetPhaseView() (string, error) {
	var b strings.Builder

	// Header
	b.WriteString("Phase\tExe Node\tFn\tParam\n")

	// Topologically sort phases
	sorted, err := wf.topoSortByPhase()
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
		fn := p.FnAlias
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

func (wf *Workflow) GetTierView(tier [][]Phase, logger logx.Logger) (string, error) {
	// sortedTiers, err := wf.topoSortByTier(logger)
	// if err != nil {
	// 	fmt.Println("Error sorting workflow by tiers:", err)
	// 	return "", err
	// }

	var b strings.Builder

	// Table header (no Params column anymore)
	b.WriteString("Tier\tIdP\tPhase\tExe Node\tDescription\tDependencies\n")

	// Iterate through tiers
	for tierIndex, tier := range tier {
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

			b.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%s\t%s\n",
				tierID, idp, p.Name, node, p.Description, deps))
		}
	}

	return b.String(), nil
}
