// File in gocore/phase/adapter.go
package phase2

import (
	"github.com/abtransitionit/gocore/logx"
)

func (phase *Phase) Execute(logger logx.Logger) (string, error) {
	// log
	logger.Infof("🅟 Starting phase: %s", phase.Name)
	logger.Infof("Node: %s", phase.Node)
	logger.Infof("Fn: %s", phase.Fn)

	// success
	return "", nil
}
