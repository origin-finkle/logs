package models

import (
	"context"

	"github.com/origin-finkle/logs/internal/logger"
)

type Gem struct {
	CommonConfig

	ID       int64            `json:"id"`
	Name     string           `json:"name"`
	Quality  int64            `json:"quality"`
	Color    string           `json:"color"`
	Requires *GemRequirements `json:"requires,omitempty"`
}

type GemRequirements struct {
	Rule  string `json:"rule"`
	Count []struct {
		Color string `json:"color"`
		Value int64  `json:"value"`
	} `json:"count,omitempty"`
	X string `json:"x,omitempty"`
	Y string `json:"y,omitempty"`
}

func (g *Gem) IsRestricted(ctx context.Context, fa *FightAnalysis) bool {
	if g.Quality < 3 {
		logger.FromContext(ctx).Debugf("gem %s has quality %d, which is < 3", g.Name, g.Quality)
		return true
	}
	ctx = logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithField("gem_name", g.Name))
	return g.CommonConfig.IsRestricted(ctx, fa)
}
