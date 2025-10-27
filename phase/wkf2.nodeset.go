// File in gocore/phase/adapter.go
package phase

import (
	"github.com/abtransitionit/gocore/logx"
	"github.com/spf13/viper"
)

// ResolveNodeSets resolves the NodeSet field of each phase into actual node lists.
// It looks up the node set name in the provided viper configuration and logs the result.
//
// Parameters:
//   - v: the viper instance containing the configuration, e.g., nodes definitions like:
//     node:
//     all: ["o1u", "o2a", "o3r"]
//     controlPlane: ["o1u"]
//     worker: ["o2a", "o3r"]
//   - logger: logger to print informational messages about the resolved nodes
//
// Usage:
//
//	workflow.ResolveNode(v, logger)
//
// After calling this, each phase with NodeSet set will have its corresponding
// actual nodes logged. You can later modify the method to store these resolved nodes
// inside the phase struct if needed for execution.
func (wf *Workflow2) ResolveNode(v *viper.Viper, logger logx.Logger) {
	// Iterate over all phases in the workflow
	for _, phase := range wf.Phases {
		// Check if the phase has a NodeSet defined
		if phase.Node != "" {
			// Resolve the NodeSet name into the actual slice of nodes using the viper config
			nodes := v.GetStringSlice("node." + phase.Node)

			// Log the mapping from NodeSet name to actual nodes
			logger.Infof("Phase %s targets NodeSet=%s -> %v", phase.Name, phase.Node, nodes)
		}
	}
}
