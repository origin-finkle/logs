package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/origin-finkle/logs/internal/models"
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
	data.Gems = make(map[int64]*models.Gem)
	for _, gem := range gems {
		SetGem(gem)
	}
	return nil
}

func GetGem(id int64) (*models.Gem, error) {
	if gem, ok := data.Gems[id]; ok {
		return gem, nil
	}
	return nil, fmt.Errorf("gem not found")
}

func SetGem(gem *models.Gem) {
	data.Gems[gem.ID] = gem
}
