package analyze

import (
	"context"

	"github.com/origin-finkle/logs/internal/analyze/gameplay"
	"github.com/origin-finkle/logs/internal/models"
)

func postProcessAnalysis(ctx context.Context, fa *models.FightAnalysis) error {
	if err := checkConsumables(ctx, fa); err != nil {
		return err
	}
	if err := processCasts(ctx, fa); err != nil {
		return err
	}
	checkTalents(fa)
	if err := checkGear(ctx, fa); err != nil {
		return err
	}
	if err := gameplay.Process(ctx, fa); err != nil {
		return err
	}
	return nil
}
