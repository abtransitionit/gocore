// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

func (wkf *Workflow) Execute(config *viperx.CViper, fr *FunctionRegistry, logger logx.Logger) error {
	//log
	logger.Infof("🅦 Starting workflow: %s > %s", wkf.Name, wkf.Description)
	logger.Info("Phases in the same tier run in parallel. Next tier starts when the previous one completes")

	// get tiers
	tiers, err := wkf.TopoTierSorted()
	nbTier := len(tiers)
	if err != nil {
		return fmt.Errorf("sorting tier: %w", err)
	}

	// loop over each tier - run each tier sequentially
	for i, tier := range tiers {
		tierIdx := i + 1
		nbPhase := len(tier)
		logger.Infof("🅣 Starting Tier %d/%d with %d concurent phase(s)", tierIdx, nbTier, nbPhase)

		var wg sync.WaitGroup
		errCh := make(chan error, nbPhase)

		// loop over each phase of the tier- run all phases concurrently
		for _, phase := range tier {
			wg.Add(1)
			go func(p Phase) {
				defer wg.Done()
				if _, err := p.Execute(config, fr, logger); err != nil {
					errCh <- fmt.Errorf("phase %s failed: %w", p.Name, err)
				}
			}(phase)
		}

		// Wait for all phases in this tier to completes
		wg.Wait()
		close(errCh)

		// Collect all errors
		var tierErrs []error
		for e := range errCh {
			logger.Errorf(e.Error())
			tierErrs = append(tierErrs, e)
		}
		if len(tierErrs) > 0 {
			return fmt.Errorf("errors occurred in tier %d", tierIdx)
		}

		// if len(errCh) > 0 {
		// 	for e := range errCh {
		// 		logger.Errorf(e.Error())
		// 	}
		// 	return fmt.Errorf("occurring in tier %d", tierIdx)
		// }

		logger.Infof("👉 Completed Tier %d", tierIdx)
	}

	// success
	logger.Info("• Workflow completed successfully")
	return nil
}
