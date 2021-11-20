package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/origin-finkle/logs/internal/models"
)

func loadEnchants(folder string) error {
	file, err := os.Open(folder + "/enchants.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var m []*models.Enchant
	if err := json.NewDecoder(file).Decode(&m); err != nil {
		return err
	}
	data.Enchants = make(map[int64]*models.Enchant)
	for _, enchant := range m {
		data.Enchants[enchant.ID] = enchant
	}
	return nil
}

func GetEnchant(enchantID int64) (*models.Enchant, error) {
	if v, ok := data.Enchants[enchantID]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("enchant %d not found", enchantID)
}
