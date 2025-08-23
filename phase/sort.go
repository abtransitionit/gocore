// File in gocore/phase/run.go
package phase

import (
	"fmt"
	"sort"
)

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
//
// Example of output:
//
//   - sortedTiers = [][]Phase{tier1, tier2, tier3, tier4}
//   - tier1 = [phase1, phase2]
//   - tier2 = [phase3, phase4]
//   - tier3 = [phase5, phase6]
//   - tier4 = [phase7, phase8]
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
	// Ensure deterministic ordering of phases in the same tier
	sort.Strings(queue)

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

		// Sort again to keep deterministic order across tiers
		sort.Strings(nextQueue)
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
