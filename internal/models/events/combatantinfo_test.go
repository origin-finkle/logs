package events_test

import (
	"context"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/events"
	"github.com/origin-finkle/logs/internal/models/remark"
	"github.com/origin-finkle/logs/internal/testutils"
)

func TestCombatantInfo_MetaNotActivated(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo_meta_not_activated.json")
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(27, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(27)
	pa.SetFight(&models.FightAnalysis{Name: "test"})
	fa := pa.GetFight("test")
	err := events.Process(context.TODO(), ev, analysis, "test")
	td.CmpNoError(t, err)
	found := false
	td.CmpLen(t, fa.Remarks, td.Gt(0))
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_MetaNotActivated {
			found = true
			break
		}
	}
	td.CmpTrue(t, found, "found meta_not_activated")
}

func TestCombatantInfo_NoEnchant(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo_no_enchant.json")
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(12, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(12)
	pa.SetFight(&models.FightAnalysis{Name: "test"})
	fa := pa.GetFight("test")
	err := events.Process(context.TODO(), ev, analysis, "test")
	td.CmpNoError(t, err)
	found := false
	td.CmpLen(t, fa.Remarks, td.Gt(0))
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_NoEnchant {
			found = true
			td.Cmp(t, r.Slot, "Tête")
			break
		}
	}
	td.CmpTrue(t, found, "found meta_not_activated")
}

func TestCombatantInfo_MissingGems(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo_missing_gems.json")
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(4, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(4)
	pa.SetFight(&models.FightAnalysis{Name: "test"})
	fa := pa.GetFight("test")
	err := events.Process(context.TODO(), ev, analysis, "test")
	td.CmpNoError(t, err)
	found := false
	td.CmpLen(t, fa.Remarks, td.Gt(0))
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_MissingGems {
			found = true
			td.Cmp(t, r.ItemWowheadAttr, "domain=fr.tbc&gems=28461&item=28193")
			break
		}
	}
	td.CmpTrue(t, found, "found missing_gems")
}

func TestCombatantInfo_InvalidGem(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo_invalid_gem.json")
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(10, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(10)
	pa.SetFight(&models.FightAnalysis{Name: "test"})
	fa := pa.GetFight("test")
	err := events.Process(context.TODO(), ev, analysis, "test")
	td.CmpNoError(t, err)
	found := false
	td.CmpLen(t, fa.Remarks, td.Gt(0))
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_InvalidGem {
			found = true
			td.Cmp(t, r.ItemWowheadAttr, "domain=fr.tbc&ench=2657&gems=23095%3A24058&item=28608")
			break
		}
	}
	td.CmpTrue(t, found, "found invalid_gem")
}

func TestCombatantInfo_ComplexRestriction(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo_warlock_leotheras.json")
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(26, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(26)
	pa.SetFight(&models.FightAnalysis{Name: "Leotheras the Blind"})
	fa := pa.GetFight("Leotheras the Blind")
	err := events.Process(context.TODO(), ev, analysis, "Leotheras the Blind")
	td.CmpNoError(t, err)
	for _, r := range fa.Remarks {
		td.CmpNot(t, r.Type, remark.Type_InvalidEnchant, "no invalid_enchant remark found")
	}
}

func TestCombatantInfo_ComplexRestrictionShouldTrigger(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo_warlock_leotheras.json")
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(26, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(26)
	pa.SetFight(&models.FightAnalysis{Name: "should trigger"})
	fa := pa.GetFight("should trigger")
	err := events.Process(context.TODO(), ev, analysis, "should trigger")
	td.CmpNoError(t, err)
	found := false
	td.CmpLen(t, fa.Remarks, td.Gt(0))
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_InvalidEnchant && r.Slot == "Mains" {
			found = true
			td.Cmp(t, r.ItemWowheadAttr, "domain=fr.tbc&ench=2613&item=30764")
			break
		}
	}
	td.CmpTrue(t, found, "found invalid_enchant")
}
