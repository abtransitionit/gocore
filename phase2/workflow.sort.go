package phase2

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

// Description: topological sort by phase
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
func (wf *Workflow) TopoSortByTier(logger logx.Logger) ([][]Phase, error) {
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

func (wf *Workflow) filterPhase(logger logx.Logger, tierList [][]Phase, skipRetainRange string) ([][]Phase, error) {

	if skipRetainRange == "" {
		return tierList, nil
	}
	logger.Infof("• workflow phases fltering activated with : %s", skipRetainRange)
	mode := ""
	ranges := ""

	if strings.HasPrefix(skipRetainRange, "-r") {
		mode = "retain"
		ranges = skipRetainRange[2:]
	} else if strings.HasPrefix(skipRetainRange, "-s") {
		mode = "skip"
		ranges = skipRetainRange[2:]
	} else {
		return nil, fmt.Errorf("invalid filter flag: %s", skipRetainRange)
	}

	// Parse ranges WITHOUT a helper
	allowed := make(map[int]struct{})
	for _, p := range strings.Split(ranges, ",") {
		if strings.Contains(p, "-") {
			b := strings.Split(p, "-")
			if len(b) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(strings.TrimSpace(b[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(b[1]))
			if err1 != nil || err2 != nil || start > end {
				continue
			}
			for i := start; i <= end; i++ {
				allowed[i] = struct{}{}
			}
		} else {
			if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				allowed[n] = struct{}{}
			}
		}
	}

	// Filter
	var filtered [][]Phase
	globalIndex := 1

	for _, tier := range tierList {
		var newTier []Phase
		for _, ph := range tier {

			_, exists := allowed[globalIndex]

			include := false
			switch mode {
			case "retain":
				include = exists
			case "skip":
				include = !exists
			default:
				return nil, fmt.Errorf("unknown mode: %s", mode)
			}

			if include {
				newTier = append(newTier, ph)
			}

			globalIndex++
		}

		if len(newTier) > 0 {
			filtered = append(filtered, newTier)
		}
	}

	return filtered, nil
}
