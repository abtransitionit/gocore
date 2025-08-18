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
// Description: denotes a function (with signature define by PhaseFunc) to be executed with meta data
//
// Notes:
// - The function is designed to accept a variable number of string arguments
type Phase struct {
	Name        string
	Description string
	fn          PhaseFunc
}

// Name: PhaseList
//
// Description: a slice of Phase structs.
//
// Notes:
//   - Can be easily iterated over to execute each phase in a defined order.
type PhaseList []Phase
