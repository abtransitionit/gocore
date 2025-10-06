package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

// PrettyPrintTable takes plain-text Helm output and prints it as a formatted table.
func PrettyPrintTable(raw string) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) == 0 {
		fmt.Println("(no data)")
		return
	}

	// Create table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight) // or StyleRounded, StyleColoredBright, etc.

	// Parse headers (first line)
	// headers := strings.Fields(lines[0])
	headers := strings.Split(lines[0], "\t")
	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Parse and append data rows
	for _, line := range lines[1:] {
		// fields := strings.Fields(line)
		fields := strings.Split(line, "\t")
		row := make(table.Row, len(fields))
		for i, f := range fields {
			row[i] = f
		}
		t.AppendRow(row)
	}

	t.Render()
}
