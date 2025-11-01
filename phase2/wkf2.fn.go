// File in gocore/phase/adapter.go
package phase2

import "fmt"

// Description: add a string (that denote a function) to a map
//
// Notes:
// - this key will be then be used to manipulate a function
//
// Example usage:
// - var registry = map[string]PhaseFunc2
//
// - var registry = map[string]HttpFunc
func RegisterSingleFunc[T any](fnRegistry map[string]T, name string, f T) {
	if _, exists := fnRegistry[name]; exists {
		fmt.Printf("⚠️  overwriting registration for %q\n", name)
	}
	fnRegistry[name] = f
}

// func RegisterSinglePhaseFn(name string, f PhaseFunc2) {
// 	var registry = make(map[string]PhaseFunc2)

// 	if _, exists := registry[name]; exists {
// 		fmt.Printf("⚠️  overwriting registration for %q\n", name)
// 	}
// 	registry[name] = f
// }

func ResolveFn(fnRegistry map[string]any, name string) (any, error) {
	fn, ok := fnRegistry[name]
	if !ok {
		return nil, fmt.Errorf("no function registered for %q", name)
	}
	return fn, nil
}

// var ResolveFn = map[string]func([]string){
// 	"checkVmAccess": func(nodes []string) { fmt.Println("mock: checkVmAccess") },
// 	"copyAgent":     func(nodes []string) { fmt.Println("mock: copyAgent") },
// 	"osUpgrade":     func(nodes []string) { fmt.Println("mock: osUpgrade") },
// }

// var ResolveFn = map[string]func([]string){
// 	"checkVmAccess": func(nodes []string) { fmt.Println("mock: checkVmAccess") },
// 	"copyAgent":     func(nodes []string) { fmt.Println("mock: copyAgent") },
// 	"osUpgrade":     func(nodes []string) { fmt.Println("mock: osUpgrade") },
// }
