// File in gocore/list/list.go
package list

import (
	"fmt"
	"sort"
	"strings"
)

// # Purpose
//
// - get the list of a Map:key from a map
//
// # Parameters
//
// - m: map[string]any
//
// # Return
//
// - []string: containing the keys of the map
//
// # Notes
//
// - This function accepts any kind of map
func GetMapKeys[V any](m map[string]V) []string {
	// defining size is more efficient
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key) // Add the current key to our slice.
	}

	// Use the sort package to order the keys alphabetically.
	sort.Strings(keys)

	return keys
}

func GetSlicefromStringWithSep(ListAsString string, sep string) []string {
	return strings.Split(ListAsString, sep)
}

func GetSlicefromStringWithSpace(ListAsString string) []string {
	return strings.Fields(ListAsString)
}

func GetStringWithSepFromSlice(ListString []string, separator string) string {
	return strings.Join(ListString, separator)
}

func GetStringWithSpaceFromSlice(ListString []string) string {
	return GetStringWithSepFromSlice(ListString, " ")
}

// Name: GetFieldByID
//
// Description: returns the field value at column index `fieldIndex` for row number `id`.
//
// Example Usage:
//
//	value, err := GetFieldByID(raw, 3, 2) // in
func GetFieldByID(rawContent string, id int, fieldIndex int) (string, error) {
	lines := strings.Split(strings.TrimSpace(rawContent), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("no data")
	}

	if id <= 0 || id >= len(lines) {
		return "", fmt.Errorf("id %d out of range", id)
	}

	fields := strings.Fields(lines[id]) // works even if columns are space-aligned
	if fieldIndex < 0 || fieldIndex >= len(fields) {
		return "", fmt.Errorf("field index %d out of range", fieldIndex)
	}

	return fields[fieldIndex], nil
}

func GetFieldByID2(raw string, id int, fieldIndex int) (string, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no data")
	}

	if id <= 0 || id > len(lines) {
		return "", fmt.Errorf("id %d out of range", id)
	}

	line := lines[id]                   // IDs start at 1
	fields := strings.Split(line, "\t") // split by tab

	if fieldIndex < 0 || fieldIndex >= len(fields) {
		return "", fmt.Errorf("field index %d out of range", fieldIndex)
	}

	return fields[fieldIndex], nil
}
