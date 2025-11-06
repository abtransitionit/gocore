package phase2

import (
	"fmt"
	"sort"
)

func (wf *Workflow) topoSortByPhase() ([]Phase, error) {
	inDegree, graph, err := wf.buildDependencyGraph()
	if err != nil {
		return nil, err
	}

	// Queue for zero in-degree nodes
	queue := []string{}
	for key, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, key)
		}
	}
	sort.Strings(queue) // deterministic

	var sorted []Phase
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

// Description: topological sort by tier
//
// Notes:
//
// - skipRetainRange is optional.
// - If nothing is passed, all tiers are returned.
// - Strings must have a prefix: "r" for retain, "s" for skip (like "r1-3", "s2-4").
// - Internally, parse it and separate retain and skip lists.
func (wf *Workflow) topoSortByTier(skipRetainRange ...string) ([][]Phase, error) {
	inDegree, graph, err := wf.buildDependencyGraph()
	if err != nil {
		return nil, err
	}

	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue)

	var tiers [][]Phase
	for len(queue) > 0 {
		tierSize := len(queue)
		currentTier := make([]Phase, 0, tierSize)
		nextQueue := make([]string, 0)

		for i := 0; i < tierSize; i++ {
			name := queue[i]
			phase := wf.Phases[name]
			phase.Name = name
			currentTier = append(currentTier, phase)

			for _, neighbor := range graph[name] {
				inDegree[neighbor]--
				if inDegree[neighbor] == 0 {
					nextQueue = append(nextQueue, neighbor)
				}
			}
		}

		sort.Strings(nextQueue)
		tiers = append(tiers, currentTier)
		queue = nextQueue
	}

	// Validate for cycles
	total := 0
	for _, tier := range tiers {
		total += len(tier)
	}
	if total != len(wf.Phases) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return tiers, nil
}

// description: helper function to build the dependency graph.
func (wf *Workflow) buildDependencyGraph() (map[string]int, map[string][]string, error) {
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	for name := range wf.Phases {
		inDegree[name] = 0
		graph[name] = []string{}
	}

	for name, phase := range wf.Phases {
		for _, dep := range phase.Dependency {
			if _, ok := wf.Phases[dep]; !ok {
				return nil, nil, fmt.Errorf("dependency %q for phase %q does not exist", dep, name)
			}
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}

	return inDegree, graph, nil
}
