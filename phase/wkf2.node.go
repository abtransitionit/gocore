// File in gocore/phase/adapter.go
package phase

func ResolveNodeSet(nodeName string) []string {
	// placeholder/mocked NodeSet mapping
	mock := map[string][]string{
		"all":      {"mockNode1", "mockNode2"},
		"frontend": {"mockFront1", "mockFront2"},
	}
	if nodes, ok := mock[nodeName]; ok {
		return nodes
	}
	// default: return nodeName as single node
	return []string{nodeName}
}
