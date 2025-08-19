// File to create in gocore/phase/run.go
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
// Description: a minimal implementation that executes phases concurrently, tier by tier.
//
// Parameters:
//
//   - ctx: The context for the workflow. This allows for cancellation and timeouts.
//   - logger: The logger to use for printing messages.
//   - skipPhases: A slice of integer IDs representing the phases to be skipped.
//
// Returns:
//
//   - An error if the workflow fails to execute.
//
// Notes:
//
//   - This version executes all phases in a tier concurrently.
//   - It stops execution on the first error it encounters.
func (w *Workflow) Execute(ctx context.Context, logger logx.Logger, skipPhases []int) error {
	logger.Info("Starting workflow execution...")

	// Logging the received IDs, as requested.
	logger.Info("Received phase IDs to skip: %v", skipPhases)

	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	// Filter out the phases to be skipped.
	filteredTiers, err := w.filterPhases(sortedTiers, skipPhases)
	if err != nil {
		return fmt.Errorf("failed to filter phases: %w", err)
	}

	// Show the final phase list
	w.ShowPhaseList(filteredTiers, logger)

	logger.Info("--- Starting concurrent execution ---")

	for tierID, tier := range filteredTiers {
		// Before starting a new tier, check if the context has been canceled.
		if ctx.Err() != nil {
			logger.Info("Workflow canceled by user.")
			return ctx.Err()
		}

		logger.Info("Executing Tier %d with %d phases concurrently...", tierID+1, len(tier))

		// Create a slice of functions that fit the syncx.Func signature.
		concurrentTasks := make([]syncx.Func, 0, len(tier))
		for _, phase := range tier {
			// Use the new adapter function to convert our PhaseFunc into a syncx.Func.
			// Corrected: Add the '...' to unpack the empty slice.
			task := adaptToSyncxFunc(phase.fn, ctx, []string{}...)

			// Wrap the task to add logging for this specific phase.
			wrappedTask := func(phaseName string) syncx.Func {
				return func() error {
					logger.Info("  -> Executing phase '%s'...", phaseName)
					if err := task(); err != nil {
						return err
					}
					logger.Info("  -> Phase '%s' finished.", phaseName)
					return nil
				}
			}(phase.Name)

			concurrentTasks = append(concurrentTasks, wrappedTask)
		}

		// Use syncx package to run all tasks in the tier concurrently.
		// Corrected: Pass the context as the first argument.
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
