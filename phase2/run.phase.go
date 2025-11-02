// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

func (phase *Phase) Execute(config *viperx.CViper, fr *FunctionRegistry, logger logx.Logger) (string, error) {
	// resolve member
	nodeList := resolveNode(phase.Node, config)
	paramList := resolveParam(phase.Param, config)
	nbParam := len(paramList)

	// log
	logger.Debugf("🅟 Starting Phase : %s > NodeSet:  %s (%v)", phase.Name, phase.Node, nodeList)
	// logger.Debugf("Phase : %s > Fn      reolve to:  %s", phase.Name, phase.Fn)
	if nbParam > 0 {
		logger.Debugf("Phase : %s > %d Param(s) reolve to:  %s", phase.Name, nbParam, paramList)
	}
	if fr.Has(phase.Fn) {
		logger.Debugf("Phase : %s > Function %s exists in registry", phase.Name, phase.Fn)
	} else {
		return "", fmt.Errorf("function %q not found in registry", phase.Fn)
	}
	// success
	return "", nil
}
