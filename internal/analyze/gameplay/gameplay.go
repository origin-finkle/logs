package gameplay

import (
	"context"

	"github.com/origin-finkle/logs/internal/analyze/gameplay/elementalshaman"
	"github.com/origin-finkle/logs/internal/models"
)

var gameplays = map[models.Specialization]func(context.Context, *models.FightAnalysis) error{
	models.Specialization_ElementalShaman: elementalshaman.Process,
}

func Process(ctx context.Context, fa *models.FightAnalysis) error {
	if fa.Talents == nil {
		return nil
	}
	if fn, ok := gameplays[fa.Talents.Spec]; ok {
		return fn(ctx, fa)
	}
	return nil
}
