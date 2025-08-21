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
func (w *Workflow) Execute(ctx context.Context, logger logx.Logger, skipPhases []int) error {
	logger.Info("Starting workflow execution...")

	// Logging the received IDs, as requested.
	logger.Infof("Received phase IDs to skip: %v", skipPhases)

	// Get the sorted phases
	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	// Get filtered phases
	filteredTiers, err := w.filterPhases(sortedTiers, skipPhases)
	if err != nil {
		return fmt.Errorf("failed to filter phases: %w", err)
	}

	// Show the filtered phases ordered by tier
	filteredTiers.Show(logger)
	// w.ShowPhaseList(filteredTiers, logger)

	logger.Info("--- Starting concurrent execution ---")

	// loop over each tier
	for tierID, tier := range filteredTiers {
		// Before starting a new tier, check if the context has been canceled.
		if ctx.Err() != nil {
			logger.Info("Workflow canceled by user.")
			return ctx.Err()
		}

		logger.Infof("Executing Tier %d with %d phases concurrently...", tierID+1, len(tier))

		// Create a slice of functions (for each tier)
		concurrentTasks := make([]syncx.Func, 0, len(tier))
		for _, phase := range tier {
			// create the closure (needed by syncx) from the phase's function - pass the context
			task := adaptToSyncxFunc(phase.fn, ctx, logx.GetLogger(), []string{}...)

			// Wrap the task to add logging for this specific phase.
			wrappedTask := func(phaseName string) syncx.Func {
				return func() error {
					logger.Infof("  -> Executing phase '%s'...", phaseName)
					if err := task(); err != nil {
						return err
					}
					logger.Infof("  -> Phase '%s' finished.", phaseName)
					return nil
				}
			}(phase.Name)

			concurrentTasks = append(concurrentTasks, wrappedTask)
		}

		// run all phases in the tier concurently - with all the same ctx.
		if errs := syncx.RunConcurrently(ctx, concurrentTasks); errs != nil {
			// Check the first error to determine the reason for cancellation.
			// This is the correct place to log the cancellation event.
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
			sb.WriteString(fmt.Sprintf("tier %d failed with the following errors:", tierID+1))
			for _, e := range errs {
				sb.WriteString(fmt.Sprintf("\n- %v", e))
			}
			logger.ErrorWithNoStack(errs[0], "%s", sb.String())
			return errs[0]
		}
	}

	logger.Info("Workflow execution finished.")
	return nil
}
