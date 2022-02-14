package events

import (
	"context"

	"github.com/origin-finkle/logs/internal/models"
)

type Cast struct {
	Timestamp     int64  `json:"timestamp"`
	Type          string `json:"type"`
	SourceID      int64  `json:"sourceID"`
	TargetID      int64  `json:"targetID"`
	AbilityGameID int64  `json:"abilityGameID"`
}

func (c *Cast) Process(ctx context.Context, analysis *models.Analysis, pa *models.PlayerAnalysis, fa *models.FightAnalysis) error {
	fa.AddCast(c.AbilityGameID, c.Timestamp, 1)
	return nil
}

func (c *Cast) GetSource() int64 { return c.SourceID }
