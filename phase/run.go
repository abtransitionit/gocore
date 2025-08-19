// File to create in gocore/phase/run.go
package phase

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Name: Execute
//
// Description: a minimal implementation that executes phases sequentially.
//
// Parameters:
//
//   - ctx: The context for the workflow. This allows for cancellation and timeouts.
//   - logger: The logger to use for printing messages.
//
// Returns:
//
//   - An error if the workflow fails to execute.
//
// Notes:
//
//   - This version executes all phases sequentially, tier by tier.
//   - It stops execution on the first error it encounters.
func (w *Workflow) Execute(ctx context.Context, logger logx.Logger) error {
	logger.Info("Starting workflow execution...")

	sortedTiers, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	w.ShowPhaseList(sortedTiers, logger)
	logger.Info("--- Starting sequential execution ---")

	for tierID, tier := range sortedTiers {
		logger.Info("Executing Tier %d with %d phases...", tierID+1, len(tier))
		for _, phase := range tier {
			logger.Info("  -> Executing phase '%s'...", phase.Name)
			err := phase.run(ctx)
			if err != nil {
				return err
			}
			logger.Info("  -> Phase '%s' finished.", phase.Name)
		}
	}

	logger.Info("Workflow execution finished.")
	return nil
}

// Name: run
//
// Description: Executes a single phase's function.
//
// Parameters:
//
//   - ctx: The context for the phase's execution.
//   - ...args: Additional arguments to be passed to the phase's function.
//
// Returns:
//
//   - error: Returns an error if the phase's function fails.
//
// Notes:
//
//   - This method is a helper for the `Execute` method.
func (p *Phase) run(ctx context.Context, args ...string) error {
	if _, err := p.fn(ctx, args...); err != nil {
		return fmt.Errorf("phase %q failed: %w", p.Name, err)
	}
	return nil
}

// Name: topologicalSort
//
// Description:
//
//   - Performs a topological sort on the workflow's phases
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - [][]Phase: A slice of slices, where each inner slice represents a tier of phases that can be run in parallel.
//   - error: An error if a circular dependency is detected.
//
// Notes:
//
//   - Uses Kahn's algorithm for topological sorting.
//   - This function is a helper for the `Execute` method.
//   - The output is used by Execute method to run each phases of a workflow
//   - tthe method determine the correct execution order and group phases (tier) that can be run in parallel.
//   - a circular dependency is detected.
func (w *Workflow) topologicalSort() ([][]Phase, error) {
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	for name := range w.Phases {
		inDegree[name] = 0
		graph[name] = []string{}
	}

	for name, phase := range w.Phases {
		for _, depName := range phase.Dependencies {
			if _, exists := w.Phases[depName]; !exists {
				return nil, fmt.Errorf("dependency %q for phase %q does not exist", depName, name)
			}
			graph[depName] = append(graph[depName], name)
			inDegree[name]++
		}
	}

	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	sortedTiers := make([][]Phase, 0)
	for len(queue) > 0 {
		tierSize := len(queue)
		currentTier := make([]Phase, 0, tierSize)
		nextQueue := make([]string, 0)

		for i := 0; i < tierSize; i++ {
			name := queue[i]
			currentTier = append(currentTier, w.Phases[name])
			for _, neighbor := range graph[name] {
				inDegree[neighbor]--
				if inDegree[neighbor] == 0 {
					nextQueue = append(nextQueue, neighbor)
				}
			}
		}
		sortedTiers = append(sortedTiers, currentTier)
		queue = nextQueue
	}

	if len(sortedTiers) > 0 {
		totalSortedPhases := 0
		for _, tier := range sortedTiers {
			totalSortedPhases += len(tier)
		}
		if totalSortedPhases != len(w.Phases) {
			return nil, fmt.Errorf("circular dependency detected in workflow")
		}
	}

	return sortedTiers, nil
}
