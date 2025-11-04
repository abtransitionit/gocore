// File in gocore/phase/type.go
package phase2

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
)

// Description: represents a set of phases.
type Workflow struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Phases      map[string]Phase `yaml:"phases"`
}

// Description: represents a phase
type Phase struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Fn          string   `yaml:"fn"`
	Dependency  []string `yaml:"dependency,omitempty"`
	Param       []string `yaml:"param,omitempty"`
	Node        string   `yaml:"node,omitempty"`
}

// Description: represents a Go function
type GoFunc struct {
	PhaseFuncName string
	Func          func(ctx context.Context, params any, logger logx.Logger) error
}

// Description: represents a registry to store functions
//
// Notes:
//   - The registry is a map of functions indexed by their PhaseName.
type FunctionRegistry struct {
	funcs map[string]*GoFunc
}

// description: constructor that return an instance of GoFunc
func GetGoFunc(phaseFunction string) *GoFunc {
	return &GoFunc{
		PhaseFuncName: phaseFunction,
	}
}

// description: constructor that return an instance of a FunctionRegistry
func GetFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		funcs: make(map[string]*GoFunc),
	}
}

// Manage Registry
func (r *FunctionRegistry) Add(name string, fn func(context.Context, any, logx.Logger) error) {
	r.funcs[name] = &GoFunc{
		PhaseFuncName: name,
		Func:          fn,
	}
}
func (r *FunctionRegistry) Get(name string) (*GoFunc, bool) {
	f, ok := r.funcs[name]
	return f, ok
}

// List returns the names of all registered functions.
func (r *FunctionRegistry) List() []string {
	names := make([]string, 0, len(r.funcs))
	for k := range r.funcs {
		names = append(names, k)
	}
	return names
}

// description: check if a PhaseFuncName is in the registry
func (fr *FunctionRegistry) Has(key string) bool {
	_, ok := fr.funcs[key]
	return ok
}

// ///////// New code

// type GoFunc2 struct {
// 	PhaseFuncName string
// 	Func          func(ctx context.Context, params any, logger logx.Logger) error
// }

// type FunctionRegistry2 struct {
// 	funcs map[string]*GoFunc2
// }

// func GetFunctionRegistry2() *FunctionRegistry2 {
// 	return &FunctionRegistry2{
// 		funcs: make(map[string]*GoFunc2),
// 	}
// }

// // Add registers a GoFunc2
// func (r *FunctionRegistry2) Add(gf *GoFunc2) error {
// 	if gf == nil || gf.PhaseFuncName == "" {
// 		return fmt.Errorf("invalid GoFunc2 or name")
// 	}
// 	if _, exists := r.funcs[gf.PhaseFuncName]; exists {
// 		return fmt.Errorf("function %q already registered", gf.PhaseFuncName)
// 	}
// 	r.funcs[gf.PhaseFuncName] = gf
// 	return nil
// }

// // Get retrieves a GoFunc2 by name
// func (r *FunctionRegistry2) Get(name string) *GoFunc2 {
// 	if gf, ok := r.funcs[name]; ok {
// 		return gf
// 	}
// 	// fallback if function not found
// 	return &GoFunc2{
// 		PhaseFuncName: name,
// 		Func: func(ctx context.Context, params any, logger logx.Logger) error {
// 			return fmt.Errorf("function not registered: %s", name)
// 		},
// 	}
// }
