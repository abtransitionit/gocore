package phase

import "context"

// Name: PhaseFunc
//
// Description: a type that represents a function to be executed as a phase.
//
// Notes:
// - The function is designed to accept a context and a variable number of string arguments
type PhaseFunc func(ctx context.Context, cmd ...string) (string, error)

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
