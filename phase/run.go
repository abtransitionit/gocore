// File in gocore/phase/run.go
package phase

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/syncx"
)

// Name: Execute
//
// Description: executes the phases of a workflow
//
// Parameters:
//
//   - ctx: The context for the workflow.
//   - logger: The logger to use for printing messages.
//   - skipPhases: A slice of integer IDs representing the phases to be skipped.
//
// Returns:
//
//   - An error if the workflow fails to execute.
//
// Notes:
//
//   - executes all phases in a tier concurrently.
//   - executes each tier sequentially.
//   - the order of the tier and their phases is determined by the topological function
func (w *Workflow) Execute(ctx context.Context, logger logx.Logger, targets []Target, skipPhases []int, retainPhases []int) error {
	logger.Info("🚀 Starting workflow execution")

	// check some paramaters, deletgate other checks to the filterPhase function
	if len(w.Phases) == 0 {
		err := fmt.Errorf("workflow is empty: no phases are defined")
		logger.ErrorWithNoStack(err, "Cannot execute workflow")
		return err
	}

	// display the wokflow (ie. the list of phases)
	logger.Info("The worflow phases are:")
	w.Show(logger)

	// sort workflow's phases (by tier)
	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	// filter the phases in the sorted tiers
	filteredTiers, err := w.filterPhase(logger, sortedTiers, skipPhases, retainPhases)
	if err != nil {
		return fmt.Errorf("failed to filter phases: %w", err)
	}

	// Show the filtered sorted tiers
	logger.Info("Execution plan after filtering:")
	filteredTiers.Show(logger)

	// Create as many slice of function as tiers
	allTierTasks := make([][]syncx.Func, len(filteredTiers))
	// loop over each tier
	for tierId, tier := range filteredTiers {
		// Create a slice of functions
		tasks, err := w.createSliceFunc(ctx, logger, tierId, tier, targets)
		if err != nil {
			return err
		}
		// add the slice to allTierTasks
		allTierTasks[tierId] = tasks
	}

	// Executes the phases: tier by tier
	logger.Info("📌 Starting execution strategy: concurrent within tiers, sequential across tiers")
	// loop over each tier
	for tierId, tier := range filteredTiers {
		tierIdx := tierId + 1
		nbPhase := len(tier)
		// check : if the context has been canceled before starting a new tier
		if ctx.Err() != nil {
			logger.Warnf("Workflow canceled (by user) before starting Tier %d", tierIdx)
			return ctx.Err()
		}

		logger.Infof("👉 Tier %d: %d has concurent phase(s)", tierIdx, nbPhase)
		tasks := allTierTasks[tierId]

		if errs := syncx.RunConcurrently(ctx, tasks); errs != nil {
			switch errs[0] {
			case context.Canceled:
				logger.Warn("Context activation: user canceled workflow execution using ctrl-c")
				return errs[0]
			case context.DeadlineExceeded:
				logger.Warn("Context activation: deadline exceeded defined timeout")
				return errs[0]
			default:
				logger.ErrorWithNoStack(errs[0], "👉 🔴 tier %d/%d : some phases failed with the following errors:", tierIdx, nbPhase)
				return errs[0]
			}
		}
		logger.Infof("👉 🟢 Tier %d : All phases completed successfully", tierIdx)
	}

	logger.Info("📌 🟢 Workflow execution finished successfully")
	return nil
}
