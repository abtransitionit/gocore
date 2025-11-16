package phase2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: manage the execution of a phase
func (phase *Phase) run(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, logger logx.Logger) error {

	// 1 - get target
	targetList, err := getTargetList(phase.Node, cfg)
	if err != nil {
		logger.Warnf("skipping phase %s, err: %v, ", phase.Name, err)
		// return fmt.Errorf("skipping phase %s, err: %v, ", phase.Name, err)
		return nil
	}

	// 2 - get parameter
	paramList, err := getParamList(phase.Param, cfg, logger)
	if err != nil {
		logger.Warnf("skipping phase %s, err: %v", phase.Name, err)
		return nil
		// return fmt.Errorf("skipping phase %s, err: %v", phase.Name, err)
	}

	// 3 - get PhaseFn
	phaseFn, err := getPhaseFn(phase.WkfName, phase.FnAlias, fnRegistry)
	if err != nil {
		logger.Warnf("skipping phase %s, err: %v", phase.Name, err)
		// return fmt.Errorf("skipping phase %s, err: %v", phase.Name, err)
		return nil
	}
	// 4 - get PhaseFn package and name
	goFnPkg, goFnName := describeFn(phaseFn, logger)

	// log
	logger.Debugf("↪ phase: %s > target:  %s > %v", phase.Name, phase.Node, targetList)
	logger.Debugf("↪ phase: %s > fnAlias: %s > %s/%s", phase.Name, phase.FnAlias, goFnPkg, goFnName)
	logger.Debugf("↪ phase: %s > param:   %s > %v", phase.Name, phase.Param, paramList)

	// 5 - create a GoFunc instance
	goFunction := &GoFunction{
		PhaseName: phase.Name,
		Name:      goFnName,
		ParamList: paramList,
		Func:      phaseFn,
	}

	// 6 - manage goroutines concurrency
	// nbTarget := len(targetList)
	var wgPhase sync.WaitGroup                      // (define) a WaitGroup instance for each tier : wait for all (concurent) goroutines (one per target) to complete
	errChPhase := make(chan error, len(targetList)) // channel to collect goroutines errors (if any or nil)
	// 61 - loop over each target of the phase AND create as many goroutines as targets
	// 61 - some goroutines will do SSH to play CLI remotely, other don't SSH and just play CLI locally
	for _, target := range targetList {
		wgPhase.Add(1)              // Increment the WaitGroup:counter for each target
		go func(oneTarget string) { // create as goroutine (that will run concurrently) as target in the phase AND pass it the target as an argument
			defer wgPhase.Done()                                        // Increment the WaitGroup:counter - when the goroutine (on the target) completes
			grErr := goFunction.run(ctx, phase.Name, oneTarget, logger) // delegate the execution of the function to this method
			if grErr != nil {                                           // send goroutines error if any into the chanel
				// log
				logger.Errorf("(goroutine) phase : %s >  %v", phase.Name, grErr)
				// send goroutines error if any into the chanel
				errChPhase <- grErr
			}
		}(target) // pass the target to the goroutine
	} // target loop

	// 62 - Synchronisation point: Wait for all goroutines (one per target) to finish/complete - done with the help of the WaitGroup:counter
	wgPhase.Wait()
	close(errChPhase)

	// 63 - collect errors
	var errList []error
	for e := range errChPhase {
		errList = append(errList, e)
	}

	// 64 - handle errors
	nbGroutineFailed := len(errList)
	errCombined := errors.Join(errList...)
	if nbGroutineFailed > 0 {
		logger.Errorf("❌ phase : %s > nb target that failed: %d", phase.Name, nbGroutineFailed)
		return errCombined
	}

	// 7 - handle success
	logger.Infof("✅ phase: %s (function: %s) > completes succesfully", phase.Name, goFunction.Name)
	return nil
}

// Description: resolves target node (ie. on which to run the function)
func getTargetList(phaseNode string, cfg *viperx.Viperx) ([]string, error) {
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
func getParamList(phaseParam []string, cfg *viperx.Viperx, logger logx.Logger) ([]string, error) {
	if cfg == nil || len(phaseParam) == 0 {
		return nil, fmt.Errorf("looking up > param > %q > cfg or phaseParam is empty", phaseParam)
	}

	var resolved []string

	for _, key := range phaseParam {
		value := cfg.Get(key)
		if value == nil {
			logger.Warnf("   ↪ looking up > param > %q not found in config", key)
			resolved = append(resolved, "")
			continue
		}

		switch v := value.(type) {
		case string:
			resolved = append(resolved, v)
		case []interface{}:
			for _, item := range v {
				resolved = append(resolved, fmt.Sprint(item))
			}
		case map[string]interface{}:
			b, _ := json.Marshal(v) // stable JSON
			resolved = append(resolved, string(b))
		default:
			resolved = append(resolved, fmt.Sprint(v))
		}
	}

	return resolved, nil
}

// return strings.TrimSpace(strings.Join(resolved, " ")), nil

func getPhaseFn(workflowName, fnAlias string, fnRegistry *FnRegistry) (PhaseFn, error) {
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
func describeFn(phaseFn PhaseFn, logger logx.Logger) (pkg string, name string) {
	if phaseFn == nil {
		logger.Warnf("describeFn: nil function")
		return "?", "?"
	}

	// PC = program counter
	pc := reflect.ValueOf(phaseFn).Pointer()
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

// logger.Debugf("   ↪ running phase > %s > node: %s (fn=%s, paramList=%v)", onePhase.Name, onePhase.FnAlias, onePhase.Param)
// Here, call your actual execution function
// runPhaseOnNode(ph, n)
// - As an example the function loop over a list or map of item, and for each item, thethe function does:
//    - actions locally
//    - connect to ssh on the node and play a CLI

// loggerstrings.Join(nodes, ",")
// if param != "" {
// 	logger.Debugf("   ↪ phase > %s > param: %v (%s)", onePhase.Name, onePhase.Param, param)
// }

// 	// // 2 - create a GoFunc instance
// 	// goFunction := &GoFunction{
// 	// 	Name:      goFnName,
// 	// 	ParamList: paramList,
// 	// 	Func:      phaseFn,
// 	// }

// 	return nil
// }

// errChReader := make(chan struct{}) // (define) a channel to read errors from the error chanel
// go func() {                        // (create) a goroutine to collect errors
// 	for e := range errChTarget {
// 		logger.Errorf("reader %v", e)
// 		ErrList = append(ErrList, e)
// 	}
// 	close(errChReader) // signal that no more errors will be sent
// }() // no parameter passed

// 5 - collect goroutines errors
// var ErrList []error // collect goroutines errors

// errChTarget := make(chan error, nbTarget) // (define) a channel to collect errors from goroutines
// close(errChTarget) // close the channel that collect error - signal that no more errors will be sent
// <-errChReader      // mean wait to receive a value from the channel named "done” -  wait until the error collector finishes processing
// errChTarget <- fmt.Errorf("phase: %s > target: %s >  function: %s > %w", phase.Name, oneTarget, goFunction.Name, err)
// logger.Errorf("hello %v", err)

// 7 - manage goroutines error
// 71 - collect goroutines errors

// for e := range errChTarget {
// 	logger.Errorf("hello %v", e)
// 	ErrList = append(ErrList, e)
// }

// // 72 - handle errors
// if len(ErrList) > 0 {
// 	combinedErr := errors.Join(ErrList...)
// 	return fmt.Errorf("%w", combinedErr)
// }
