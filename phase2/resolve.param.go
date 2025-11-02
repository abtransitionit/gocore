// File in gocore/phase/adapter.go
package phase2

import (
	"github.com/abtransitionit/gocore/viperx"
)

func resolveParam(PhaseParam []string, config *viperx.CViper) map[string]interface{} {
	if config == nil || len(PhaseParam) == 0 {
		return nil
	}

	params := make(map[string]any)
	for _, key := range PhaseParam {
		if key == "" {
			continue
		}
		// fetch the value from config; could be string, list, or struct
		value := config.Get(key)
		if value != nil {
			params[key] = value
		}
	}
	return params
}
