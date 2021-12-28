package elementalshaman

import (
	"context"

	"github.com/origin-finkle/logs/internal/analyze/gameplay/gameplayutils"
	"github.com/origin-finkle/logs/internal/models"
)

const spell_CHAIN_LIGHTNING = 25442

var chainLightningThreshold = map[string]int64{
	"Hydross the Unstable": 20, // because of nature phase?
	"default":              60,
	"Lady Vashj":           30, // most of the fight is waiting for elems to spawn
}

func Process(ctx context.Context, fa *models.FightAnalysis) error {
	threshold, ok := chainLightningThreshold[fa.Name]
	if !ok {
		threshold = chainLightningThreshold["default"]
	}
	if err := gameplayutils.UseAsMuchAsPossible(ctx, fa, spell_CHAIN_LIGHTNING, threshold); err != nil {
		return err
	}
	return nil
}
