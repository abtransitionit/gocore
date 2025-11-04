package phase2

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/yamlx"
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
	FnAlias     string   `yaml:"fn"`
	Dependency  []string `yaml:"dependency,omitempty"`
	Param       []string `yaml:"param,omitempty"`
	Node        string   `yaml:"node,omitempty"`
}

// Description: represents a Go function
type PhaseFn struct {
	Name string
	Func func(ctx context.Context, params any, logger logx.Logger) error
}

// Description: represents a registry to store functions
//
// Notes:
//   - The registry is a map of functions indexed by their PhaseName.
type FnRegistry struct {
	funcs map[string]*PhaseFn
}

// ---------- CONSTRUCTOR ----------

// Description: constructor that returns an instance of a Workflow
func GetWorkflow(fileName, cmdPathName string, logger logx.Logger) (*Workflow, error) {

	// 1. Define the path of the workflow YAML
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("getting caller information")
	}
	workflowFilePath := filepath.Join(path.Dir(file), "..", cmdPathName, fileName)

	// log
	logger.Debugf("found workflow file: %s", workflowFilePath)

	// 2. Load the yaml file into a struct
	workflow, err := yamlx.LoadFile[Workflow](workflowFilePath)
	if err != nil {
		return nil, fmt.Errorf("loading workflow yaml file from %s: %w", workflowFilePath, err)
	}

	return workflow, nil

}

// description: constructor that return an instance of PhaseFn
func GetPhaseFn(phaseFunction string) *PhaseFn {
	return &PhaseFn{
		Name: phaseFunction,
	}
}

// description: constructor that return an instance of a FnRegistry
func GetFnRegistry() *FnRegistry {
	return &FnRegistry{
		funcs: make(map[string]*PhaseFn),
	}
}
