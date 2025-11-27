// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/abtransitionit/gocore/mock/yamlx"
)

func GetWorkflow(cmdPathName string) (*Workflow, error) {
	// 1. Define YAML workflow file path
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("could not get caller information")
	}
	workflowPath := filepath.Join(path.Dir(file), "..", cmdPathName, "wkf.phase.yaml")

	fmt.Println("workflowPath:", workflowPath)

	// 2. Load the yaml using the generic function from lib.go
	workflow, err := yamlx.LoadFile[Workflow](workflowPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load workflow from %s: %w", workflowPath, err)
	}

	return workflow, nil

}
