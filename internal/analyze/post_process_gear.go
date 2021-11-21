package analyze

import (
	"context"

	"github.com/origin-finkle/logs/internal/config"
	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/remark"
)

var (
	slotsWithMandatoryItem = []string{
		"Tête",
		"Cou",
		"Épaule",
		"Torse",
		"Taille",
		"Jambes",
		"Pieds",
		"Poignets",
		"Mains",
		"Doigt",
		"Bijou",
		"Dos",
	}
)

func checkGear(ctx context.Context, fa *models.FightAnalysis) error {
	if len(fa.Gear) == 0 {
		logger.FromContext(ctx).Debug("no gear found")
		return nil
	}
	slots := map[string]int{
		"Tête":                   1,
		"Cou":                    1,
		"Épaule":                 1,
		"Torse":                  1,
		"Taille":                 1,
		"Jambes":                 1,
		"Pieds":                  1,
		"Poignets":               1,
		"Mains":                  1,
		"Doigt":                  2,
		"Bijou":                  2,
		"Dos":                    1,
		"Main droite":            1,
		"Main gauche":            1,
		"À distance":             1,
		"À une main":             2,
		"Relique":                1,
		"Deux mains":             1,
		"Tenu(e) en main gauche": 1,
		"Armes de jet":           1,
	}
	for _, gear := range fa.Gear {
		wowheadData, err := config.GetWowheadItem(ctx, gear.ID)
		if err != nil {
			return err
		}
		switch wowheadData.Slot {
		case "Tabard", "Chemise":
			continue
		}
		slots[wowheadData.Slot]--
		if slots[wowheadData.Slot] == 0 {
			delete(slots, wowheadData.Slot)
		}
	}
	for _, slot := range slotsWithMandatoryItem {
		if slots[slot] > 0 {
			fa.AddRemark(remark.MissingItemInSlot{
				Slot: slot,
			})
		}
		if slots["Relique"] > 0 && slots["Armes de jet"] > 0 && slots["À distance"] > 0 {
			fa.AddRemark(remark.MissingItemInSlot{
				Slot: "Relique/Armes de jet/À distance",
			})
		}
	}
	if slots["Main gauche"] > 0 || slots["Main droite"] > 0 || slots["Tenu(e) en main gauche"] > 0 || slots["Deux mains"] > 0 || slots["À une main"] > 0 {
		// so we need to figure out
		// possible options:
		// - Main droite + (Main gauche | Tenu(e) en main gauche | À une main)
		// - À une main + (À une main | Main gauche | Tenu(e) en main gauche)
		// - Deux mains
		valid := false
		if slots["Deux mains"] == 0 {
			valid = true
		}
		if slots["Main droite"] == 0 && ((slots["Main gauche"] > 0 || slots["Tenu(e) en main gauche"] > 0) || slots["À une main"] > 1) {
			valid = true
		}
		if slots["À une main"] <= 1 && (slots["Main gauche"] > 0 || slots["Tenu(e) en main gauche"] > 0) {
			valid = true
		}
		if !valid {
			fa.AddRemark(remark.MissingItemInSlot{
				Slot: "Armes",
			})
		}
	}
	return nil
}
