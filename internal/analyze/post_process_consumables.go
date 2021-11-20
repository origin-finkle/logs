package analyze

import (
	"context"
	"fmt"

	"github.com/origin-finkle/logs/internal/config"
	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/remark"
)

var (
	consumableExcludeList = map[string]bool{
		"Chess Event": true,
	}
	remarkForConsumableType = map[string]struct {
		Missing func() remark.Metadata
		Invalid func(id int64) remark.Metadata
	}{
		"food": {
			Missing: func() remark.Metadata { return remark.MissingFood{} },
			Invalid: func(id int64) remark.Metadata {
				return remark.InvalidFood{
					WowheadAttr: fmt.Sprintf("spell=%d", id),
				}
			},
		},
		"battle_elixir": {
			Missing: func() remark.Metadata { return remark.MissingBattleElixir{} },
			Invalid: func(id int64) remark.Metadata {
				return remark.InvalidBattleElixir{
					WowheadAttr: fmt.Sprintf("spell=%d", id),
				}
			},
		},
		"guardian_elixir": {
			Missing: func() remark.Metadata { return remark.MissingGuardianElixir{} },
			Invalid: func(id int64) remark.Metadata {
				return remark.InvalidGuardianElixir{
					WowheadAttr: fmt.Sprintf("spell=%d", id),
				}
			},
		},
	}
)

func checkConsumables(ctx context.Context, fa *models.FightAnalysis) error {
	if consumableExcludeList[fa.Name] {
		logger.FromContext(ctx).Debugf("fight %s is excluded from consumable checks", fa.Name)
		return nil
	}
	missing := map[string]bool{
		"food":            true,
		"battle_elixir":   true,
		"guardian_elixir": true,
	}
	invalid := map[string]int64{}
	for _, aura := range fa.Auras {
		consumable, err := config.GetConsumable(aura.Ability)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Debug("consumable not found, skipping")
			continue
		}
		for consumableType := range missing {
			if consumable.Is(consumableType) {
				delete(missing, consumableType)
				if consumable.IsRestricted(ctx, fa) {
					invalid[consumableType] = consumable.ID
				}
			}
		}
	}
	for consumableType := range missing {
		fa.AddRemark(remarkForConsumableType[consumableType].Missing())
	}
	for invalidType, id := range invalid {
		fa.AddRemark(remarkForConsumableType[invalidType].Invalid(id))
	}
	return nil
}
