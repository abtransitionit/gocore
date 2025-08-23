// file  gocore/ctx/ctx.go
package ctx

// Custom key type (best practice)
type contextKey string

// Keys (constants)
const (
	ExecutionIDKey contextKey = "executionID"
	WorkflowKey    contextKey = "workflow" // 👈 new one for *Workflow
)
