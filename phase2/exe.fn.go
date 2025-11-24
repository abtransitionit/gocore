package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Description: executes code for 1 host
//
// Notes:
// - this function is executed inside a goroutine
func (goFunction *GoFunction) runOnHOst(ctx context.Context, phaseName, hostName string, logger logx.Logger) error {

	// 1 - define the return value of the function to be executed on the host
	var ok bool
	var err error

	// log
	// logger.Infof("↪ (gofunc) phase: %s > host:%s > running", phaseName, hostName)

	// 2 - execute the function
	// 21 - tell this function to MANAGE this: goFunction.ParamList can be nil, not defined, empty st can be nil, not defined, empty
	// logger.Infof("↪ (goroutine) %s/%s > running > ParamList: %v", phaseName, hostName, goFunction.ParamList)
	ok, err = goFunction.Func(goFunction.PhaseName, hostName, goFunction.ParamList, logger) // execute the task:PhaseFn (signature is important here)

	// handle system eroor
	if err != nil {
		return err
	}

	// handle logic eroor
	if !ok {
		return fmt.Errorf("↪ (gofunc) phase: %s > host:%s > phase:%s > go:%s > param: %s", phaseName, hostName, goFunction.PhaseName, goFunction.Name, goFunction.ParamList)
	}

	// handle success
	return nil
}
