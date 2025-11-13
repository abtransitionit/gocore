package phase2

import (
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
	WkfName     string   `yaml:"wkfName,omitempty"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	FnAlias     string   `yaml:"fn"`
	Dependency  []string `yaml:"dependency,omitempty"`
	Param       []string `yaml:"param,omitempty"`
	Node        string   `yaml:"node,omitempty"`
}

// Description: constructor that returns an instance of a Workflow
//
// Parameters:
//
// - fileName: The name of the YAML file containing the workflow definition.
// - cmdPathName: The name of the directory containing the YAML file (relative to the worflow folder).
// - logger: The logger to use for logging.
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

func GetPhase(workflowName string) *Phase {
	return &Phase{
		WkfName: workflowName,
	}
}
