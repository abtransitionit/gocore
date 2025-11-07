package phase2

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute a workflow
func (wkf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, retainSkipRange string, logger logx.Logger) error {
	// log
	logger.Infof("🅦 Runing workflow %q to %s", wkf.Name, wkf.Description)
	// logger.Info("• Phases in the same tier run concurrently. Next tier starts when the previous one completes.")
	// logger.Info("• Phases run concurently on each node.")
	logger.Info("• Tier concurrency:    all phases in tier run at the same time. Next tier starts when the previous one completes")
	logger.Info("• Phase concurrency:   each phase can run on multiple nodes at the same time")

	// get the tiers
	tierList, err := wkf.TopoSortByTier(logger)
	if err != nil {
		return fmt.Errorf("cannot sort tiers: %w", err)
	}

	// filter the phases in the tiers
	tierListFiltered, err := wkf.filterPhase(logger, tierList, retainSkipRange)
	if err != nil {
		return err
	}

	// display the workflow
	phaseView, err := wkf.GetTierView(tierListFiltered, logger)
	if err != nil {
		return fmt.Errorf("getting phase table: %w", err)
	}
	list.PrettyPrintTable(phaseView)

	// loop over each tier
	nbTier := len(tierListFiltered)
	for tierId, tier := range tierListFiltered {
		tierIdx := tierId + 1
		nbPhase := len(tier)
		logger.Infof("👉 Starting Tier %d / %d with %d concurrent phase(s)", tierIdx, nbTier, nbPhase)

		var wgTier sync.WaitGroup

		for _, p := range tier {
			wgTier.Add(1)

			go func(ph Phase) {
				defer wgTier.Done()

				// Determine target nodes
				nodes := resolveNode(ph.Node, cfg)
				if len(nodes) == 0 {
					logger.Warnf("   ↪ phase %q has no nodes resolved, skipping", ph.Name)
					return
				}

				// Determine function parameters
				param := resolveParam(ph.Param, cfg, logger)
				if len(nodes) == 0 {
					logger.Warnf("   ↪ phase %q has no nodes resolved, skipping", ph.Name)
					return
				}

				// log
				logger.Debugf("   ↪ phase %q > fn: %s", ph.Name, ph.FnAlias)
				if param != "" {
					logger.Debugf("   ↪ phase %q > param: %v (%s)", ph.Name, ph.Param, param)
				}

				var wgNodes sync.WaitGroup
				for _, node := range nodes {
					wgNodes.Add(1)
					go func(n string) {
						defer wgNodes.Done()
						logger.Debugf("   ↪ running phase %q > node: %s", ph.Name, node)
						// logger.Debugf("   ↪ running phase %q > node: %s (fn=%s, paramList=%v)", ph.Name, ph.FnAlias, ph.Param)
						// Here, call your actual execution function
						// runPhaseOnNode(ph, n)
						// - As an example the function loop over a list or map of item, and for each item, thethe function does:
						//    - actions locally
						//    - connect to ssh on the node and play a CLI
					}(node)
				}

				wgNodes.Wait()
			}(p)
		}

		wgTier.Wait()
		logger.Infof("✔ Tier %d complete. Waiting for next tier...", tierIdx)
	} // tier loop

	// success
	return nil
}

func resolveNode(phaseNode string, cfg *viperx.Viperx) []string {
	if cfg == nil || phaseNode == "" {
		return nil
	}
	return cfg.GetStringSlice(phaseNode)
}

func resolveParam(phaseParam []string, cfg *viperx.Viperx, logger logx.Logger) string {
	if cfg == nil || len(phaseParam) == 0 {
		return ""
	}

	var resolved []string
	for _, key := range phaseParam {
		// logger.Debugf("   ↪ resolving param: %s", key)
		value := cfg.Get(key)
		if value == nil {
			return ""
		}

		switch v := value.(type) {
		case string:
			resolved = append(resolved, v)
		case []interface{}:
			for _, item := range v {
				resolved = append(resolved, fmt.Sprintf("%v", item))
			}
		default:
			resolved = append(resolved, fmt.Sprintf("%v", v))
		}
	}

	return strings.Join(resolved, " ")
}

// var wgTier sync.WaitGroup
// for _, phase := range tier {
//     wgTier.Add(1)
//     go func(p Phase) {
//         defer wgTier.Done()
//         nodes := resolveNode(p.Node, cfg) // could return multiple nodes
//         var wgNodes sync.WaitGroup
//         for _, n := range nodes {
//             wgNodes.Add(1)
//             go func(node string) {
//                 defer wgNodes.Done()
//                 runPhaseOnNode(p, node)
//             }(n)
//         }
//         wgNodes.Wait()
//     }(phase)
// }
// wgTier.Wait()
