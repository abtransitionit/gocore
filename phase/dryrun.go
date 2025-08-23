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
	logger.Info("Starting workflow planning (dry run)...")

	// check parameters: skipPhases  and retainPhases are mutually exclusive
	if len(skipPhases) > 0 && len(retainPhases) > 0 {
		return fmt.Errorf("invalid parameters: skipPhases and retainPhases cannot be set at the same time")
	}

	// display the wokflow (ie. the list of phases)
	logger.Info("The worflow contains the following phases:")
	w.Show(logger)

	// Logging the received IDs (skippeds or retained)
	logger.Infof("Received phase IDs to skip: %v", skipPhases)

	// Get the sorted phases
	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	// Get filtered phases
	filteredTiers, err := w.filterSkipPhase(sortedTiers, skipPhases)
	if err != nil {
		return fmt.Errorf("failed to filter phases: %w", err)
	}

	// Show the filtered phases ordered by tier
	logger.Info("The execution plan is:")

	filteredTiers.Show(logger)

	logger.Info("Workflow planning finished. No phases were executed.")
	return nil
}
