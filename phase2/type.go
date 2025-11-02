// File in gocore/phase/type.go
package phase2

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
)

// Description: represents a set of phases.
type Workflow struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Phases      map[string]Phase `yaml:"phases"`
}

// Description: represents a phase
type Phase struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Fn          string   `yaml:"fn"`
	Dependency  []string `yaml:"dependency,omitempty"` // replace Next
	Param       []string `yaml:"param,omitempty"`
	Node        string   `yaml:"node,omitempty"`
}

// type PhaseFunc func(ctx context.Context, node []string, l logx.Logger) (string, error)
// type GoFunc func(ctx context.Context, l logx.Logger) error

type GoFunc struct {
	Name string
	Func func(ctx context.Context, logger logx.Logger) error
}
