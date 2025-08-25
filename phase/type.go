// File in gocore/phase/type.go
package phase

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
)

// Name: PhaseFunc
//
// Description: a type that represents a function to be executed as a phase.
//
// Parameters:
//   - ctx: The context for the phase's execution.
//   - cmd: Additional arguments to be passed to the phase's function.
// 	- targets: A slice of Target to be passed to the phase's function.

// Notes:
// - The function is designed to play some code on a Target (VM, Container, etc).
type PhaseFunc func(ctx context.Context, l logx.Logger, targets []Target, cmd ...string) (string, error)

// Name: Phase
//
// Description: denotes a task to be executed within a worflow
//
// Notes:
// - The function is designed to accept a variable number of string arguments
// - primarily a function to be executed
type Phase struct {
	Name         string
	Description  string
	fn           PhaseFunc
	Dependencies []string
}

// Name: Workflow
//
// Description: represents a set of phases with defined dependencies.
//
// Notes:
//   - Can be easily iterated over to execute each phase in a defined order.
//   - primarily a map of Phases/tasks.
//   - designed to be a Directed Acyclic Graph (DAG) of tasks.
type Workflow struct {
	Name   string
	Phases map[string]Phase
}

// // Name: PhaseList
// //
// // Description: denotes a worflow.
// //
// // Notes:
// //   - Can be easily iterated over to execute each phase in a defined order.
// //   - primarily a sequence of Phases/tasks.
// type PhaseList []Phase

// Name: PhaseTiers
//
// Description: represents a set of phases with defined dependencies.
//
// Notes:
//   - Can be easily iterated over to execute each phase in a defined order.
//   - primarily a map of Phases/tasks.
//   - designed to be a Directed Acyclic Graph (DAG) of tasks.
type PhaseTiers [][]Phase

// Name: Target
//
// Description: represents an abstract entity that a phase can operate on (VM, Container).
//
// Notes:
//   - It can be a VM, container, or anything else the workflow targets.
type Target interface {
	Name() string
	Type() string // e.g., "VM", "Container"
	// Optional: add methods like Address(), SSHConfig() if needed
}

// Name: VM
//
// Description: represents a virtual machine in the workflow.
type Vm struct {
	NameStr string
	Addr    string // optional: IP or hostname
	// Add more fields if needed, e.g., SSH config, OS type, etc.
}

// Name: Name
//
// Description: returns the VM's name (implements Target interface)
func (v Vm) Name() string {
	return v.NameStr
}

// name: Type
//
// Description: returns the type of target: "VM"
func (v Vm) Type() string {
	return "Vm"
}
