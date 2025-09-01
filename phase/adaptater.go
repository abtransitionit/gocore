// File in gocore/phase/adapter.go
package phase

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/syncx"
)

// Name: adaptToSyncxFunc
//
// Description: adapts a PhaseFunc to a syncx.Func, passing the list of targets.
//
// Parameters:
//   - fn: PhaseFunc to adapt
//   - ctx: context for cancellation/timeouts
//   - l: logger
//   - targets: list of targets (VMs, etc.) for this phase
//   - cmd: optional arguments
//
// Returns:
//   - syncx.Func suitable for concurrent execution
// func adaptToSyncxFunc(fn PhaseFunc, ctx context.Context, l logx.Logger, targets []Target, cmd ...string) syncx.Func {
// 	return func() error {
// 		_, err := fn(ctx, l, targets, cmd...)
// 		return err
// 	}
// }

// Name: adaptToSyncxFunc
//
// Description: takes a PhaseFunc and returns a syncx.Func.
//
// Parameters:
//
//   - fn: PhaseFunc to adapt
//   - ctx: context for cancellation/timeouts
//   - l: logger
//   - targets: list of targets (VMs, etc.) for this phase
//   - cmd: optional arguments//
//
// Returns:
//
//   - syncx.Func: A function that represents the adapted PhaseFunc.
//
// Notes:
//   - This acts as an adapter, making a PhaseFunc compatible with the syncx.RunConcurrently function's signature.
//   - It wraps the PhaseFunc in a syncx.Func: we create a closure that is an anonymous function.
//
// Todo
//   - pass logging to the closure that is an anonymous function.
//   - make that function not an anonymous.
func adaptToSyncxFunc(fn PhaseFunc, ctx context.Context, l logx.Logger, targets []Target, cmd ...string) syncx.Func {
	return func() error {
		_, err := fn(ctx, l, targets, cmd...)
		return err
	}
}

// Name: createSliceFunc
//
// Description:
//
//	Builds a slice of wrapped syncx.Func tasks for all phases in a given tier.
//	Each phase is adapted, wrapped with logging (before/after execution), and added
//	to the slice. The resulting slice can then be executed concurrently.
//
// Parameters:
//   - ctx:    The context used to control cancellation and timeouts.
//   - logger: Logger instance used for structured logging.
//   - tierId: The index of the current tier (0-based).
//   - tier:   The list of phases belonging to this tier.
//
// Returns:
//   - []syncx.Func: A slice of wrapped functions ready to be executed concurrently.
//   - error:        An error if building the slice fails.
//
// Notes:
//   - Each wrapped task logs the start and end of its execution.
//   - If a task fails, it logs the error and returns it.
//   - Nothing is executed at this stage — the slice only prepares tasks for later execution.
//
// Todo:
//   - Externalize the wrapping logic into a named function instead of using an inline closure.
//   - Allow custom log message formatting for phases (e.g., configurable success/failure symbols).
func (w *Workflow) createSliceFunc(ctx context.Context, logger logx.Logger, tierId int, tier []Phase, targets []Target) ([]syncx.Func, error) {
	nbPhase := len(tier)
	concurrentTasks := make([]syncx.Func, 0, nbPhase)

	for phaseId, phase := range tier {
		phaseIdx := phaseId + 1
		phaseName := phase.Name

		// Adapt PhaseFunc -> syncx.Func
		task := adaptToSyncxFunc(phase.fn, ctx, logger, targets, []string{}...)

		// Wrap task with logging
		wrappedTask := func() error {
			logger.Debugf("➡️ Tier %d: running phase %d/%d named %s", tierId+1, phaseIdx, nbPhase, phaseName)
			if err := task(); err != nil {
				return fmt.Errorf("➡️ 🔴 phase %d/%d of tier %d (%s) failed: %w", phaseIdx, nbPhase, tierId+1, phaseName, err)
			}
			logger.Debugf("➡️ 🟢 phase %d/%d of tier %d: %s completed successfully", phaseIdx, nbPhase, tierId+1, phaseName)
			return nil
		}

		concurrentTasks = append(concurrentTasks, wrappedTask)
	}

	return concurrentTasks, nil
}
