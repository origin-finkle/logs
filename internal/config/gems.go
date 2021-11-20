package config

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/wowhead"
)

var (
	gemMu sync.RWMutex
)

func loadGems(folder string) error {
	file, err := os.Open(folder + "/gems.json")
	if err != nil {
		return err
	}
	defer file.Close()
	var gems []*models.Gem
	if err := json.NewDecoder(file).Decode(&gems); err != nil {
		return err
	}
	gemMu.Lock()
	data.Gems = make(map[int64]*models.Gem)
	gemMu.Unlock()
	for _, gem := range gems {
		SetGem(gem)
	}
	return nil
}

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
