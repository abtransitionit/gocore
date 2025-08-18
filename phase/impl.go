package phase

import (
	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/logx"
)

// Name: Run
//
// Description: Executes the function associated with the Phase.
//
// Notes:
// - This function is the primary way to execute a Phase.
// - It returns the same values as the underlying `Func`.
func (p *Phase) Run() (string, error) {
	output, err := p.Func()
	if err != nil {
		// handle FAILURE
		return "", errorx.Wrap(err, "phase '%s' failed", p.Name)
	}
	return output, nil
}

// Name: Run
//
// Description: Executes all phases of a PhaseList sequentially.
//
// Notes:
//
//   - It is the main orchestrator for the entire workflow. It uses a logger
//     to provide clear progress updates and wraps errors using the errorx
//     library to provide a rich error chain with stack traces.
func (pl PhaseList) Run(l logx.Logger) error {
	l.Info("=== Starting sequence of phases ===")
	for _, p := range pl {
		l.Info("--- Starting Phase: %s ---", p.Name)
		output, err := p.Run()
		if err != nil {
			// handle FAILURE
			l.ErrorWithStack(err, "Sequence failed on phase '%s'.", p.Name)
			return err
		}
		l.Info("--- Phase '%s' completed successfully. ---", p.Name)
		l.Info("Output: %s", output)
	}
	l.Info("=== Sequence of phases completed successfully ===")
	return nil
}

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
//   - A pointer to the created Phase.
func SetPhase(name string, fn func(cmd ...string) (string, error), desc string) *Phase {
	return &Phase{Name: name, Func: fn, Description: desc}
}

// Name: NewPhaseList
//
// Description: a helper function that creates a PhaseList from a variable number of Phase pointers.
//
// Notes:
// - the function simplifies the creation of a PhaseList and makes the code morereadable and maintainable.
func NewPhaseList(phases ...*Phase) PhaseList {
	l := make(PhaseList, len(phases))
	for i, p := range phases {
		l[i] = *p
	}
	return l
}
