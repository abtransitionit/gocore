// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Description: map the phaseName to a GoFunction
func resolveFunction(PhaseFuncName string, fr *FunctionRegistry, logger logx.Logger) (*GoFunc, error) {

	// check parameter
	if fr == nil || PhaseFuncName == "" {
		return nil, fmt.Errorf("provide a PhaseFuncName and the FunctionRegistry")
	}

	// Check function is in the registry
	if !fr.Has(PhaseFuncName) {
		return nil, fmt.Errorf("PhaseFunction %q not found in registry", PhaseFuncName)
	}

	// success
	return fr.getFunction(PhaseFuncName), nil
}
