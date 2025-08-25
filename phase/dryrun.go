// File in gocore/phase/run.go
package phase

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Name: DryRun
//
// Description: performs a dry-run of the workflow by displaying the sorted execution plan.
//
// Parameters:
//
//	ctx:    The context for the phase.
//	logger: The logger to use for printing messages.
//
// Notes:
//
//	This method does not execute any phases. It is intended to be used for
//	planning and verification purposes by a CLI or other external system.
//
// TODO: Extend this method to include a more realistic dry run by executing
//
//	mock phases that log their execution without performing any real work.
func (w *Workflow) DryRun(ctx context.Context, logger logx.Logger, skipPhases []int, retainPhases []int) error {
	logger.Info("Dry Running workflow...")

	// check parameters: skipPhases  and retainPhases are mutually exclusive
	if len(skipPhases) > 0 && len(retainPhases) > 0 {
		return fmt.Errorf("invalid parameters: skipPhases and retainPhases cannot be set at the same time")
	}

	// display the wokflow (ie. the list of phases)
	logger.Info("The worflow phases are:")
	w.Show(logger)

	// Log the received IDs to skipping or retaining if any
	if len(skipPhases) > 0 {
		logger.Infof("Received phase IDs to skip: %v", skipPhases)
	} else if len(retainPhases) > 0 {
		logger.Infof("Received phase IDs to retain: %v", retainPhases)
	}

	// Get the Tier (ie. topo dependency sorted phases)
	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	// get filtered phases
	filteredTiers, err := w.filterPhase(logger, sortedTiers, skipPhases, retainPhases)
	if err != nil {
		return fmt.Errorf("failed to filter phases: %w", err)
	}

	// Show the filtered phases ordered by tier
	logger.Info("The execution plan is:")

	filteredTiers.Show(logger)

	logger.Info("Dry runned the worflow.")
	return nil
}
