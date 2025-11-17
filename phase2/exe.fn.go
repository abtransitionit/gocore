package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Description: executes code for 1 target
//
// Notes:
// - this function is executed inside a goroutine
func (goFunction *GoFunction) runOnTarget(ctx context.Context, phaseName, targetName string, logger logx.Logger) error {

	// define the return value of the function to be executed on the target
	var oko bool
	var err error

	// log
	logger.Infof("↪ (gofunc) phase: %s > target:%s > running", phaseName, targetName)

	// Execute the function
	oko, err = goFunction.Func(targetName, goFunction.ParamList, logger) // execute the task:PhaseFn (signature is important here)

	// handle system eroor
	if err != nil {
		return err
	}

	// handle logic eroor
	if !oko {
		return fmt.Errorf("↪ (gofunc) phase: %s > target:%s > phase:%s > go:%s > param: %s", phaseName, targetName, goFunction.PhaseName, goFunction.Name, goFunction.ParamList)
	}

	// handle success
	return nil
}
