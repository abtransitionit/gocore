package phase2

import (
	"fmt"
	"sort"
	"strings"
)

// Description: adds a function to the registry - the key is "Phase:FnAlias"
func (registry *FnRegistry) Add(workflowName string, FnAlias string, phaseFn PhaseFn) {
	key := fmt.Sprintf("%s:%s", workflowName, FnAlias)
	registry.functionMap[key] = phaseFn
}

// Description: returns the go function with the given "Phase:FnAlias"
func (registry *FnRegistry) Get(workflowName, fnAlias string) (PhaseFn, bool) {
	key := fmt.Sprintf("%s:%s", workflowName, fnAlias)
	fn, ok := registry.functionMap[key]
	return fn, ok
}

// Description: returns the "Phase:FnAlias" of all registered functions
func (registry *FnRegistry) List(workflowName string) []string {
	prefix := workflowName + ":"
	nameList := make([]string, 0, len(registry.functionMap))
	for k := range registry.functionMap {
		if strings.HasPrefix(k, prefix) {
			nameList = append(nameList, strings.TrimPrefix(k, prefix))
		}
	}
	sort.Strings(nameList) // optional, keep table neat
	return nameList
}

// description: check a "Phase:FnAlias" is in the registry
func (registry *FnRegistry) Has(workflowName, fnAlias string) bool {
	key := fmt.Sprintf("%s:%s", workflowName, fnAlias)
	_, ok := registry.functionMap[key]
	return ok
}
