// File in gocore/phase/adapter.go
package phase2

import (
	"github.com/abtransitionit/gocore/logx"
)

func (goFunc *GoFunc) Execute(logger logx.Logger) (string, error) {
	// log
	logger.Infof("🅕 Executing function: %s", goFunc.Name)

	// success
	return "", nil
}
