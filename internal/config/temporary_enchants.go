package config

import (
	"fmt"

	"github.com/origin-finkle/logs/internal/models"
)

func SetTemporaryEnchant(enchant *models.TemporaryEnchant) {
	if data.TemporaryEnchants[enchant.ID] != nil {
		panic(fmt.Errorf("TemporaryEnchant %d is already stored in cache", enchant.ID))
	}
	data.TemporaryEnchants[enchant.ID] = enchant
}

func GetTemporaryEnchant(enchantID int64) (*models.TemporaryEnchant, error) {
	if v, ok := data.TemporaryEnchants[enchantID]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("temporary enchant %d not found", enchantID)
}
