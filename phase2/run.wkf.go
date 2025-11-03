// File in gocore/phase/adapter.go
package phase2

import (
	"context"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

func (wkf *Workflow) Execute(ctx context.Context, config *viperx.CViper, fr *FunctionRegistry, logger logx.Logger) error {
	//log
	logger.Infof("🅦 Starting workflow: %s > %s", wkf.Name, wkf.Description)
	logger.Info("Phases in the same tier run in parallel. Next tier starts when the previous one completes")

	// get tiers - that is a set of phases
	tiers, err := wkf.TopoTierSorted()
	if err != nil {
		return fmt.Errorf("sorting tier: %w", err)
	}

	// loop over each tier - run each tier sequentially - and all phases of a tier concurrently
	nbTier := len(tiers)
	for i, tier := range tiers {
		tierIdx := i + 1
		nbPhase := len(tier)
		logger.Infof("🅣 Starting Tier %d/%d with %d concurent phase(s)", tierIdx, nbTier, nbPhase)

		var wg sync.WaitGroup              // Creates a WaitGroup instance - will wait for all goroutines to complete
		errCh := make(chan error, nbPhase) // channel to collect as many errors as phases - one per phase

		for _, phase := range tier { // loop over each phase of the tier
			wg.Add(1)          // increment the WaitGroup counter - signal that a new goroutine is starting
			go func(p Phase) { // launch a goroutines per phase - goroutine works concurrently
				defer wg.Done()                                               // decrement the WaitGroup counter - when the goroutine completes
				if _, err := p.Execute(ctx, config, fr, logger); err != nil { // execute this code
					errCh <- fmt.Errorf("phase %s failed: %w", p.Name, err) // way to get all goroutine's error (nil or error)
				}
			}(phase)
		} // phase loop

		wg.Wait()    // Wait for all phases in this tier to completes
		close(errCh) // close the channel - signal that no more error will be sent

		// loop over the channel to log any error reported by the goroutines
		var tierErrs []error
		for e := range errCh {
			logger.Errorf(e.Error())
			tierErrs = append(tierErrs, e)
		}

		// if an error occurred in this tier, return
		if len(tierErrs) > 0 {
			return fmt.Errorf("errors occurred in tier %d", tierIdx)
		}

		logger.Infof("👉 Completed Tier %d", tierIdx)
	} // tier loop

	// success
	logger.Info("• Workflow completed successfully")
	return nil
}
