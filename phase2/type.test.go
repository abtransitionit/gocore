package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

// Desscription: represents a node
type Node struct {
	Name string
}

// Desscription: represents a function
type FuncT func(ctx context.Context, l logx.Logger, nodeList []Node) (string, error)

// Desscription: represents a phase
type PhaseT struct {
	Name        string
	Description string
	fn          FuncT
	Dependency  []string
}

// Desscription: represents a workflow
type WorkflowT struct {
	Name   string
	Phases map[string]PhaseT
}

// Desscription: represents a registry of functions
type FunctionRegistryT struct {
	Phases map[string]PhaseT
}

// Desscription: Constructor that returns an instance of a Workflow
func GetWorkflowT() *WorkflowT {
	return &WorkflowT{
		Phases: make(map[string]PhaseT),
	}
}

// Desscription: adds a Phase to a workflow
func (w *WorkflowT) AddPhase(p PhaseT) error {
	if _, exists := w.Phases[p.Name]; exists {
		return fmt.Errorf("a phase with the name '%s' already exists", p.Name)
	}
	w.Phases[p.Name] = p
	return nil
}

// Desscription: Constructor that returns an instance of a Phase
func GetPhase(name, description string, fn FuncT, dependency []string) PhaseT {
	return PhaseT{
		Name:        name,
		Description: description,
		fn:          fn,
		Dependency:  dependency,
	}
}
