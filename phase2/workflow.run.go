package phase2

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute a workflow
//
// Parameters:
//   - ctx
//   - cfg : the workflow config file (ie. a strunct that contains all the vars used by the workflow)
//   - fnRegistry : the registry of functions used by the workflow
//   - retainSkipRange : the range of phases to retain/skip
//   - logger
//
// Return:
//   - error
//
// Example Usage:
//
//	err := workflow.Execute(ctx, cfg, fnRegistry, retainSkipRange, logger)
//
// Notes:
//
// - run the workflow
// - log the workflow
func (wkf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, retainSkipRange string, logger logx.Logger) error {
	// TODO: Before running the workflow:
	//   - check syntax then/and	/or do:
	//     - check fn    is resolved in the config
	//     - check param is resolved in the config
	//     - check node  is resolved in the config
	//     - Then
	//       - check fn is registred

	// log
	logger.Infof("🅦 Runing workflow %q to %s", wkf.Name, wkf.Description)
	logger.Info("• Tier concurrency:    all phases in tier run at the same time. Next tier starts when the previous one completes")
	logger.Info("• Phase concurrency:   each phase runs at the same time on multiple nodes")

	// 1 - get the tiers
	tierList, err := wkf.TopoSortByTier(logger)
	if err != nil {
		return fmt.Errorf("cannot sort tiers: %w", err)
	}

	// 2 - filter the tier phases according to retainSkipRange
	tierListFiltered, err := wkf.filterPhase(logger, tierList, retainSkipRange)
	if err != nil {
		return err
	}

	// 3 - display the workflow
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

		// log
		logger.Infof("👉 Starting Tier %d:%d:%d concurrent phase(s)", tierIdx, nbTier, nbPhase)

		var wgTier sync.WaitGroup // Creates a WaitGroup instance - will wait for all (concurent) goroutines (one per phase) to complete
		// loop over each phase
		for _, p := range tier {
			wgTier.Add(1) // Increment the WaitGroup counter
			// define/start a goroutine for each phase
			go func(ph Phase) {
				defer wgTier.Done() // Decrement the WaitGroup counter - when the goroutine (the phase) completes

				// 1 - for this phase
				// 11 - resolve the target nodes
				targetList, err := resolveNode(ph.Node, cfg)
				if err != nil {
					logger.Warnf("skipping phase %s, err: %v, ", ph.Name, err)
					return
				}

				// 12 - resolve function parameters
				paramList, err := resolveParam(ph.Param, cfg, logger)
				if err != nil {
					logger.Warnf("skipping phase %s, err: %v", ph.Name, err)
					return
				}

				// 13 - resolve function name
				goFn, err := resolveFn(wkf.Name, ph.FnAlias, fnRegistry)
				if err != nil {
					logger.Warnf("skipping phase %s, err: %v", ph.Name, err)
					return
				}
				// 14 - for this function of this phase: resolve package and name
				goFnPkg, goFnName := describeFn(goFn, logger)

				// 15 - log all this info
				logger.Debugf("   ↪ phase %q > lookup > target   > %s > %v", ph.Name, ph.Param, targetList)
				logger.Debugf("   ↪ phase %q > lookup > function > %s > %s/%s", ph.Name, ph.FnAlias, goFnPkg, goFnName)
				logger.Debugf("   ↪ phase %q > lookup > param    > %s > %v", ph.Name, ph.Param, paramList)
				// loggerstrings.Join(nodes, ",")
				// if param != "" {
				// 	logger.Debugf("   ↪ phase %q > param: %v (%s)", ph.Name, ph.Param, param)
				// }

				var wgNodes sync.WaitGroup // Creates a SECOND WaitGroup instance - will wait for all (concurent) goroutines (one per node) to complete
				for _, target := range targetList {
					wgNodes.Add(1) // Increment the SECOND WaitGroup counter
					go func(n string) {
						defer wgNodes.Done() // Decrement the SECOND WaitGroup counter - when the goroutine (on the node) completes
						logger.Debugf("   ↪ phase %q > RUNNINING function on target > %s", ph.Name, target)
						// Todo
					}(target)
				}

				wgNodes.Wait() // Wait for all goroutines (launched on node) to complete
			}(p)
		}

		wgTier.Wait() // Wait for all goroutines (launched by phase) to complete
		logger.Infof("✔ Tier %d complete.", tierIdx)
	} // tier loop

	// success
	return nil
}

