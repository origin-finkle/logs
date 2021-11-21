package config

import (
	"fmt"

	"github.com/origin-finkle/logs/internal/models"
)

func SetConsumable(consumable *models.Consumable) {
	data.Consumables[consumable.ID] = consumable
}

func GetConsumable(id int64) (*models.Consumable, error) {
	if c, ok := data.Consumables[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("consumable for ability=%d not found", id)
}
