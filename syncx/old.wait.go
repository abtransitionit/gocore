// File in gocore/syncx/runner.go
package syncx

import (
	"context"
	"fmt"
	"time"

	"github.com/abtransitionit/gocore/logx"
)

// Name: WaitForReady
//
// Description: polls a resource until the healthFunc returns true or the context times out.
func WaitForReady(ctx context.Context, logger logx.Logger, checkInterval, firstPollDelay time.Duration, healthFunc func() (bool, error)) error {
	// delay the first poll (eg. when rebuilding a VM, let the time to the VM to start being updated)
	if firstPollDelay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(firstPollDelay):
		}
	}
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for resource: %w", ctx.Err())
		default:
			ready, err := healthFunc()
			if err != nil {
				logger.Warnf("health check failed: %v", err)
			} else if ready {
				return nil // resource is ready
			}

			time.Sleep(checkInterval)
		}
	}
}
