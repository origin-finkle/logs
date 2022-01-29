package models_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/origin-finkle/logs/internal/models"
)

func TestFightAnalysis_FightName(t *testing.T) {
	fa := &models.FightAnalysis{
		Name: "Illidan Stormrage",
	}
	td.Cmp(t, fa.FightName(), "Illidan Stormrage")
	fa.Name = "Illidan Stormrage - Wipe 1 (66.0%)"
	td.Cmp(t, fa.FightName(), "Illidan Stormrage")
}
