package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/origin-finkle/logs/internal/models"
)

func loadTemporaryEnchants(folder string) error {
	file, err := os.Open(folder + "/temporary_enchants.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var m []*models.TemporaryEnchant
	if err := json.NewDecoder(file).Decode(&m); err != nil {
		return err
	}
	data.TemporaryEnchants = make(map[int64]*models.TemporaryEnchant)
	for _, enchant := range m {
		data.TemporaryEnchants[enchant.ID] = enchant
	}
	return nil
}

func GetTemporaryEnchant(enchantID int64) (*models.TemporaryEnchant, error) {
	if v, ok := data.TemporaryEnchants[enchantID]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("temporary enchant %d not found", enchantID)
}
