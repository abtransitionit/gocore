package phase2

import (
	"context"
	"fmt"

	"github.com/abtransitionit/gocore/list"
	"github.com/abtransitionit/gocore/logx"
	"github.com/abtransitionit/gocore/viperx"
)

// Description: execute a workflow
func (wkf *Workflow) Execute(ctx context.Context, cfg *viperx.Viperx, fnRegistry *FnRegistry, retainSkipRange string, logger logx.Logger) error {
	// log
	logger.Infof("🅦 Runing workflow %q to %s", wkf.Name, wkf.Description)
	logger.Info("• Phases in the same tier run concurrently. Next tier starts when the previous one completes.")
	logger.Info("• Phases run concurently on each node.")

	// get the tiers
	tierList, err := wkf.TopoSortByTier(logger)
	if err != nil {
		return fmt.Errorf("cannot sort tiers: %w", err)
	}

	// filter the phases in the tiers
	tierListFiltered, err := wkf.filterPhase(logger, tierList, retainSkipRange)
	if err != nil {
		return err
	}

	// display the workflow
	phaseView, err := wkf.GetTierView(tierListFiltered, logger)
	if err != nil {
		return fmt.Errorf("getting phase table: %w", err)
	}
	list.PrettyPrintTable(phaseView)

	// loop over each tier
	nbTier := len(tierListFiltered)
	for tierId, tier := range tierListFiltered {
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

func resolveNode(PhaseNode string, config *viperx.Viperx) []string {
	if config == nil || PhaseNode == "" {
		return nil
	}
	return config.GetStringSlice("node." + PhaseNode)
}
