/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

*/

package filex

import "strings"

// deletes the left spaces on a multiline string and returns it
func DeleteLeftSpace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, " \t")
	}
	return strings.Join(lines, "\n")
}

func DeleteLeftTab(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, "\t") // trims only leading tabs
		if strings.TrimSpace(trimmed) != "" {   // skip empty lines
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}
