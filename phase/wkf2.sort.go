// File in gocore/phase/adapter.go
package phase

import (
	"fmt"
	"sort"
)

func (wf *Workflow2) TopoSorted2() ([]PhaseYAML, error) {
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	// Initialize graph and in-degree map
	for name := range wf.Phases {
		inDegree[name] = 0
		graph[name] = []string{}
	}

	// Build graph based on dependencies
	for name, phase := range wf.Phases {
		for _, dep := range phase.Dependencies {
			graph[dep] = append(graph[dep], name) // dep -> name
			inDegree[name]++
		}
	}

	// Initialize queue with nodes having zero in-degree
	queue := []string{}
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue) // deterministic order

	var sorted []PhaseYAML
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		sorted = append(sorted, wf.Phases[name])

		for _, neighbor := range graph[name] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
		sort.Strings(queue) // deterministic ordering
	}

	if len(sorted) != len(wf.Phases) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return sorted, nil
}
