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

	ShowPhaseList(sortedPhases, logger)

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
// Description: Performs a topological sort on the workflow's phases to determine
//
//	the correct execution order.
//
// Parameters:
//
//   - none
//
// Returns:
//
//   - []Phase: A slice of phases in the correct execution order.
//   - error: An error if a circular dependency is detected.
//
// Notes:
//
//   - Uses Kahn's algorithm for topological sorting.
func (w *Workflow) topologicalSort() ([]Phase, error) {
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

	sortedList := make([]Phase, 0, len(w.Phases))
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		sortedList = append(sortedList, w.Phases[name])

		for _, neighbor := range graph[name] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(sortedList) != len(w.Phases) {
		return nil, fmt.Errorf("circular dependency detected in workflow")
	}

	return sortedList, nil
}
