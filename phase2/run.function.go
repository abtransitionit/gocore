// File in gocore/phase/adapter.go
package phase2

import (
	"context"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
)

// description: fetch a function in the registry and execute it
func (gf *GoFunc) Execute(ctx context.Context, fr *FunctionRegistry, nodeList []string, logger logx.Logger) (string, error) {
	var goFunc *GoFunc

	// resolve
	goFunc, err := resolveFunction(gf.PhaseFuncName, fr, logger)
	if err != nil {
		return "", err
	}

	logger.Infof("🅕 Executing function: %s concurently on nodes %v", goFunc.PhaseFuncName, nodeList)
	// loop over each node - the function is executed on each node concurrently
	errCh := make(chan error, len(nodeList)) // channel to collect as many errors as nodes - one per node
	var wg sync.WaitGroup                    // Creates a WaitGroup instance - will wait for all goroutines to complete
	for _, node := range nodeList {
		wg.Add(1)    // increment the WaitGroup counter - signal that a new goroutine is starting
		node := node // capture loop variable
		logger.Infof("🅝 Node %s > Starting function", node)
		go func() { // launch a goroutines per node - goroutine works concurrently
			defer wg.Done()                                  // decrement the WaitGroup counter - when the goroutine completes
			if err := goFunc.Func(ctx, logger); err != nil { // execute this code
				errCh <- fmt.Errorf("🅝 Node %s > function failed: %w", node, err) // way to get all goroutine's error (nil or error)
			}
		}()
	} // node loop
	wg.Wait()    // Wait for all phases in this tier to completes
	close(errCh) // close the channel - signal that no more error will be sent

	// loop over the channel to log any error reported by the goroutines
	var NodeErrs []error
	for e := range errCh {
		logger.Errorf(e.Error())
		NodeErrs = append(NodeErrs, e)
	}

	// if an error occurred on a node, return
	if len(NodeErrs) > 0 {
		return "", fmt.Errorf("errors occurred on a %d node", len(NodeErrs))
	}

	// success
	return "ok", nil
}
