// File in gocore/phase/adapter.go
package phase2

import (
	"github.com/abtransitionit/gocore/viperx"
)

// Description: map the phase parameters to a map[string]any
func resolveParam(PhaseParam []string, config *viperx.CViper) map[string]any {
	if config == nil || len(PhaseParam) == 0 {
		return nil
	}

	paramMap := make(map[string]any)
	for _, key := range PhaseParam {
		if key == "" {
			continue
		}
		// fetch the value from config; could be string, list, or struct
		value := config.Get(key)
		if value != nil {
			paramMap[key] = value
		}
	}
	return paramMap
}
