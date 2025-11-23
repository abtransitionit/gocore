package phase2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: resolves host node (ie. on which to run the function)
func getHostList(phaseNode string, cfg *viperx.Viperx) ([]string, error) {
	if cfg == nil || phaseNode == "" {
		return nil, fmt.Errorf("looking up > host > %q > cfg or phaseNode is empty", phaseNode)
	}

	nodes := cfg.GetStringSlice(phaseNode)
	if len(nodes) == 0 {
		// logger.Warnf("   ↪ looking up> host > %q > not found or empty in config", phaseNode)
		return nil, fmt.Errorf("looking up > host > %q > not found or empty in config", phaseNode)
	}
	// success
	return nodes, nil
}

// Description: resolves phase parameters

func getParamList(phaseParam []string, cfg *viperx.Viperx, logger logx.Logger) ([][]any, error) {
	// check parameters
	if cfg == nil {
		return nil, fmt.Errorf("cfg is nil")
	}
	if len(phaseParam) == 0 {
		return [][]any{}, nil // return empty slice instead of error
	}

	// resolve
	resolved := make([][]any, len(phaseParam))
	for i, key := range phaseParam {
		val := cfg.Get(key)
		if val == nil {
			logger.Warnf("param %q not found", key)
			resolved[i] = []any{""}
			continue
		}

		switch v := val.(type) {
		case string:
			resolved[i] = []any{v}
		case []any:
			anySlice := make([]any, len(v))
			copy(anySlice, v)
			resolved[i] = anySlice
		case map[string]any:
			b, _ := json.Marshal(v)
			resolved[i] = []any{string(b)}
		default:
			resolved[i] = []any{fmt.Sprint(v)}
		}
	}

	return resolved, nil
}

func getParamList1(phaseParam []string, cfg *viperx.Viperx, logger logx.Logger) ([]string, error) {
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

// return strings.TrimSpace(strings.Join(resolved, " ")), nil

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
