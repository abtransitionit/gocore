package phase2

// var globalRegistry = &FnRegistry{
// 	functionMap: make(map[string]PhaseFn),
// }

// Description: returns an instance of FnRegistry
// func GetFnRegistry() *FnRegistry {
// 	return globalRegistry
// }

func GetFnRegistry() *FnRegistry {
	return &FnRegistry{
		functionMap: make(map[string]PhaseFn),
	}
}

// Description: adds a function to the registry
func (registry *FnRegistry) Add(name string, phaseFn PhaseFn) {
	registry.functionMap[name] = phaseFn
}

// Description: returns the function with the given name
func (registry *FnRegistry) Get(name string) (PhaseFn, bool) {
	fn, ok := registry.functionMap[name]
	return fn, ok
}

// Description: returns the names of all registered functions
func (registry *FnRegistry) List() []string {
	names := make([]string, 0, len(registry.functionMap))
	for k := range registry.functionMap {
		names = append(names, k)
	}
	return names
}

// description: check a PhaseFuncName is in the registry
func (registry *FnRegistry) Has(key string) bool {
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
