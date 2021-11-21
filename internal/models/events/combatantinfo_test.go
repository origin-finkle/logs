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
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo.json")
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
	for _, r := range fa.Remarks {
		if r.Type == remark.Type_MetaNotActivated {
			found = true
			break
		}
	}
	td.CmpTrue(t, found)
}
