// File in gocore/phase/adapter.go
package phase2

import (
	"context"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

func (phase *Phase) Execute(ctx context.Context, config *viperx.CViper, fr *FunctionRegistry, logger logx.Logger) (string, error) {
	nodeList := resolveNode(phase.Node, config)
	paramMap := resolveParam(phase.Param, config)

	logger.Debugf("🅟 Starting Phase : %s > NodeSet: %v", phase.Name, nodeList)

	// Fetch the registered GoFunc
	goFunc, ok := fr.Get(phase.Fn)
	if !ok {
		return "", fmt.Errorf("function %q not registered", phase.Fn)
	}

	// Execute on all nodes concurrently
	errCh := make(chan error, len(nodeList))
	var wg sync.WaitGroup

	for _, node := range nodeList {
		wg.Add(1)
		node := node
		go func() {
			defer wg.Done()
			logger.Infof("🅝 Node %s > Starting PhaseFunction %s", node, goFunc.PhaseFuncName)

			// Pass parameters to the function
			if err := goFunc.Func(ctx, paramMap, logger); err != nil {
				errCh <- fmt.Errorf("🅝 Node %s > function failed: %w", node, err)
			}
		}()
	}

	wg.Wait()
	close(errCh)

	var nodeErrs []error
	for e := range errCh {
		logger.Errorf(e.Error())
		nodeErrs = append(nodeErrs, e)
	}

	if len(nodeErrs) > 0 {
		return "", fmt.Errorf("errors occurred on %d node(s)", len(nodeErrs))
	}

	return "ok", nil
}

// func (phase *Phase) Execute(ctx context.Context, config *viperx.CViper, fr *FunctionRegistry, logger logx.Logger) (string, error) {
// 	// resolve
// 	nodeList := resolveNode(phase.Node, config)
// 	paramMap := resolveParam(phase.Param, config)

// 	// log
// 	logger.Debugf("🅟 Starting Phase : %s > NodeSet:  %s (%v)", phase.Name, phase.Node, nodeList)

// 	// Execute the function
// 	_, err := GetGoFunc(phase.Fn).Execute(ctx, fr, nodeList, paramMap, logger)
// 	if err != nil {
// 		return "", err
// 	}

// 	// success
// 	return "", nil
// }
