package phase2

import (
	"fmt"
	"sort"
	"strings"
)

// Description: returns an instance of FnRegistry as a singleton
var globalRegistry = &FnRegistry{
	functionMap: make(map[string]PhaseFn),
}

// Description: returns an instance of FnRegistry
func GetFnRegistry() *FnRegistry {
	return globalRegistry
}

// func GetFnRegistry() *FnRegistry {
// 	return &FnRegistry{
// 		functionMap: make(map[string]PhaseFn),
// 	}
// }

// Description: adds a function to the registry
func (registry *FnRegistry) Add(workflowName string, FnAlias string, phaseFn PhaseFn) {
	key := fmt.Sprintf("%s:%s", workflowName, FnAlias)
	registry.functionMap[key] = phaseFn
}

// Description: returns the function with the given name
func (registry *FnRegistry) Get(workflowName, fnAlias string) (PhaseFn, bool) {
	key := fmt.Sprintf("%s:%s", workflowName, fnAlias)
	fn, ok := registry.functionMap[key]
	return fn, ok
}

// Description: returns the names of all registered functions
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

// description: check a PhaseFnAlias is in the registry
func (registry *FnRegistry) Has(workflowName, fnAlias string) bool {
	key := fmt.Sprintf("%s:%s", workflowName, fnAlias)
	_, ok := registry.functionMap[key]
	return ok
}

// // // Manage Registry
// //
// //	func (r *FnRegistry) Add(name string, fn func(context.Context, any, logx.Logger) error) {
// //		r.funcs[name] = &PhaseFn{
// //			Name: name,
// //			Func: fn,
// //		}
// //	}
// // func (r *FnRegistry) Get(name string) (*PhaseFn, bool) {
// // 	f, ok := r.funcs[name]
// // 	return f, ok
// // }

// // List returns the names of all registered functions.
// func (r *FnRegistry) List() []string {
// 	names := make([]string, 0, len(r.funcs))
// 	for k := range r.funcs {
// 		names = append(names, k)
// 	}
// 	return names
// }

// // description: check if a PhaseFuncName is in the registry
// func (fr *FnRegistry) Has(key string) bool {
// 	_, ok := fr.funcs[key]
// 	return ok
// }
