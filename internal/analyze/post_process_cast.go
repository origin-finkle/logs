package analyze

import (
	"context"
	"fmt"
	"sort"

	"github.com/origin-finkle/logs/internal/config"
	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/remark"
)

func processCasts(ctx context.Context, fa *models.FightAnalysis) error {
	for spellID, count := range fa.Casts {
		fc := models.FightCast{
			SpellID: spellID,
			Count:   count,
		}
		castInFight, err := config.GetCastInFight(spellID)
		if err != nil {
			if err == config.ErrCastInFightNotFound {
				fa.Analysis.Unknown = append(fa.Analysis.Unknown, fc)
				continue
			}
			return err
		}
		if castInFight.IsRestricted(ctx, fa) {
			switch castInFight.InvalidReason {
			case "cast_higher_rank_available":
				fa.AddRemark(remark.CastHigherRankAvailable{
					SpellID:                      spellID,
					SuggestedSpellID:             castInFight.SuggestedSpellID,
					SpellWowheadAttr:             fmt.Sprintf("spell=%d", spellID),
					HigherRankedSpellWowheadAttr: fmt.Sprintf("spell=%d", castInFight.SuggestedSpellID),
					Count:                        int(count),
				})
			}
		}
		if !castInFight.Display {
			logger.FromContext(ctx).Debugf("cast_in_fight %d should not be displayed", spellID)
			continue
		}
		switch castInFight.Type {
		case "spell":
			fa.Analysis.Spells = append(fa.Analysis.Spells, fc)
		case "consumable":
			fa.Analysis.Consumables = append(fa.Analysis.Consumables, fc)
		case "item":
			fa.Analysis.Items = append(fa.Analysis.Items, fc)
		default:
			return fmt.Errorf("unknown cast_in_fight type %s", castInFight.Type)
		}
	}
	sort.SliceStable(fa.Analysis.Spells, func(i, j int) bool { return fa.Analysis.Spells[i].SpellID < fa.Analysis.Spells[j].SpellID })
	sort.SliceStable(fa.Analysis.Consumables, func(i, j int) bool { return fa.Analysis.Consumables[i].SpellID < fa.Analysis.Consumables[j].SpellID })
	sort.SliceStable(fa.Analysis.Items, func(i, j int) bool { return fa.Analysis.Items[i].SpellID < fa.Analysis.Items[j].SpellID })
	sort.SliceStable(fa.Analysis.Unknown, func(i, j int) bool { return fa.Analysis.Unknown[i].SpellID < fa.Analysis.Unknown[j].SpellID })
	return nil
}
