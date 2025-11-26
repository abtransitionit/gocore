package phase2

import (
	"context"
	"errors"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: manage the execution of a phase
func (phase *Phase) run(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, logger logx.Logger) error {

	// 1 - get host
	hostList, err := getHostList(phase.Node, cfg)
	if err != nil {
		logger.Warnf("skipping phase %s, err: %v, ", phase.Name, err)
		return nil
	}

	// 2 - get parameter
	paramList, err := getParamList(phase.Param, cfg, logger)
	if err != nil {
		logger.Warnf("skipping phase %s, err: %v", phase.Name, err)
		return nil
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
	logger.Debugf("↪ %s > host:  %s > %v", phase.Name, phase.Node, hostList)
	logger.Debugf("↪ %s > fnAlias: %s > %s/%s", phase.Name, phase.FnAlias, goFnPkg, goFnName)
	if len(phase.Param) > 0 {
		for i, key := range phase.Param {
			if i < len(paramList) {
				logger.Debugf("↪ %s > param: %s > %v", phase.Name, key, paramList[i])
			}
		}
	}

	// 5 - create a GoFunc instance
	goFunction := &GoFunction{
		PhaseName: phase.Name,
		Name:      goFnName,
		ParamList: paramList,
		Func:      phaseFn,
	}
	// 6 - manage goroutines concurrency
	nbItem := len(hostList)
	var wgPhase sync.WaitGroup             // define a WaitGroup instance for each item in the list : wait for all (concurent) goroutines to complete
	errChPhase := make(chan error, nbItem) // define a channel to collect errors from each goroutine
	// log
	// 61 - loop over each host of the phase AND create as many goroutines as hosts
	// 61 - some goroutines will do SSH to play CLI remotely, other don't SSH and just play CLI locally
	for _, host := range hostList {
		wgPhase.Add(1)            // Increment the WaitGroup:counter for each host
		go func(oneItem string) { // create as many goroutine (that will run concurrently) as item AND pass the item as an argument
			defer func() {
				logger.Debugf("↩ (%s) > %s > complete", phase.Name, oneItem)
				wgPhase.Done() // Decrement the WaitGroup counter - when the goroutine complete
			}()
			logger.Debugf("↪ (%s) > %s > ongoing", phase.Name, oneItem)
			grErr := goFunction.runOnHOst(ctx, phase.Name, oneItem, logger) // delegate the execution of the function to this method
			if grErr != nil {                                               // send goroutines error if any into the chanel
				// log
				logger.Errorf("phase : %s >  %v", phase.Name, grErr)
				// send goroutines error if any into the chanel
				errChPhase <- grErr
			}
		}(host) // pass the host to the goroutine
	} // host loop

	// 62 - Synchronisation point: Wait for all goroutines (one per host) to finish/complete - done with the help of the WaitGroup:counter
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
		logger.Errorf("❌ %s > nb host that failed: %d", phase.Name, nbGroutineFailed)
		return errCombined
	}

	// 7 - handle success
	logger.Infof("✅ %s > completes", phase.Name)
	return nil
}
