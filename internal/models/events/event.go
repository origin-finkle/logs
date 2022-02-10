package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
)

type eventType struct {
	Type string `json:"type"`
}

type Event interface {
	Process(context.Context, *models.Analysis, *models.PlayerAnalysis, *models.FightAnalysis) error
	GetSource() int64
}

func Process(ctx context.Context, ev json.RawMessage, analysis *models.Analysis, fightName string) error {
	var eType eventType
	if err := json.Unmarshal(ev, &eType); err != nil {
		return err
	}
	var event Event
	switch eType.Type {
	case "combatantinfo":
		event = &CombatantInfo{}
	case "cast":
		event = &Cast{}
	case "applybuff":
		event = &ApplyBuff{}
	default:
		return fmt.Errorf("unknown event type: %s", eType.Type)
	}
	if err := json.Unmarshal(ev, event); err != nil {
		return err
	}
	playerID := event.GetSource()
	pa := analysis.GetPlayerAnalysis(playerID)
	if pa == nil {
		return nil // event should not be processed as player does not exist
	}
	fa := pa.GetFight(fightName)
	if fa == nil {
		// sometimes fight is not found, can be due to MC or whatever, just skip it
		return nil
	}
	ctx = logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithField("player", pa.Actor.Name))
	return event.Process(ctx, analysis, pa, fa)
}
