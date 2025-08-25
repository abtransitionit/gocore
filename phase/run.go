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
func (w *Workflow) Execute(ctx context.Context, logger logx.Logger, skipPhases []int, retainPhases []int) error {
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

	// sort workflow's phases
	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	// filter workflow's phases
	filteredTiers, err := w.filterPhase(logger, sortedTiers, skipPhases, retainPhases)
	if err != nil {
		return fmt.Errorf("failed to filter phases: %w", err)
	}

	// Optional: Show the filtered phases ordered by tiers
	logger.Debug("Execution plan after filtering:")
	filteredTiers.Show(logger)

	logger.Info("📌 Starting execution strategy: concurrent within tiers, sequential across tiers")

	// loop over each tier - run tiers sequentially
	for tierId, tier := range filteredTiers {
		tierIdx := tierId + 1
		nbPhase := len(tier)
		// check : if the context has been canceled before starting a new tier
		if ctx.Err() != nil {
			logger.Warnf("Workflow canceled (by user) before starting Tier %d", tierIdx)
			return ctx.Err()
		}

		logger.Infof("👉 Executing Tier %d with %d phase(s) concurrently...", tierIdx, nbPhase)

		// Create a slice of functions (for each tier)
		// Build tasks for concurrent execution
		concurrentTasks := make([]syncx.Func, 0, nbPhase)
		for phaseId, phase := range tier {
			phaseIdx := phaseId + 1
			phaseName := phase.Name
			// create a syncx.Func from PhaseFunc
			// convert Phases that al have a PhaseFunc sigature to a syncx.Func signature
			task := adaptToSyncxFunc(phase.fn, ctx, logger, []string{}...)

			wrappedTask := func() error {
				logger.Debugf("➡️ running phase %d/%d of tier %d: %s", phaseIdx, nbPhase, tierIdx, phaseName)
				if err := task(); err != nil {
					return fmt.Errorf("➡️ 🔴 phase %d/%d of tier %d (%s) failed: %w", phaseIdx, nbPhase, tierIdx, phaseName, err)
				}
				logger.Debugf("➡️ 🟢 phase %d/%d of tier %d: %s completed successfully", phaseIdx, nbPhase, tierIdx, phaseName)
				return nil
			}
			concurrentTasks = append(concurrentTasks, wrappedTask)
		}

		// Run all phases in this tier concurrently using syncx with the same context
		if errs := syncx.RunConcurrently(ctx, concurrentTasks); errs != nil {
			switch errs[0] {
			case context.Canceled:
				logger.Warn("Context activation: user canceled workflow execution using ctrl-c")
				return errs[0]
			case context.DeadlineExceeded:
				logger.Warn("Context activation: deadline exceeded defined timeout")
				return errs[0]
			default:
				logger.ErrorWithNoStack(errs[0], "👉 🔴 tier %d/%d : some phases failed with the following errors:", tierIdx, nbPhase)
				// for _, e := range errs {
				// 	logger.ErrorWithNoStack(e, "👉 🔴 tier %d/%d : some phases failed with the following errors:", tierIdx, nbPhase)
				// }
				return errs[0]
			}
			// // Log all collected errors and return the first one to stop the workflow.
			// var sb strings.Builder
			// sb.WriteString(fmt.Sprintf("👉 🔴 tier %d/%d : some phases failed with the following errors:", tierIdx, nbPhase))
			// for _, e := range errs {
			// 	sb.WriteString(fmt.Sprintf("\n- %v", e))
			// }
			// // logger.ErrorWithNoStack(errs[0], "%s", sb.String())
			// return errs[0]
		}
		logger.Infof("👉 🟢 Tier %d : All %phases completed successfully", tierIdx)
	}

	logger.Info("🎉 Workflow execution finished successfully")
	return nil
}
