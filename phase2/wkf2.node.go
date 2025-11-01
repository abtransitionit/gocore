// File in gocore/phase/adapter.go
package phase2

import (
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/viperx"
)

// func ResolveNode(nodeName string) []string {
// 	// placeholder/mocked NodeSet mapping
// 	mock := map[string][]string{
// 		"all":      {"mockNode1", "mockNode2"},
// 		"frontend": {"mockFront1", "mockFront2"},
// 	}
// 	if nodes, ok := mock[nodeName]; ok {
// 		return nodes
// 	}
// 	// default: return nodeName as single node
// 	return []string{nodeName}
// }

// ResolveNode takes a node selector string and resolves it to a list of node names
// using the provided workflow configuration.
func ResolveNode(cfg *viperx.Config, nodeSpec string) []string {
	if cfg == nil {
		fmt.Println("⚠️  ResolveNode: config is nil, returning localhost")
		return []string{"localhost"}
	}

	nodeSpec = strings.TrimSpace(nodeSpec)
	if nodeSpec == "" {
		return []string{"localhost"}
	}

	// Try to find in config under "node"
	if cfg.IsSet("node." + nodeSpec) {
		nodes := cfg.GetStringSlice("node." + nodeSpec)
		if len(nodes) > 0 {
			return nodes
		}
	}

	// Allow comma-separated list directly in YAML or inline spec
	if strings.Contains(nodeSpec, ",") {
		parts := strings.Split(nodeSpec, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}

	// If it's a single node name not found in config, return it as-is
	return []string{nodeSpec}
}

// func (wf *Workflow2) ResolveNode(v *viper.Viper, logger logx.Logger) {
// 	// Iterate over all phases in the workflow
// 	for _, phase := range wf.Phases {
// 		// Check if the phase has a NodeSet defined
// 		if phase.Node != "" {
// 			// Resolve the NodeSet name into the actual slice of nodes using the viper config
// 			nodes := v.GetStringSlice("node." + phase.Node)

// 			// Log the mapping from NodeSet name to actual nodes
// 			logger.Infof("Phase %s targets NodeSet=%s -> %v", phase.Name, phase.Node, nodes)
// 		}
// 	}
// }
