package phase2

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute a workflow
func (wf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, logger logx.Logger) error {

	// log
	logger.Infof("🅦 Runing workflow %q to %s", wf.Name, wf.Description)
	logger.Info("Phases in the same tier run concurrently. Next tier starts when the previous one completes.")
	// success
	return nil
}
