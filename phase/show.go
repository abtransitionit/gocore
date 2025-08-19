// File to create in gocore/phase/display.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com
*/
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
