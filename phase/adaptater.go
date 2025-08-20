// File to create in gocore/phase/adapter.go
package phase

import (
	"context"

	"github.com/abtransitionit/gocore/syncx"
)

// Name: adaptToSyncxFunc
//
// Description: takes a PhaseFunc and returns a syncx.Func.
// Parameters:
//
//   - fn: The PhaseFunc to be adapted.
//   - ctx: The context for the phase's execution.
//   - cmd: Additional arguments to be passed to the phase's function.
//
// Returns:
//
//   - syncx.Func: A syncx.Func that represents the adapted PhaseFunc.
//
// Notes:
//   - This acts as an adapter, making a PhaseFunc compatible with the syncx.RunConcurrently function's signature.
//   - It wraps the PhaseFunc in a syncx.Func: we create a closure that is an anonymous function.
//
// Todo
//   - pass logging to the closure that is an anonymous function.
//   - make that function not an anonymous.
func adaptToSyncxFunc(fn PhaseFunc, ctx context.Context, cmd ...string) syncx.Func {
	return func() error {
		_, err := fn(ctx, cmd...)
		return err
	}
}
