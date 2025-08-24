// File in gocore/phase/run.go
package phase

import (
	"context"
	"fmt"
	"strings"

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
func (w *Workflow) Execute(ctx context.Context, logger logx.Logger, skipPhases []int, retainPhases []int) error {
	logger.Info("Starting workflow execution...")

	// check paramaters. Deletgate check of other parameters to the filterPhase function
	if len(w.Phases) == 0 {
		return fmt.Errorf("workflow is empty: no phases defined")
		// logger.ErrorWithNoStack(err, "cannot execute empty workflow")
		// return err
	}

	// if len(w.Phases) == 0 {
	// 	logger.Info("the workflow cannot be emty, please add phases to the workflow")
	// 	return fmt.Errorf("failed to sort phases: %w")
	// }

	// Get the sorted phases
	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	// Get filtered phases
	filteredTiers, err := w.filterPhase(logger, sortedTiers, skipPhases, retainPhases)
	if err != nil {
		return fmt.Errorf("failed to filter phases: %w", err)
	}

	// Show the filtered phases ordered by tier
	filteredTiers.Show(logger)
	// w.ShowPhaseList(filteredTiers, logger)

	logger.Info("--- Starting execution: concurrent within tiers, sequential across tiers ---")

	// loop over each tier
	for tierId, tier := range filteredTiers {
		tierIdx := tierId + 1
		NbPhase := len(tier)
		// Before starting a new tier, check if the context has been canceled.
		if ctx.Err() != nil {
			logger.Info("Workflow canceled by user.")
			return ctx.Err()
		}
		logger.Infof("👉 Executing Tier %d with %d phases concurrently...", tierIdx, NbPhase)

		// Create a slice of functions (for each tier)
		concurrentTasks := make([]syncx.Func, 0, NbPhase)
		for phaseId, phase := range tier {
			phaseIdx := phaseId + 1
			phaseName := phase.Name // <- capture the variable here
			task := adaptToSyncxFunc(phase.fn, ctx, logx.GetLogger(), []string{}...)

			// convert Phases that al have a PhaseFunc sigature to a syncx.Func signature
			wrappedTask := func() error {
				logger.Infof("➡️ Executing phase %d/%d of tier %d : '%s'...", phaseIdx, NbPhase, tierIdx, phaseName)
				if err := task(); err != nil {
					return err
				}
				logger.Infof("➡️🔹 Terminated phase %d/%d of tier %d : '%s'.", phaseIdx, NbPhase, tierIdx, phaseName)
				return nil
			}

			concurrentTasks = append(concurrentTasks, wrappedTask)
		}

		// run all functions in the tier concurently - with all the same ctx.
		if errs := syncx.RunConcurrently(ctx, concurrentTasks); errs != nil {
			// TODO: CHATGPT: COMMENT:Check the first error to determine the reason for cancellation. This is the correct place to log the cancellation event.
			switch errs[0] {
			case context.Canceled:
				logger.Info("Context activation: canceled by user (e.g., via Ctrl+C).")
				return errs[0]
			case context.DeadlineExceeded:
				logger.Info("Context activation: deadline exceeded (e.g., via timeout).")
				return errs[0]
			}

			// Log all collected errors and return the first one to stop the workflow.
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("tier %d / Phase %d (TODO:phase name) failed with the following errors:", tierIdx, NbPhase))
			for _, e := range errs {
				sb.WriteString(fmt.Sprintf("\n- %v", e))
			}
			logger.ErrorWithNoStack(errs[0], "❌ %s", sb.String())
			return errs[0]
		}
		logger.Infof("✅ Tier %d: All %d phases completed successfully.", tierIdx, NbPhase)

	}

	logger.Info("Workflow execution finished.")
	return nil
}
