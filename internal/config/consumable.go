package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
)

func loadConsumables(folder string) error {
	file, err := os.Open(folder + "/consumables.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var m []*models.Consumable
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		logger.FromContext(context.TODO()).WithError(err).Warn("could not load consumables")
		return err
	}
	data.Consumables = make(map[int64]*models.Consumable)
	for _, consumable := range m {
		data.Consumables[consumable.ID] = consumable
	}
	return nil
}

func GetConsumable(id int64) (*models.Consumable, error) {
	if c, ok := data.Consumables[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("consumable for ability=%d not found", id)
}
