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

	// 1 - define the return value of the function to be executed on the target
	var ok bool
	var err error

	// log
	// logger.Infof("↪ (gofunc) phase: %s > target:%s > running", phaseName, targetName)

	// 2 - execute the function
	// 21 - tell this function to MANAGE this: goFunction.ParamList can be nil, not defined, empty st can be nil, not defined, empty
	logger.Infof("↪ (gofunc) phase: %s > target:%s > running > ParamList: %v", phaseName, targetName, goFunction.ParamList)
	ok, err = goFunction.Func(targetName, goFunction.ParamList, logger) // execute the task:PhaseFn (signature is important here)

	// handle system eroor
	if err != nil {
		return err
	}

	// handle logic eroor
	if !ok {
		return fmt.Errorf("↪ (gofunc) phase: %s > target:%s > phase:%s > go:%s > param: %s", phaseName, targetName, goFunction.PhaseName, goFunction.Name, goFunction.ParamList)
	}

	// handle success
	return nil
}
