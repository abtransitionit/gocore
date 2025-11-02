// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
)

func resolveFunction(PhaseFn []string, registry *FunctionRegistry) string {
	return ""
}

// description: creates an empty function registry
func GetFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		funcs: make(map[string]interface{}),
	}
}

// description: registers a new function in the registry
func (r *FunctionRegistry) Add(phaseFunction string, fn any) error {
	if phaseFunction == "" || fn == nil {
		return fmt.Errorf("invalid function or name")
	}
	if _, exists := r.funcs[phaseFunction]; exists {
		return fmt.Errorf("function %q already registered", phaseFunction)
	}
	r.funcs[phaseFunction] = fn
	return nil
}

func (fr *FunctionRegistry) Has(key string) bool {
	_, ok := fr.funcs[key]
	return ok
}

// GetPhaseFunc returns a callable function from the registry
// func (fr *FunctionRegistry) GetGoFunc(PhaseFunction string) (GoFunc, error) {
// func GetPhaseFunc(fnName string) (func() (string, error), error) {
// 	raw, ok := Registry[fnName]
// 	if !ok {
// 		return nil, fmt.Errorf("fn %q not found in registry", fnName)
// 	}

// 	// wrap different simple function types into uniform func() (string, error)
// 	switch f := raw.(type) {
// 	case func() string:
// 		return func() (string, error) {
// 			return f(), nil
// 		}, nil

// 	case func() (string, error):
// 		return f, nil

// 	default:
// 		return nil, fmt.Errorf("unsupported function type for %q", fnName)
// 	}
// }

// description: returns a callable GoFunc from the registry
// func (fr *FunctionRegistry) GetGoFunc(PhaseFunction string, params map[string]any) (GoFunc, error) {
// 	raw, ok := fr.funcs[PhaseFunction]
// 	if !ok {
// 		return nil, fmt.Errorf("function %q not found in registry", PhaseFunction)
// 	}

// 	case func([]string) GoFunc:
// 		var args []string
// 		if p, ok := params["param"].([]string); ok {
// 			args = p
// 		}
// 		return f(args), nil

// 	default:
// 		return nil, fmt.Errorf("unsupported function type for %q", PhaseFunction)
// 	}
// }

// switch f := raw.(type) {
// case func(ctx context.Context, logger logx.Logger, targets []Target, cmd ...string) (string, error):
// 	return func(ctx context.Context, logger logx.Logger, targets []Target) (string, error) {
// 		var cmd []string
// 		if p, ok := params["cmd"].([]string); ok {
// 			cmd = p
// 		}
// 		return f(ctx, logger, targets, cmd...)
// 	}, nil
