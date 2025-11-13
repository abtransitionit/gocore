package phase2

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
)

// Description: execute a function
func (goFunction *GoFunction) run(ctx context.Context, target string, logger logx.Logger) error {
	logger.Debugf("↪ %s > %s > %s > %s", target, goFunction.PhaseName, goFunction.Name, goFunction.ParamList)
	_, err := goFunction.Func(goFunction.ParamList, logger)
	if err != nil {
		return err
	}
	return nil
}
