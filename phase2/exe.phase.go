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

	// 1 - get target
	targetList, err := getTargetList(phase.Node, cfg)
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
	logger.Debugf("↪ phase: %s > target:  %s > %v", phase.Name, phase.Node, targetList)
	logger.Debugf("↪ phase: %s > fnAlias: %s > %s/%s", phase.Name, phase.FnAlias, goFnPkg, goFnName)
	if len(phase.Param) > 0 {
		for i, key := range phase.Param {
			if i < len(paramList) {
				logger.Debugf("↪ phase: %s > param: %s > %v", phase.Name, key, paramList[i])
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
	// nbTarget := len(targetList)
	var wgPhase sync.WaitGroup                      // (define) a WaitGroup instance for each tier : wait for all (concurent) goroutines (one per target) to complete
	errChPhase := make(chan error, len(targetList)) // channel to collect goroutines errors (if any or nil)
	// 61 - loop over each target of the phase AND create as many goroutines as targets
	// 61 - some goroutines will do SSH to play CLI remotely, other don't SSH and just play CLI locally
	for _, target := range targetList {
		wgPhase.Add(1)              // Increment the WaitGroup:counter for each target
		go func(oneTarget string) { // create as goroutine (that will run concurrently) as target in the phase AND pass it the target as an argument
			defer wgPhase.Done()                                                // Increment the WaitGroup:counter - when the goroutine (on the target) completes
			grErr := goFunction.runOnTarget(ctx, phase.Name, oneTarget, logger) // delegate the execution of the function to this method
			if grErr != nil {                                                   // send goroutines error if any into the chanel
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
