package config

import (
	"github.com/origin-finkle/logs/internal/models"
)

var data struct {
	Gems        map[int64]*models.Gem
	CastInFight map[int64]*models.CastInFight
	Consumables map[int64]*models.Consumable
}

func Init(folder string) error {
	for _, loader := range loaders {
		if err := loader(folder); err != nil {
			return err
		}
	}
	return nil
}

type configLoader func(string) error

var loaders = []configLoader{
	loadGems,
	loadCastInFight,
	loadConsumables,
}
