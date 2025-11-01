// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/abtransitionit/gocore/yamlx"
)

func GetWorkflow() (*Workflow, error) {
	// 1. Define YAML workflow file path
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("could not get caller information")
	}
	workflowPath := filepath.Join(path.Dir(file), "wkf.yaml")

	fmt.Println("workflowPath:", workflowPath)

	// 2. Load the yaml using the generic function from lib.go
	workflow, err := yamlx.LoadFile[Workflow2](workflowPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load workflow from %s: %w", workflowPath, err)
	}

	return workflow, nil

}
