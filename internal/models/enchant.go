package models

import (
	"context"

	"github.com/origin-finkle/logs/internal/logger"
	"github.com/sirupsen/logrus"
)

type Enchant struct {
	CommonConfig

	ID      int64  `json:"id"`
	Name    string `json:"name"`
	SpellID int64  `json:"spell_id"` // TODO:implement
}

func (e *Enchant) IsRestricted(ctx context.Context, fa *FightAnalysis) bool {
	ctx = logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithFields(logrus.Fields{
		"enchant_id":   e.ID,
		"enchant_name": e.Name,
	}))
	return e.CommonConfig.IsRestricted(ctx, fa)
}
