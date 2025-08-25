// file  gocore/ctx/ctx.go
package ctx

import (
	"context"
	"os"
	"os/signal"
)

// Custom key type (best practice)
type contextKey string

// Keys (constants)
const (
	StringKeyId   contextKey = "stringId"
	WorkflowKeyId contextKey = "workflowId"
)

// Name: NewPhaseCtx
//
// Description: creates a base context for workflow phases
//
// Parameters:
//
// - stringID: string that identifies the phase
// - workflowID: workflow metadata
//
// Returns:
//
//   - context.Context, cancel function
//
// Example Usage:
//
//	ctx, cancel := NewPhaseCtx("exec-123", workflow)
//	defer cancel()
//
// Notes:
//
// - sets up cancellation on Ctrl+C
// - injects stringId and workflowId into context
func NewPhaseCtx(stringID string, workflowID any) (context.Context, context.CancelFunc) {
	baseCtx := context.Background()
	// wrap context with user interrupt signal support
	ctxWithSignal, cancel := signal.NotifyContext(baseCtx, os.Interrupt)
	// add workflow metadata
	ctxWithSignal = context.WithValue(ctxWithSignal, StringKeyId, stringID)
	ctxWithSignal = context.WithValue(ctxWithSignal, WorkflowKeyId, workflowID)

	return ctxWithSignal, cancel
}
