package analyze

import (
	"context"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/events"
	"github.com/origin-finkle/logs/internal/models/remark"
	"github.com/origin-finkle/logs/internal/testutils"
)

func TestAnalyze_MissingItemInSlot_MissingWrists(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/missing_item_in_slot.json")
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
	err = checkGear(context.TODO(), fa)
	td.CmpNoError(t, err)
	missingItems := make([]*remark.Remark, 0)
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_MissingItemInSlot {
			missingItems = append(missingItems, r)
			td.Cmp(t, r.Slot, "Poignets")
		}
	}
	td.CmpLen(t, missingItems, 1, "missing_item_in_slot remark not found")
}

func TestAnalyze_MissingItemInSlot_MissingWeapon(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/missing_item_in_slot_weapon.json")
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
	err = checkGear(context.TODO(), fa)
	td.CmpNoError(t, err)
	missingItems := make([]*remark.Remark, 0)
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_MissingItemInSlot {
			missingItems = append(missingItems, r)
			td.Cmp(t, r.Slot, "Armes")
		}
	}
	td.CmpLen(t, missingItems, 1, "missing_item_in_slot remark not found")
}
