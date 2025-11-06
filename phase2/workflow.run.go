package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute a workflow
func (wkf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, retainRange string, skipRange string, logger logx.Logger) error {
	// check range parameter
	if retainRange != "" && skipRange != "" {
		return fmt.Errorf("only one of retainRange or skipRange can be set, not both")
	}

	// log
	logger.Infof("🅦 Runing workflow %q to %s", wkf.Name, wkf.Description)
	logger.Info("• Phases in the same tier run concurrently. Next tier starts when the previous one completes.")

	// Get the provided range AND log
	var rangeVal string
	if retainRange != "" {
		rangeVal = retainRange
		logger.Info("• workflow running with retained phase(s)")
	} else if skipRange != "" {
		rangeVal = skipRange
		logger.Info("• workflow running with skipped phase(s)")
	}

	// // display the workflow
	// phaseView, err := wkf.GetTierView(logger, rangeVal)
	// if err != nil {
	// 	return fmt.Errorf("getting phase table: %w", err)
	// }
	// list.PrettyPrintTable(phaseView)
	// os.Exit(0)

	// get the tiers
	tiers, err := wkf.topoSortByTier(logger, rangeVal)
	if err != nil {
		return fmt.Errorf("cannot sort tiers: %w", err)
	}

	// loop over each tier
	nbTier := len(tiers)
	for tierId, tier := range tiers {
		tierIdx := tierId + 1
		nbPhase := len(tier)
		logger.Infof("👉 Starting Tier %d / %d with %d concurent phase(s)", tierIdx, nbTier, nbPhase)

		for _, p := range tier {
			logger.Debugf("   ↪ would run concurrently: phase %q (node=%s, fn=%s)", p.Name, p.Node, p.FnAlias)
		}

		logger.Infof("✔ Tier %d complete. Waiting for next tier...", tierIdx)
	} // tier loop

	// success
	return nil
}
