// File in gocore/phase/run.go
package phase

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
)

// Name: filterPhases
//
// Description: Returns a new set of tiers with specified phases removed.
//
// Parameters:
//   - sortedPhases: The full list of topologically sorted phases.
//   - skipIDs: A map of integer IDs to be skipped.
//
// Returns:
//   - [][]Phase: A new list of tiers with the skipped phases removed.
//   - error: Returns an error if a requested ID does not exist.
//
// Notes:
//   - This does not re-run the topological sort.
// func (w *Workflow) filterPhases(sortedPhases PhaseTiers, skipPhases []int) (PhaseTiers, error) {
// 	l := logx.GetLogger()
// 	l.Info(">>> Entering filterPhases")

// 	var skippedPhaseName []string // list of phases name to skip
// 	var filteredPhases PhaseTiers // list of filtered phases

// 	// manage parameters
// 	if len(skipPhases) == 0 || len(w.Phases) == 0 {
// 		return sortedPhases, nil
// 	}
// 	// log
// 	l.Infof("Received phase IDs to skip: %v", skipPhases)

// 	// Get Map:Key in a slice
// 	var ListPhase = list.GetMapKeys(w.Phases)
// 	l.Infof("List phase ordered: %v", ListPhase)

// 	// create the list of phase name to skip
// 	for _, phaseID := range skipPhases {
// 		// check if the phase ID exist
// 		if phaseID > len(ListPhase) {
// 			l.Errorf("Phase ID %d does not exist in the workflow", phaseID)
// 			os.Exit(1)
// 		}

// 		skippedPhaseName = append(skippedPhaseName, ListPhase[phaseID-1])
// 	}
// 	// log
// 	l.Infof("List phase name to skip: %v", skippedPhaseName)

// 	// Create a map from slice for efficient lookups.
// 	skippedPhasesMapTemp := make(map[string]bool)
// 	for _, name := range skippedPhaseName {
// 		skippedPhasesMapTemp[name] = true
// 	}

//		// create the filtered phases
//		for _, tier := range sortedPhases {
//			var newTier []Phase
//			for _, phase := range tier {
//				// Check if the phase is in the skipped map. If not, add it.
//				if !skippedPhasesMapTemp[phase.Name] {
//					newTier = append(newTier, phase)
//				}
//			}
//			// Only append the new tier if it's not empty.
//			if len(newTier) > 0 {
//				filteredPhases = append(filteredPhases, newTier)
//			}
//		}
//		// return
//		return filteredPhases, nil
//	}
func (w *Workflow) filterSkipPhase(sortedPhases PhaseTiers, skipPhases []int) (PhaseTiers, error) {
	l := logx.GetLogger()
	l.Info(">>> Entering filterSkipPhase")

	var skippedPhaseName []string // list of phases name to skip
	var filteredPhases PhaseTiers // list of filtered phases

	// manage parameters
	if len(skipPhases) == 0 || len(w.Phases) == 0 {
		return sortedPhases, nil
	}
	// log
	l.Infof("Received phase IDs to skip: %v", skipPhases)

	// Get Map:Key in a slice
	var ListPhase = list.GetMapKeys(w.Phases)
	l.Infof("List phase ordered: %v", ListPhase)

	// create the list of phase name to skip
	for _, phaseID := range skipPhases {
		// check if the phase ID exist
		if phaseID > len(ListPhase) {
			l.Errorf("Phase ID %d does not exist in the workflow", phaseID)
			os.Exit(1)
		}

		skippedPhaseName = append(skippedPhaseName, ListPhase[phaseID-1])
	}
	// log
	l.Infof("List phase name to skip: %v", skippedPhaseName)

	// Create a map from slice for efficient lookups.
	skippedPhasesMapTemp := make(map[string]bool)
	for _, name := range skippedPhaseName {
		skippedPhasesMapTemp[name] = true
	}

	// create the filtered phases
	for _, tier := range sortedPhases {
		var newTier []Phase
		for _, phase := range tier {
			// Check if the phase is in the skipped map. If not, add it.
			if !skippedPhasesMapTemp[phase.Name] {
				newTier = append(newTier, phase)
			}
		}
		// Only append the new tier if it's not empty.
		if len(newTier) > 0 {
			filteredPhases = append(filteredPhases, newTier)
		}
	}
	// return
	return filteredPhases, nil
}

// exported wrapper
func (sortedPhases PhaseTiers) Filter(wkfl Workflow, l logx.Logger, skipPhases []int) PhaseTiers {
	// Get filtered phases
	filteredTiers, _ := wkfl.filterSkipPhase(sortedPhases, skipPhases)

	return filteredTiers
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

// Name: SortedPhases
//
// Description: Sort phases of a worflow
//
//   - Returns a slice of slices, where each inner slice represents a tier of phases that can be run in parallel.
//
// Parameters:
//
//   - ctx: The context for the workflow. This allows for cancellation and timeouts.
//
// Returns:
//
//   - [][]Phase: A slice of slices that denotes each a set of phases
//   - error: An error if a circular dependency is detected.
func (w *Workflow) TopoSort(ctx context.Context) (PhaseTiers, error) {
	sortedByTier, err := w.topologicalSort()
	if err != nil {
		return nil, fmt.Errorf("failed to sort phases: %w", err)
	}

	// filteredTiers, err := w.filterPhases(sortedTiers, []int{}) // no skipped phases
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to filter phases: %w", err)
	// }

	return sortedByTier, nil
}
