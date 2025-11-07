package phase2

import (
	"github.com/abtransitionit/gocore/logx"
)

// Description: represents a GO function.
type PhaseFn func([]string, logx.Logger) (bool, error)

// Description: represents a set of function (that are registered and can be executed).
type FnRegistry struct {
	functionMap map[string]PhaseFn
}
