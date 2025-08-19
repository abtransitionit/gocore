// File to create in gocore/phase/show.go
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
	t.AppendHeader(table.Row{"Phase", "Description", "Dependencies"})

	// Get phase names and sort them
	names := make([]string, 0, len(w.Phases))
	for name := range w.Phases {
		names = append(names, name)
	}
	sort.Strings(names)

	// Append rows
	for _, name := range names {
		phase := w.Phases[name]
		deps := "none"
		if len(phase.Dependencies) > 0 {
			deps = fmt.Sprintf("%v", phase.Dependencies)
		}
		t.AppendRow(table.Row{phase.Name, phase.Description, deps})
	}
	l.Info("Available phases:")
	// Render the table
	t.Render()
	fmt.Println() // Add a newline for better readability after the table.
}

// Name: ShowPhaseList
//
// Description: Displays a given list of phases in a pretty, human-readable table.
//
// Parameters:
//   - phases: The slice of phases to be displayed.
//   - l: The logger to use for printing information.
//
// Notes:
//   - This function is reusable and can display any slice of Phase objects.
func ShowPhaseList(phases []Phase, l logx.Logger) {
	// Create a new table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Phase", "Description", "Dependencies"})

	// Append rows from the provided phase list
	for _, phase := range phases {
		deps := "none"
		if len(phase.Dependencies) > 0 {
			deps = fmt.Sprintf("%v", phase.Dependencies)
		}
		t.AppendRow(table.Row{phase.Name, phase.Description, deps})
	}
	l.Info("Phases to be executed in order:")
	// Render the table
	t.Render()
	fmt.Println() // Add a newline for better readability after the table.
}
