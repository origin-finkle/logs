package config

import (
	"encoding/json"
	"os"

	"github.com/origin-finkle/logs/internal/models"
)

func loadConsumables(folder string) error {
	file, err := os.Open(folder + "/consumables.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var m []*models.Consumable
	if err := json.NewDecoder(file).Decode(&m); err != nil {
		return err
	}
	data.Consumables = make(map[int64]*models.Consumable)
	for _, consumable := range m {
		data.Consumables[consumable.ID] = consumable
	}
	return nil
}
