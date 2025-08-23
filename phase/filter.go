// File in gocore/phase/run.go
package phase

import (
	"fmt"
	"os"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
)

// Name: filterPhase
//
// Description: Returns a new set of tiers with specified phases removed or retained.
//
// Parameters:
//   - sortedPhases: The full list of topologically sorted phases.
//   - skipIDs: 		A list of integer IDs to be skipped.
//   - retainIDs: 	A list of integer IDs to be retained.
//
// Returns:
//   - [][]Phase: A new list of tiers
//   - error: Returns an error if a requested ID does not exist.
// Notes:
//   - This does not re-run the topological sort.

func (w *Workflow) filterPhase(logger logx.Logger, sortedPhases PhaseTiers, skipPhases []int, retainPhases []int) (PhaseTiers, error) {
	// l := logx.GetLogger()
	logger.Info(">>> Entering filterPhase")

	// Parameter checks
	if len(skipPhases) > 0 && len(retainPhases) > 0 {
		return nil, fmt.Errorf("invalid parameters: skipPhases and retainPhases cannot be set at the same time")
	}

	if (len(skipPhases) == 0 && len(retainPhases) == 0) || len(w.Phases) == 0 {
		return sortedPhases, nil
	}

	// Get ordered phase names
	ListPhase := list.GetMapKeys(w.Phases)
	logger.Infof("All phases ordered: %v", ListPhase)

	var filterPhaseName []string
	mode := "skip"

	if len(skipPhases) > 0 {
		logger.Infof("Phase(s) to skip by id: %v", skipPhases)
		filterPhaseName = make([]string, len(skipPhases))
		for i, id := range skipPhases {
			if id > len(ListPhase) {
				logger.Errorf("Phase ID %d does not exist in the workflow", id)
				os.Exit(1)
			}
			filterPhaseName[i] = ListPhase[id-1]
		}
		logger.Infof("Phase(s) to skip by id %v", filterPhaseName)
	} else {
		mode = "retain"
		logger.Infof("Phase(s) to retain by id: %v", retainPhases)
		filterPhaseName = make([]string, len(retainPhases))
		for i, id := range retainPhases {
			if id > len(ListPhase) {
				logger.Errorf("Phase(s) ID %d does not exist in the workflow", id)
				os.Exit(1)
			}
			filterPhaseName[i] = ListPhase[id-1]
		}
		logger.Infof("Phases to skip/retain by name: %v", filterPhaseName)
	}

	// Build lookup map
	filterMap := make(map[string]bool)
	for _, name := range filterPhaseName {
		filterMap[name] = true
	}

	// Filter phases
	var filteredPhases PhaseTiers
	for _, tier := range sortedPhases {
		var newTier []Phase
		for _, phase := range tier {
			if (mode == "skip" && !filterMap[phase.Name]) || (mode == "retain" && filterMap[phase.Name]) {
				newTier = append(newTier, phase)
			}
		}
		if len(newTier) > 0 {
			filteredPhases = append(filteredPhases, newTier)
		}
	}

	return filteredPhases, nil
}
