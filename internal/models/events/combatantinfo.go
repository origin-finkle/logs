package events

import (
	"context"
	"fmt"

	"github.com/origin-finkle/logs/internal/config"
	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/remark"
)

type CombatantInfo struct {
	Timestamp int64          `json:"timestamp"`
	Type      string         `json:"type"`
	SourceID  int64          `json:"sourceID"`
	Gear      []*models.Gear `json:"gear"`
	Auras     []*models.Aura `json:"auras"`
	Talents   []struct {
		ID int64 `json:"id"`
	}
}

func (ci *CombatantInfo) Process(ctx context.Context, analysis *models.Analysis, pa *models.PlayerAnalysis, fa *models.FightAnalysis) error {
	points := [3]int64{}
	for idx, t := range ci.Talents {
		points[idx] = t.ID
	}
	fa.Talents = models.NewTalents(fa, points)
	logger.FromContext(ctx).Debugf("player spec is %s", fa.Talents.Spec)
	fa.Auras = make(map[int64]*models.Aura)
	for _, aura := range ci.Auras {
		aura.Events = make([]struct{}, 0)
		fa.Auras[aura.Ability] = aura
	}
	var metaChecker func() error
	fa.Gear = make([]*models.Gear, 0)
	gc := gemCount{}
	for _, gear := range ci.Gear {
		gear.ComputeWowheadAttr()
		gear.ComputeUUID()
		if gear.ID == 0 {
			continue
		}
		fa.Gear = append(fa.Gear, gear)
		wowheadData, err := config.GetWowheadItem(ctx, gear.ID)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Debugf("could not load item %d", gear.ID)
			return err
		}
		gear.WowheadData = wowheadData
		if count := gear.CountMissingGems(); count > 0 {
			fa.AddRemark(remark.MissingGems{
				ItemWowheadAttr: gear.WowheadAttr,
				Count:           int(wowheadData.Sockets) - len(gear.Gems),
			})
		}
		for _, gem := range gear.Gems {
			gemData, err := config.GetGem(gem.ID)
			if err != nil {
				logger.FromContext(ctx).WithError(err).Debugf("could not load gem %d", gem.ID)
				return err
			}
			gem.Color = gemData.Color
			gc.Add(gem.Color)
			if gemData.IsRestricted(ctx, fa) {
				logger.FromContext(ctx).Debugf("gem %s is restricted", gemData.Name)
				fa.AddRemark(remark.InvalidGem{
					ItemWowheadAttr: gear.WowheadAttr,
					WowheadAttr:     fmt.Sprintf("item=%d", gem.ID),
				})
			}
			if gemData.Color == "meta" {
				gear := gear // for closure
				metaChecker = func() error {
					if gemData.Requires != nil {
						valid := true
						switch gemData.Requires.Rule {
						case "count_at_least":
							for _, requires := range gemData.Requires.Count {
								valid = valid && gc.HasAtLeast(requires.Color, requires.Value)
							}
						case "more_x_than_y":
							valid = gc.MoreXThanY(gemData.Requires.X, gemData.Requires.Y)
						default:
							return fmt.Errorf("meta rule %s not handled", gemData.Requires.Rule)
						}
						if !valid {
							fa.AddRemark(remark.MetaNotActivated{
								ItemWowheadAttr: gear.WowheadAttr,
								WowheadAttr:     fmt.Sprintf("item=%d", gemData.ID),
							})
						}
					}
					return nil
				}
			}
		}
		if gear.PermanentEnchant != nil {
			if enchant, err := config.GetEnchant(*gear.PermanentEnchant); err != nil {
				logger.FromContext(ctx).WithError(err).Debugf("could not load enchant %d", *gear.PermanentEnchant)
				// for now, silently ignore
				fa.AddRemark(remark.InvalidEnchant{
					ItemWowheadAttr: gear.WowheadAttr,
					Slot:            gear.WowheadData.Slot,
					EnchantID:       *gear.PermanentEnchant,
				})
			} else if enchant.IsRestricted(ctx, fa) {
				fa.AddRemark(remark.InvalidEnchant{
					ItemWowheadAttr: gear.WowheadAttr,
					Slot:            gear.WowheadData.Slot,
					EnchantID:       *gear.PermanentEnchant,
				})
			}
		} else if gear.ShouldBeEnchanted() {
			fa.AddRemark(remark.NoEnchant{
				ItemWowheadAttr: gear.WowheadAttr,
				Slot:            gear.WowheadData.Slot,
			})
		}

		if gear.TemporaryEnchant != nil {
			if enchant, err := config.GetTemporaryEnchant(*gear.TemporaryEnchant); err != nil {
				logger.FromContext(ctx).WithError(err).Debugf("could not load enchant %d", *gear.TemporaryEnchant)
				// for now, silently ignore
				fa.AddRemark(remark.InvalidTemporaryEnchant{
					ItemWowheadAttr: fmt.Sprintf("item=%d&ench=%d", gear.ID, *gear.TemporaryEnchant),
				})
			} else if enchant.IsRestricted(ctx, fa) {
				fa.AddRemark(remark.InvalidTemporaryEnchant{
					ItemWowheadAttr: fmt.Sprintf("item=%d&ench=%d", gear.ID, *gear.TemporaryEnchant),
				})
			}
		} else if gear.ShouldHaveTemporaryEnchant() {
			// check if there could be a potential windfury
			if fa.CouldBenefitFromWindfury(analysis) {
				fa.AddRemark(remark.NoTemporaryEnchantButWindfury{
					ItemWowheadAttr: fmt.Sprintf("item=%d", gear.ID),
				})
			} else {
				fa.AddRemark(remark.NoTemporaryEnchant{
					ItemWowheadAttr: fmt.Sprintf("item=%d", gear.ID),
				})
			}
		}
	}
	if metaChecker != nil {
		if err := metaChecker(); err != nil {
			return err
		}
	}
	return nil
}

func (ci *CombatantInfo) GetSource() int64 { return ci.SourceID }

var (
	gemColorToRealColors = map[string][]string{
		"purple": {"red", "blue"},
		"red":    {"red"},
		"blue":   {"blue"},
		"yellow": {"yellow"},
		"green":  {"blue", "yellow"},
		"orange": {"red", "yellow"},
	}
)

type gemCount map[string]int64

func (g gemCount) HasAtLeast(color string, value int64) bool {
	return g[color] >= value
}

func (g gemCount) Add(color string) {
	for _, c := range gemColorToRealColors[color] {
		g[c]++
	}
}

func (g gemCount) MoreXThanY(x, y string) bool {
	return g[x] > g[y]
}
