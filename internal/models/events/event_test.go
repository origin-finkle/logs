package events_test

import (
	"context"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/events"
	"github.com/origin-finkle/logs/internal/testutils"
)

func TestProcessEvent_CombatantInfo(t *testing.T) {
	ev := testutils.LoadJSONData(t, "testdata/combatantinfo.json")
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(9, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(9)
	pa.SetFight(&models.FightAnalysis{Name: "test"})
	fa := pa.GetFight("test")
	err := events.Process(context.TODO(), ev, analysis, "test")
	td.CmpNoError(t, err)
	td.Cmp(t, fa.Talents.Points, [3]int64{8, 11, 42})
}
