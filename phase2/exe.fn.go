package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Description: execute a function on 1 target
func (goFunction *GoFunction) run(ctx context.Context, target string, logger logx.Logger) error {
	logger.Debugf("↪ target:%s > phase:%s > go:%s > param: %s", target, goFunction.PhaseName, goFunction.Name, goFunction.ParamList)
	status, err := goFunction.Func(goFunction.ParamList, logger)
	// handle error
	if err != nil {
		return err
	} else if !status {
		return fmt.Errorf("in goFunction.run > %s returned false", goFunction.Name)
	}
	// handle success
	return nil
}
