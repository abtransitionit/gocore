package phase2

import (
	"context"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute/run a workflow
//
// Parameters:
//   - ctx
//   - cfg : the workflow config file as a struct that contains all the vars used by the workflow
//   - fnRegistry : contains a map of all GO functions used by the workflow
//   - retainSkipRange : manage the phases to retain/skip in the workflow
//   - logger
//
// Return:
//   - an error if something went wrong
//
// Example Usage:
//
//	err := workflow.Execute(ctx, cfg, &fnRegistry, retainSkipRange, logger)
//
// Notes:
//
// - the order of the phases is based on their dependency ==> topological sort => concept of tier
// - a tier is a ordered set of phases
// - a phase denotes a GO function that takes parameters
// - running a worflow consist of running:
//   - one tier after another
//   - each phase of a same tier concurently/asynchrous/at the same time
//   - a phase runs concurently/asynchrous/at the same time on (potentially) several targets (ie. VM, container, local)
//
// Todo:
//  - check syntax of the config yaml
//  - check syntax of the workflow yaml
//  - check fn    is resolved in the config
//  - check param is resolved in the config
//  - check node  is resolved in the config
//  - check fn  is registred

func (wkf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, retainSkipRange string, logger logx.Logger) error {

	// log
	logger.Infof("🅦 Runing workflow %q to %s", wkf.Name, wkf.Description)
	logger.Info("• Tier sequencing:    next tier starts after the previous one completes")
	logger.Info("• Tier completion:    a tier completes when all its phases completes")
	logger.Info("• Host concurrency:   each phase rus concurently on all hosts")
	logger.Info("• Node concurrency:   each phase runs (also) concurently on all nodes (when defined)")
	logger.Info("• Phase completion:   a phase completes (for a host) when all its subsequent node tasks complete")

	// 1 - get the tiers
	tierList, err := wkf.TopoSortByTier(logger)
	if err != nil {
		return fmt.Errorf("cannot sort tiers: %w", err)
	}

	// 2 - filter the tier phases according to retainSkipRange
	tierListFiltered, err := wkf.filterPhase(tierList, retainSkipRange, logger)
	if err != nil {
		return err
	}

	// 3 - display the workflow
	phaseView, err := wkf.GetTierView(tierListFiltered, logger)
	if err != nil {
		return fmt.Errorf("getting phase table: %w", err)
	}
	list.PrettyPrintTable(phaseView)

	// 4 - loop over each tier
	nbTier := len(tierListFiltered)
	for tierId, phaseList := range tierListFiltered {

		tierIdx := tierId + 1
		nbPhase := len(phaseList)

		// log
		logger.Infof("👉 Starting Tier %d:%d:%d concurrent phase(s)", tierIdx, nbTier, nbPhase)

		// 6 - manage goroutines concurrency
		var wgTier sync.WaitGroup               // define a WaitGroup instance for each tier : wait for all (concurent) goroutines (one per phase in a tier) to complete
		errChPhase := make(chan error, nbPhase) // define a channel to collect errors from goroutines

		// 61 - loop over each phases in the tier AND create as many goroutines as phases
		for _, phase := range phaseList {
			wgTier.Add(1)             // Increment the WaitGroup:counter for each phase
			go func(onePhase Phase) { // create as goroutine (that will run concurrently) as phase in the tier AND pass it the phase as an argument
				defer wgTier.Done()                            // Decrement the WaitGroup counter - when the goroutine (the phase) completes
				err := phase.run(ctx, cfg, fnRegistry, logger) // delegate the execution of the phase to this method
				if err != nil {                                // send goroutines error if any into the chanel
					errChPhase <- fmt.Errorf("%w", err)
				}
			}(phase) // pass the phase to the goroutine
		} // phase loop

		wgTier.Wait()     // Wait for all goroutines (one per phase) to complete - done with the help of the WaitGroup:counter
		close(errChPhase) // close the channel - signal that no more error will be sent

		// 7 - manage goroutines error
		// 71 - Aggregate goroutines errors
		var ErrList []error
		for e := range errChPhase {
			ErrList = append(ErrList, e)
		}

		// 72 - handle error
		if len(ErrList) > 0 {
			return fmt.Errorf("errors occurred in tier %d", tierIdx)
		}

		// 8 - handle success
		logger.Infof("• Tier %d complete.", tierIdx)
	} // tier loop

	// success
	return nil
}
