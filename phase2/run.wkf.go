// File in gocore/phase/adapter.go
package phase2

import (
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

func (wf *Workflow) Execute(cfg *viperx.CViper, logger logx.Logger) error {
	// log
	logger.Infof("• Starting workflow: %s", wf.Name)
	logger.Info("• Phases in the same tier run in parallel. Next tier starts when the previous one completes")

	// toposort the phases of the workflow
	phases, _ := wf.TopoPhaseSorted()

	// Loop over sorted phases
	for _, phase := range phases {
		// Execute the phase
		phase.Execute(logger)
	}
	return nil
}
