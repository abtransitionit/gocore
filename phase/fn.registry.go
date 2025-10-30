package phase

import "fmt"

// Description: a private registry
//
// Notes:
// - map a YAML function denote as a string to a real Go functions.
// - private and not exported.
var registry = make(map[string]PhaseFunc2)

// description: add a function to the registry
func RegisterSingleFunc(name string, f PhaseFunc2) {
	if _, exists := registry[name]; exists {
		fmt.Printf("⚠️  overwriting registration for %q\n", name)
	}
	registry[name] = f
}
