package syncx

import (
	"sync"
)

// Name: Func
//
// Description: represents a function that can be executed in a concurrent group.
//
// Parameters:
//
//   - none
//
// Returns:
//   - error: Returns an error if the function fails.
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
//
//   - error: Returns an error if any of the functions fails.
func RunConcurrently(funcs []Func) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(funcs))

	for _, fn := range funcs {
		wg.Add(1)
		go func(fn Func) {
			defer wg.Done()
			if err := fn(); err != nil {
				errs <- err
			}
		}(fn)
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return <-errs
	}

	return nil
}

// TODO:
// - Implement a function to return all errors.
// - Implement a function to add a context with a timeout/deadline.
