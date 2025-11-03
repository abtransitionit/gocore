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

// Description: represents a registry to store functions
//
// Notes:
//   - The registry is a map of functions indexed by their PhaseName.
type FunctionRegistry struct {
	funcs map[string]any
}

// Description: represents a Go function
type GoFunc struct {
	PhaseFuncName string
	Func          func(ctx context.Context, logger logx.Logger) error
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
		funcs: make(map[string]any),
	}
}
