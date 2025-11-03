// File in gocore/syncx/runner.go
package syncx

import (
	"context"
	"fmt"
	"sync"

	"github.com/abtransitionit/gocore/logx"
)

// Description: runs a function concurrently for each item in a slice.
//
// Parameters:
//   - itemList: a slice of items to process concurrently.
//   - fn: a function to be executed concurrently for each item.
//   - logger: a logger instance for logging errors.
//
// ConcurrentExec runs a set of tasks concurrently and collects errors.
// Each task is a function returning error.
func ExecConcurrently[T any](context context.Context, itemList []T, logger logx.Logger, fn func(T) error) error {
	// check parameter
	if len(itemList) == 0 {
		return nil
	}

	// channel to collect as errors as goroutines - one per item in itemList
	errCh := make(chan error, len(itemList))
	// channel to collect as errors of each goroutines - one per phaseitem in item
	var wg sync.WaitGroup

	// loop over each item in the items
	for _, item := range itemList {
		// increment the WaitGroup counter - signal that a new goroutine is starting
		wg.Add(1)
		item := item // capture loop variable
		// launch a goroutine - one per item in the list - all goroutines will run concurrently
		go func() {
			// decrement the WaitGroup counter - when the goroutine completes
			defer wg.Done()
			// execute the function's code
			if err := fn(item); err != nil {
				// way to get each goroutine's error (nil or error)
				errCh <- err
			}
		}()
	}

	// Wait for all phases in this tier to completes (mechanism of WaitGroup's counter)
	wg.Wait()
	// close the channel - signal that no more error will be sent
	close(errCh)

	// loop over the channel to log any error reported by the goroutines
	var errs []error
	for e := range errCh {
		logger.Errorf(e.Error())
		errs = append(errs, e)
	}

	// if an error occurred in this tier, return
	if len(errs) > 0 {
		return fmt.Errorf("%d error(s) occurred", len(errs))
	}

	// success
	return nil
}
