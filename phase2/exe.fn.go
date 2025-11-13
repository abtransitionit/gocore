package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Description: execute a function on 1 target
func (goFunction *GoFunction) run(ctx context.Context, target string, logger logx.Logger) error {
	// define var
	var oko bool  // the return value of the function
	var err error // the return value of the function

	// HERE code is executed on localhost
	if target == "local" { // run code locally
		logger.Infof("↪ target:%s > function: %s > running", target, goFunction.Name)
		oko, err = goFunction.Func(goFunction.ParamList, logger) // returns a bool
		if err != nil {
			return err
		}
	}
	// HERE code is executed on a remote target
	logger.Infof("↪ target:%s > nothing yet configured", target)

	// handle error
	if err != nil {
		return err
	} else if !oko { // if !ok
		return fmt.Errorf("↪ target:%s > phase:%s > go:%s > param: %s", target, goFunction.PhaseName, goFunction.Name, goFunction.ParamList)
	}
	// handle success
	return nil
}
