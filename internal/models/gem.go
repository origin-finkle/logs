package models

import (
	"context"

	"github.com/origin-finkle/logs/internal/logger"
)

type Gem struct {
	CommonConfig

	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Quality int64  `json:"quality"`
	Color   string `json:"color"`
}

func (g *Gem) IsRestricted(ctx context.Context, fa *FightAnalysis) bool {
	if g.Quality < 3 {
		logger.FromContext(ctx).Debugf("gem %s has quality %d, which is < 3", g.Name, g.Quality)
		return true
	}
	return g.CommonConfig.IsRestricted(ctx, fa)
}
