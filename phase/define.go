package phase

import "fmt"

// Name: NewPhase
// Description: a constructor
//
// Parameters:
//   - name: The name of the phase.
//   - description: The description of the phase.
//   - fn: The function to be executed when the phase is run.
//   - dependencies: A list of phase names that this phase depends on.
//
// Return:
//   - Phase: A new Phase object with the specified field.
//
// Notes:
//   - It handles the initialization of the unexported `fn` field internally.
func NewPhase(name, description string, fn PhaseFunc, dependencies []string) Phase {
	return Phase{
		Name:         name,
		Description:  description,
		fn:           fn,
		Dependencies: dependencies,
	}
}

// Name: NewWorkflow
//
// Description: a constructor
//
// Return:
//
// - Workflow: a new Workflow instance.
func NewWorkflow() *Workflow {
	return &Workflow{
		Phases: make(map[string]Phase),
	}
}

// Name: AddPhase
//
// Description: adds a new Phase to the workflow.
//
// Parameters:
//
//   - p: The Phase to be added.
//
// Return:
//
//   - error: If the phase with the same name already exists.
//
// Notes:
//   - This method handles the registration of a phase and its dependencies within the workflow.
func (w *Workflow) AddPhase(p Phase) error {
	if _, exists := w.Phases[p.Name]; exists {
		// handle specific error explicitly: expected outcome: The phase already exists
		return fmt.Errorf("a phase with the name '%s' already exists", p.Name)
	}
	w.Phases[p.Name] = p
	return nil
}

// Name: NewWorkflowFromPhases
//
// Description: a helper function that creates a new Workflow instance
//
// Parameters:
//
// - phases: A variadic list of Phase structs to be added to the workflow.
//
// Returns:
//
// - *Workflow: The newly created workflow, or nil if an error occurred.
// - error: An error if a phase could not be added to the workflow.
// func NewWorkflowFromPhases(phases ...Phase) *Workflow {
// 	workflow := NewWorkflow()
// 	for _, p := range phases {
// 		// We drop the error handling here.
// 		// If AddPhase can fail, you might want to log or panic instead.
// 		_ = workflow.AddPhase(p)
// 	}
// 	return workflow
// }

func NewWorkflowFromPhases(phases ...Phase) (*Workflow, error) {
	workflow := NewWorkflow()
	for _, p := range phases {
		if err := workflow.AddPhase(p); err != nil {
			// return nil, err
			// handle specific error explicitly: expected outcome: The phase already exists
			return nil, fmt.Errorf("failed to add phase %q: %w", p.Name, err)
		}
	}
	return workflow, nil
}
