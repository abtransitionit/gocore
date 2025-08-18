// File to create in gocore/phase/display.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com
*/
package phase

import (
	"fmt"
	"os"

	"github.com/abtransitionit/gocore/logx"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Name: Show
//
// Description:
// Displays a list of phases in a pretty, formatted table.
// It is intended for use with a CLI to show users a clear
// overview of the available phases in a sequence.
//
// Parameters:
// - l: The logger instance to use for output.
func (pl PhaseList) Show(l logx.Logger) {
	// Create a new table writer with standard output.
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Set the table headers.
	t.AppendHeader(table.Row{"#", "Name", "Description"})

	// Iterate over the PhaseList and add each phase's details to the table.
	for i, p := range pl {
		t.AppendRow(table.Row{i + 1, p.Name, p.Description})
	}

	// Render the table to standard output.
	l.Info("Available Phases:")
	t.Render()
	fmt.Println() // Add a newline for better readability after the table.
}
