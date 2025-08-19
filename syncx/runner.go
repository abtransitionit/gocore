// File to create in gocore/syncx/runner.go
package syncx

import (
	"sync"
)

// Func represents a function that can be executed in a concurrent group.
// It returns an error to indicate failure.
type Func func() error

// Name: RunConcurrently
//
// Description: executes a slice of Funcs concurrently.
//
// Parameters:
//
//   - funcs: A slice of Funcs to be executed concurrently.
//
// Returns:
// []error: a slice of all errors encountered, or nil if no errors occurred.
func RunConcurrently(funcs []Func) []error {
	var wg sync.WaitGroup
	var mu sync.Mutex // A mutex to protect the errors slice
	var errs []error

	for _, fn := range funcs {
		wg.Add(1)
		go func(fn Func) {
			defer wg.Done()
			if err := fn(); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(fn)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errs
	}

	return nil
}
