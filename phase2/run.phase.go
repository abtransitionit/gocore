// File in gocore/phase/adapter.go
package phase2

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

func (phase *Phase) Execute(ctx context.Context, config *viperx.CViper, fr *FunctionRegistry, logger logx.Logger) (string, error) {
	// resolve
	nodeList := resolveNode(phase.Node, config)
	// resolve
	paramMap := resolveParam(phase.Param, config)
	nbParam := len(paramMap)

	// log
	logger.Debugf("🅟 Starting Phase : %s > NodeSet:  %s (%v)", phase.Name, phase.Node, nodeList)
	// logger.Debugf("Phase : %s > Fn      reolve to:  %s", phase.Name, phase.Fn)
	if nbParam > 0 {
		logger.Debugf("Phase : %s > %d Param(s) reolve to:  %s", phase.Name, nbParam, paramMap)
	}

	// Execute the function
	_, err := GetGoFunc(phase.Fn).Execute(ctx, fr, nodeList, logger)
	if err != nil {
		return "", err
	}

	// success
	return "", nil
}
