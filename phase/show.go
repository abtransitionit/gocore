// File in gocore/phase/show.go
package phase

import (
	"fmt"
	"os"
	"sort"

	"github.com/abtransitionit/gocore/logx"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Name: Show
//
// Description: Displays a list of phases in a pretty, human-readable table.
func (w *Workflow) Show(l logx.Logger) {
	// Create a new table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Phase ID", "Phase", "Description", "Dependencies"}) // Changed Header

	// Get phase names and sort them
	names := make([]string, 0, len(w.Phases))
	for name := range w.Phases {
		names = append(names, name)
	}
	sort.Strings(names)

	// Append rows
	for i, name := range names { // Added counter 'i'
		phase := w.Phases[name]
		deps := "none"
		if len(phase.Dependencies) > 0 {
			deps = fmt.Sprintf("%v", phase.Dependencies)
		}
		t.AppendRow(table.Row{i + 1, phase.Name, phase.Description, deps}) // Added ID
	}
	l.Info("Available phases:")
	// Render the table
	t.Render()
	fmt.Println() // Add a newline for better readability after the table.
}

// Name: ShowPhaseList
//
// Description: Displays a comprehensive human-readable table of phases to be executed for a workflow,
//
//	with an ID column indicating parallel execution tiers.
//
// Parameters:
//   - sortedPhases: The slice of slices of phases to be displayed.
//   - l: The logger to use for printing information.
//
// Notes:
//   - It display a comprehensive list of all phases that will run in parallel.
// func (w *Workflow) ShowPhaseList(sortedPhases [][]Phase, l logx.Logger) {
// 	// Create a new table writer
// 	t := table.NewWriter()
// 	t.SetOutputMirror(os.Stdout)
// 	t.AppendHeader(table.Row{"Tier ID", "Phase", "Description", "Dependencies"})

// 	// Append rows from the provided phase list
// 	for id, tier := range sortedPhases {
// 		for _, phase := range tier {
// 			deps := "none"
// 			if len(phase.Dependencies) > 0 {
// 				deps = fmt.Sprintf("%v", phase.Dependencies)
// 			}
// 			t.AppendRow(table.Row{id + 1, phase.Name, phase.Description, deps})
// 		}
// 	}
// 	l.Info("Phases to be executed in order:")
// 	// Render the table
// 	t.Render()
// 	fmt.Println() // Add a newline for better readability after the table.
// }

// Show is a method on the slice of phase tiers that formats and prints the phases to the logger.
// Show is a method on the PhaseTiers type that formats and prints the phases to the logger.
func (sortedPhases PhaseTiers) Show(l logx.Logger) {
	// Create a new table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Tier ID", "Phase", "Description", "Dependencies"})

	// Append rows from the provided phase list
	for id, tier := range sortedPhases {
		for _, phase := range tier {
			deps := "none"
			if len(phase.Dependencies) > 0 {
				deps = fmt.Sprintf("%v", phase.Dependencies)
			}
			t.AppendRow(table.Row{id + 1, phase.Name, phase.Description, deps})
		}
	}
	l.Info("Phases sorted by tier:")
	// Render the table
	t.Render()
	fmt.Println() // Add a newline for better readability after the table.
}
