package phase2

import (
	"github.com/abtransitionit/gocore/logx"
)

// Description: represents the signature of a GO function to be executed on a target.
type PhaseFn func([]string, logx.Logger) (bool, error)

// type PhaseFn func(ctx context.Context, params any, logger logx.Logger) error

// Description: represents a GO function to be executed on a Target.
type GoFunction struct {
	PhaseName string
	Name      string
	Func      PhaseFn
	ParamList []string
}

// Description: represents a map of function (that are registered and can be executed).
type FnRegistry struct {
	functionMap map[string]PhaseFn
}

// Description: define an instance of a registry as a singleton
var globalRegistry = &FnRegistry{
	functionMap: make(map[string]PhaseFn),
}

// Description: constructor that returns an instance of FnRegistry
func GetFnRegistry() *FnRegistry {
	return globalRegistry
}
