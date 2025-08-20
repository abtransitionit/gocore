// File to create in gocore/list/list.go
package list

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

	return keys
}
