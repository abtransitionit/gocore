package list

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/abtransitionit/gocore/color"
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

// PrettyPrint prints []string with rotating colors
func PrettyPrint(list []string) {
	colors := []string{
		color.Red,
		color.Green,
		color.Yellow,
		color.Blue,
		color.Magenta,
		color.Cyan,
	}
	for i, item := range list {
		fmt.Println(color.Colorize(fmt.Sprintf("- %s", item), colors[i%len(colors)]))
	}
}

func PrettyPrintKvpair(raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		fmt.Println("(no data)")
		return
	}

	lines := strings.Split(raw, "\n")

	// Table writer setup
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)

	// Prepare rows
	re := regexp.MustCompile(`^([^=]+)=(.*)$`)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		m := re.FindStringSubmatch(line)
		if len(m) == 3 {
			key := strings.TrimSpace(m[1])
			val := strings.Trim(strings.TrimSpace(m[2]), `"`)
			t.AppendRow(table.Row{key, val})
		} else {
			// fallback: single column if format unexpected
			t.AppendRow(table.Row{line})
		}
	}

	// Optional: set column headers (can be commented out if unwanted)
	t.AppendHeader(table.Row{"Key", "Value"})

	t.Render()
}
