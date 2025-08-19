// File to create in gocore/phase/run.go
package phase

import (
	"context"
	"fmt"
	"log"

	"github.com/abtransitionit/gocore/logx"
)

// Name: Execute
//
// Description: a minimal implementation that displays the
//
//	topologically sorted list of phases to be executed.
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
//   - For now, it simply logs a message and shows the execution plan.
//   - It does not yet execute the phases.
//   - will execute each phase in the workflow sequentially.

func (w *Workflow) Execute(ctx context.Context, logger logx.Logger) error {
	log.Println("Starting workflow execution...")

	sortedPhases, err := w.topologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort phases: %w", err)
	}

	w.ShowPhaseList(sortedPhases, logger)

	log.Println("Workflow execution finished.")
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
// Description: Performs a topological sort on the workflow's phases
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
//   - a circular dependency is detected.
//
// .  - The output is used by Execute to decide the phases that can be run in parallel in goroutine according the dependency.
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
