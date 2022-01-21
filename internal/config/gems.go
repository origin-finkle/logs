package config

import (
	"context"
	"sync"

	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/wowhead"
)

var (
	gemMu sync.RWMutex
)

func GetGem(id int64) (*models.Gem, error) {
	gemMu.RLock()
	if gem, ok := data.Gems[id]; ok {
		gemMu.RUnlock()
		return gem, nil
	}

	gemMu.RUnlock()
	gemWowhead, err := wowhead.GetGem(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	gem := &models.Gem{
		Quality: gemWowhead.Quality,
		Name:    gemWowhead.Name,
		Color:   gemWowhead.Color,
		ID:      id,
	}
	SetGem(gem)
	return gem, nil
}

func SetGem(gem *models.Gem) {
	gemMu.Lock()
	defer gemMu.Unlock()
	data.Gems[gem.ID] = gem
}
