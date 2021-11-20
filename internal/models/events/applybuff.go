package events

import (
	"context"

	"github.com/origin-finkle/logs/internal/models"
)

type ApplyBuff struct {
	Timestamp     int64  `json:"timestamp"`
	Type          string `json:"type"`
	SourceID      int64  `json:"sourceID"`
	TargetID      int64  `json:"targetID"`
	AbilityGameID int64  `json:"abilityGameID"`
}

func (ab *ApplyBuff) Process(context.Context, *models.Analysis, *models.PlayerAnalysis, *models.FightAnalysis) error {
	return nil
}

func (ab *ApplyBuff) GetSource() int64 { return ab.SourceID }
