package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
)

func loadTemporaryEnchants(folder string) error {
	file, err := os.Open(folder + "/temporary_enchants.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var m []*models.TemporaryEnchant
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		logger.FromContext(context.TODO()).WithError(err).Warn("could not load temporary enchants")
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
