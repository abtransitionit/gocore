// File to create in gocore/syncx/runner.go
package syncx

import (
	"context"
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
func RunConcurrently(ctx context.Context, funcs []Func) []error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, fn := range funcs {
		// Stop if the context is already canceled.
		if ctx.Err() != nil {
			return []error{ctx.Err()}
		}

		wg.Add(1)
		go func(fn Func) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// Do not run the function if the context is canceled.
				return
			default:
				if err := fn(); err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
				}
			}
		}(fn)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		// Context was canceled; return the cancellation error.
		return []error{ctx.Err()}
	case <-done:
		// All goroutines finished; return any collected errors.
		if len(errs) > 0 {
			return errs
		}
		return nil
	}
}
