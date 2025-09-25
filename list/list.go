// File in gocore/list/list.go
package list

import (
	"fmt"
	"sort"
	"strings"

	"github.com/abtransitionit/gocore/color"
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
