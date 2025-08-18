package phase

// Name: SetPhase
//
// Description: a helper function to create a new Phase instance.
//
// Parameters:
//
//   - name:   A phase ID.
//   - desc:   The purpose of the phase.
//   - fn:     The function attached to the phase. It must return an error and a string.
//
// Return:
//
//   - A New Phase struct.
func SetPhase(name string, fn PhaseFunc, desc string) Phase {
	return Phase{
		Name:        name,
		Description: desc,
		fn:          fn,
	}
}

// Name: NewPhaseList
//
// Description: a constructor
//
// Parameters:
// - phases: A variadic list of Phase structs.
//
// Returns:
//   - A new PhaseList
func NewPhaseList(phases ...Phase) PhaseList {
	return PhaseList(phases)
}

// // Name: NewPhaseList
// //
// // Description: a helper function that creates a PhaseList from a variable number of Phase pointers.
// //
// // Notes:
// // - the function simplifies the creation of a PhaseList and makes the code morereadable and maintainable.
// func NewPhaseList(phases ...*Phase) PhaseList {
// 	l := make(PhaseList, len(phases))
// 	for i, p := range phases {
// 		l[i] = *p
// 	}
// 	return l
// }
