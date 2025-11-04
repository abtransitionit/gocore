package phase2

import (
	"context"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute a workflow
func (wf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, logger logx.Logger) error {

	// log
	logger.Infof("🅦 Runing workflow %s to %s", wf.Name, wf.Description)

	// success
	return nil
}
