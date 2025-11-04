// File in gocore/phase/adapter.go
package phase2

import (
	"github.com/abtransitionit/gocore/viperx"
)

// Description: map the phase parameters to a []string
func resolveNode(PhaseNode string, config *viperx.CViper) []string {
	if config == nil || PhaseNode == "" {
		return nil
	}
	return config.GetStringSlice("node." + PhaseNode)
}
