package phase2

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
)

// Manage Registry
func (r *FnRegistry) Add(name string, fn func(context.Context, any, logx.Logger) error) {
	r.funcs[name] = &PhaseFn{
		Name: name,
		Func: fn,
	}
}
func (r *FnRegistry) Get(name string) (*PhaseFn, bool) {
	f, ok := r.funcs[name]
	return f, ok
}

// List returns the names of all registered functions.
func (r *FnRegistry) List() []string {
	names := make([]string, 0, len(r.funcs))
	for k := range r.funcs {
		names = append(names, k)
	}
	return names
}

// description: check if a PhaseFuncName is in the registry
func (fr *FnRegistry) Has(key string) bool {
	_, ok := fr.funcs[key]
	return ok
}
