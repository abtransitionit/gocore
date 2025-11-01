// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

func (wf *Workflow) Execute(config *viperx.CViper, logger logx.Logger) error {
	//log
	logger.Infof("🅦 Starting workflow: %s", wf.Name)
	logger.Info("Phases in the same tier run in parallel. Next tier starts when the previous one completes")

	// get tiers
	tiers, err := wf.TopoTierSorted()
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
				if _, err := p.Execute(logger); err != nil {
					errCh <- fmt.Errorf("phase %s failed: %w", p.Name, err)
				}
			}(phase)
		}

		// Wait for all phases in this tier
		wg.Wait()
		close(errCh)

		if len(errCh) > 0 {
			for e := range errCh {
				logger.Errorf(e.Error())
			}
			return fmt.Errorf("occurring in tier %d", tierIdx)
		}

		logger.Infof("👉 Completed Tier %d", tierIdx)
	}

	// success
	logger.Info("• Workflow completed successfully")
	return nil
}

func (wf *Workflow) Execute2(cfg *viperx.CViper, logger logx.Logger) error {
	// log
	logger.Infof("• Starting workflow: %s", wf.Name)
	logger.Info("• Phases in the same tier run in parallel. Next tier starts when the previous one completes")

	// toposort the phases of the workflow
	phases, _ := wf.TopoPhaseSorted()

	// Loop over sorted phases
	for _, phase := range phases {
		// Execute the phase
		phase.Execute(logger)
	}
	return nil
}
