package phase2

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute a workflow
func (wkf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, retainSkipRange string, logger logx.Logger) error {
	// TODO: Before running the workflow:
	//   - check syntax then/and/or do:
	//     - check fn is resolved in the config
	//     - check param is resolved in the config
	//     - check node is resolved in the config
	//     - Then
	//       - check fn is registred

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
		// logger.Infof("👉 Starting Tier %d of %d with %d concurrent phase(s)", tierIdx, nbTier, nbPhase)
		logger.Infof("👉 Starting Tier %d:%d:%d concurrent phase(s)", tierIdx, nbTier, nbPhase)

		var wgTier sync.WaitGroup

		for _, p := range tier {
			wgTier.Add(1)

			go func(ph Phase) {
				defer wgTier.Done()

				// resolve target nodes
				nodes := resolveNode(ph.Node, cfg, logger)
				if len(nodes) == 0 {
					logger.Warnf("   ↪ phase %q has no nodes resolved, skipping", ph.Name)
					return
				}

				// resolve function parameters
				// resolveParam(ph.Param, cfg, logger)
				resolveParam(ph.Param, cfg, logger)
				// paramNodes := strings.Split(paramStr, ",") // convert to slice
				// logger.Debugf("   ↪ paramNodes %q ", paramNodes)
				// if len(paramNodes) > 0 {
				// 	nodes = paramNodes // override execution nodes with param nodes
				// }

				// log
				logger.Debugf("   ↪ phase %q > fn: %s", ph.Name, ph.FnAlias)
				logger.Debugf("   ↪ phase %q > nodes: %v", ph.Name, nodes)
				// logger.Debugf("   ↪ phase %q > fn: %s", ph.Name, ph.FnAlias)
				// if param != "" {
				// 	logger.Debugf("   ↪ phase %q > param: %v (%s)", ph.Name, ph.Param, param)
				// }

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
		logger.Infof("✔ Tier %d complete.", tierIdx)
	} // tier loop

	// success
	return nil
}

func resolveNode(phaseNode string, cfg *viperx.Viperx, logger logx.Logger) []string {
	if cfg == nil || phaseNode == "" {
		return nil
	}

	nodes := cfg.GetStringSlice(phaseNode)
	if len(nodes) == 0 {
		logger.Warnf("   ↪ node key %q not found or empty in config", phaseNode)
		return nil
	}

	// Join for logging like "o1u,o2a,o3r"
	logger.Debugf("   ↪ lookup > target > %s > %s", phaseNode, strings.Join(nodes, ","))

	return nodes
}
func resolveParam(phaseParam []string, cfg *viperx.Viperx, logger logx.Logger) string {
	if cfg == nil || len(phaseParam) == 0 {
		return ""
	}

	resolved := make([]string, 0, len(phaseParam))

	for _, key := range phaseParam {
		value := cfg.Get(key)
		if value == nil {
			logger.Warnf("   ↪ param key %q not found in config", key)
			resolved = append(resolved, "")
			continue
		}

		// Lightweight conversion
		var str string
		switch v := value.(type) {
		case string:
			str = v
		case []interface{}:
			parts := make([]string, len(v))
			for i, item := range v {
				parts[i] = fmt.Sprint(item)
			}
			str = strings.Join(parts, ",")
		case map[string]interface{}:
			b, _ := json.Marshal(v) // stable JSON
			str = string(b)
		default:
			str = fmt.Sprint(v)
		}

		resolved = append(resolved, str)
		logger.Debugf("   ↪ lookup > param > %s > %s", key, str)
	}

	return strings.TrimSpace(strings.Join(resolved, " "))
}
