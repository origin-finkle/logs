package gameplayutils

import (
	"context"
	"fmt"
	"math"

	"github.com/origin-finkle/logs/internal/config"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/remark"
)

func UseAsMuchAsPossible(ctx context.Context, fa *models.FightAnalysis, spellID int64, threshold int64) error {
	if threshold < 0 || threshold > 100 {
		return fmt.Errorf("invalid threshold: %d", threshold)
	}
	spell, err := config.GetCastInFight(ctx, spellID)
	if err != nil {
		return err
	}
	if spell.Cooldown == 0 {
		return fmt.Errorf("spell %s (%d) does not have any cooldown defined", spell.Name, spellID)
	}
	possibleUsages := math.Floor(float64(fa.Duration()) / float64(spell.Cooldown))
	usages := fa.Casts[spellID]
	percentage := int64(float64(usages) * 100.0 / possibleUsages)
	if percentage < threshold {
		fa.AddRemark(remark.CouldMaximizeCasts{
			SpellID:               spellID,
			PossibleCasts:         int64(possibleUsages),
			ActualCasts:           usages,
			Threshold:             threshold,
			ActualPercentageOfUse: percentage,
			SpellWowheadAttr:      fmt.Sprintf("spell=%d", spellID),
		})
	}
	return nil
}
