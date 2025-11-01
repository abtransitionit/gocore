// File in gocore/phase/type.go
package phase

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
)

// Description: represents a set of phases.
type Workflow2 struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Phases      map[string]Phase2 `yaml:"phases"`
}

// Description: represents a phase
type Phase2 struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Fn           string            `yaml:"fn"`
	Dependencies []string          `yaml:"dependencies,omitempty"` // replace Next
	Params       map[string]string `yaml:"params,omitempty"`
	Node         string            `yaml:"node,omitempty"`
}

// Name: PhaseFunc
//
// Description: represents a function of a phase.
//
// Parameters:
//   - ctx: The context for the phase's execution.
//   - cmd: Additional arguments to be passed to the phase's function.
// 	- targets: A slice of Target to be passed to the phase's function.

// Notes:
// - The function is designed to play some code on a Target (VM, Container, etc).
// - The cmd...string here is meant to pass the same arguments to all phases of a workflow via Execute
type PhaseFunc2 func(ctx context.Context, node []string, l logx.Logger) (string, error)
type PhaseFunc3 func(ctx context.Context, l logx.Logger) error
