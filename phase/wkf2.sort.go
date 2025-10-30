// File in gocore/phase/adapter.go
package phase

import (
	"fmt"
	"sort"
)

func (wf *Workflow2) TopoSorted2() ([]Phase2, error) {
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	// Initialize graph and in-degree map
	for key := range wf.Phases {
		inDegree[key] = 0
		graph[key] = []string{}
	}

	// Build graph based on dependencies
	for key, phase := range wf.Phases {
		for _, dep := range phase.Dependencies {
			graph[dep] = append(graph[dep], key) // dep -> key
			inDegree[key]++
		}
	}

	// Queue for zero in-degree nodes
	queue := []string{}
	for key, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, key)
		}
	}
	sort.Strings(queue) // deterministic order

	var sorted []Phase2
	for len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]

		phase := wf.Phases[key]
		phase.Name = key // populate Name from map key
		sorted = append(sorted, phase)

		for _, neighbor := range graph[key] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
		sort.Strings(queue) // deterministic
	}

	if len(sorted) != len(wf.Phases) {
		return nil, fmt.Errorf("circular dependency detected")
	}
	return sorted, nil
}
