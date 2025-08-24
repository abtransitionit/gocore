// File in gocore/phase/run.go
package phase

import (
	"context"
	"fmt"
)

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
