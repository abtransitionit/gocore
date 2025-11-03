// File in gocore/syncx/runner.go
package syncx

import (
	"context"
	"fmt"
	"sync"
	"time" // Added for demonstration of a timeout context.
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
//   - ctx: one context, the same for all goroutines
//   - funcs: A slice of Funcs to be executed concurrently.
//
// Returns:
//   - A slice of errors, one for each function that failed to execute or a context cancellation error.
func RunConcurrently(ctx context.Context, funcs []Func) []error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	// Start a goroutine to log the cancellation reason.
	// This goroutine will wait until the context is done and then log the reason.
	go func() {
		<-ctx.Done()
		// Determine the reason for the cancellation and log it.
		select {
		case <-time.After(50 * time.Millisecond): // A small delay to allow other goroutines to log their status.
		default:
		}

		switch ctx.Err() {
		case context.Canceled:
			fmt.Println("INFO\tlogx/loggerZap.go:41\tContext canceled by user (e.g., via Ctrl+C).")
		case context.DeadlineExceeded:
			fmt.Println("INFO\tlogx/loggerZap.go:41\tContext canceled due to timeout.")
		default:
			fmt.Println("INFO\tlogx/loggerZap.go:41\tContext canceled for an unknown reason.")
		}
	}()

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
