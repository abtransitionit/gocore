package phase

// Name: Phase
//
// Description: Phase represents a single, executable step in a workflow.
//
// It is a concrete type that encapsulates a Go function, its name,
// and a description of its purpose.  and return a string result
// along with an error for robust failure handling.
//
// Notes:
// - The function is designed to accept a variable number of string arguments
type Phase struct {
	Name        string
	Description string
	Func        func(arg ...string) (string, error) // go function, code, that perfomrs actions and accepts 0..N arguments
}

// Name: PhaseList
//
// Description: a slice of Phase structs.
//
// Notes:
//   - Can be easily iterated over to execute each phase in a defined order.
type PhaseList []Phase
