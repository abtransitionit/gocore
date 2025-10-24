package list

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/abtransitionit/gocore/color"
	"github.com/jedib0t/go-pretty/table"
)

func CountNbLine(raw string) int {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

// PrettyPrintTable prints a table from raw input, using the first line as header.
// PrettyPrintTable prints a table from raw input, using the first line as header.
// Compatible with existing calls.
func PrettyPrintTable(raw string) {
	printTablePretty(raw, 0) // skipLines = 0 → first line is header
}

// PrettyPrintTableWithSkipLine prints a table skipping the first skipLines lines.
// Use this when you want to skip any number of lines at the top.
func PrettyPrintTableWithSkipLine(raw string, skipLines int) {
	printTablePretty(raw, skipLines)
}

// printTablePretty handles printing, skipping the first skipLines lines
func printTablePretty(raw string, skipLines int) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) == 0 {
		fmt.Println("(no data)")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)

	startLine := skipLines

	// If skipLines == 0, treat first line as header
	if skipLines == 0 {
		headers := strings.Split(lines[0], "\t")
		headerRow := make(table.Row, len(headers)+1)
		headerRow[0] = "ID"
		for i, h := range headers {
			headerRow[i+1] = h
		}
		t.AppendHeader(headerRow)
		startLine = 1
	}

	// Append data rows
	for i, line := range lines[startLine:] {
		fields := strings.Split(line, "\t")
		row := make(table.Row, len(fields)+1)
		row[0] = i + 1
		for j, f := range fields {
			row[j+1] = f
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

	// Add header with ID column
	t.AppendHeader(table.Row{"ID", "Key", "Value"})

	id := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		m := re.FindStringSubmatch(line)
		if len(m) == 3 {
			key := strings.TrimSpace(m[1])
			val := strings.Trim(strings.TrimSpace(m[2]), `"`)
			// t.AppendRow(table.Row{key, val})
			t.AppendRow(table.Row{id, key, val})
		} else {
			// fallback: single column if format unexpected
			// t.AppendRow(table.Row{line})
			t.AppendRow(table.Row{id, line, ""})

		}
		id++
	}

	// // Optional: set column headers (can be commented out if unwanted)
	// t.AppendHeader(table.Row{"Key", "Value"})

	t.Render()
}
