package config

import (
	"fmt"

	"github.com/origin-finkle/logs/internal/models"
)

func SetEnchant(enchant *models.Enchant) {
	data.Enchants[enchant.ID] = enchant
}

func GetEnchant(enchantID int64) (*models.Enchant, error) {
	if v, ok := data.Enchants[enchantID]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("enchant %d not found", enchantID)
}