// Description: resolves target node (ie. on which to run the function)
func resolveNode(phaseNode string, cfg *viperx.Viperx) ([]string, error) {
	if cfg == nil || phaseNode == "" {
		return nil, fmt.Errorf("looking up > target > %q > cfg or phaseNode is empty", phaseNode)
	}

	nodes := cfg.GetStringSlice(phaseNode)
	if len(nodes) == 0 {
		// logger.Warnf("   ↪ looking up> target > %q > not found or empty in config", phaseNode)
		return nil, fmt.Errorf("looking up > target > %q > not found or empty in config", phaseNode)
	}
	// success
	return nodes, nil
}
func resolveParam(phaseParam []string, cfg *viperx.Viperx, logger logx.Logger) ([]string, error) {
	if cfg == nil || len(phaseParam) == 0 {
		return nil, fmt.Errorf("looking up > param > %q > cfg or phaseParam is empty", phaseParam)
	}

	resolved := make([]string, 0, len(phaseParam))

	for _, key := range phaseParam {
		value := cfg.Get(key)
		if value == nil {
			logger.Warnf("   ↪ looking up > param > %q not found in config", key)
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
	}

	// return strings.TrimSpace(strings.Join(resolved, " ")), nil
	return resolved, nil

}

func resolveFn(workflowName, fnAlias string, fnRegistry *FnRegistry) (PhaseFn, error) {
	if fnAlias == "" {
		return nil, fmt.Errorf("fn alias is empty")
	}

	fn, ok := fnRegistry.Get(workflowName, fnAlias)
	if !ok {
		return nil, fmt.Errorf("getting registred function for alias %s:%s (not registered)", workflowName, fnAlias)
	}

	return fn, nil
}

// Description: returns package import path + function name for a PhaseFn
func describeFn(fn PhaseFn, logger logx.Logger) (pkg string, name string) {
	if fn == nil {
		logger.Warnf("describeFn: nil function")
		return "?", "?"
	}

	// PC = program counter
	pc := reflect.ValueOf(fn).Pointer()
	fnInfo := runtime.FuncForPC(pc)
	if fnInfo == nil {
		logger.Warnf("retrieving fnInfo")
		return "?", "?"
	}

	// Basic info from runtime.Func
	full := fnInfo.Name() // e.g. "github.com/me/project/mock/node.CheckSshConf"
	// file, line := fnInfo.FileLine(pc)

	// Log what we know
	// logger.Infof("Function: name=%s", full)
	// logger.Infof("Function: file=%s:%d", file, line)

	// Extract module/package (best effort)
	// Example full: github.com/me/project/mock/node.CheckSshConf
	// module path = everything before the last slash
	module := full[:strings.LastIndex(full, "/")]
	moduleX := strings.TrimPrefix(module, "github.com/abtransitionit/")
	// logger.Infof("Function: module=%s", moduleX)

	// Extract only last section: mock/node.CheckSshConf
	last := full[strings.LastIndex(full, "/")+1:]

	// Split on dots
	parts := strings.Split(last, ".")
	// ["mock/node", "CheckSshConf"]
	// Or ["mock/node", "CheckSshConf", "CheckSshConf"]

	// Package = first part
	pkg = parts[0]

	// Name = last part
	name = parts[len(parts)-1]
	fullName := fmt.Sprintf("%s.%s", pkg, name)

	// logger.Infof("Function: package=%s", pkg)
	// logger.Infof("Function: fn=%s", name)

	return moduleX, fullName
}

// logger.Debugf("   ↪ running phase %q > node: %s (fn=%s, paramList=%v)", ph.Name, ph.FnAlias, ph.Param)
// Here, call your actual execution function
// runPhaseOnNode(ph, n)
// - As an example the function loop over a list or map of item, and for each item, thethe function does:
//    - actions locally
//    - connect to ssh on the node and play a CLI
