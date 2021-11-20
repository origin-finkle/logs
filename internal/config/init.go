package config

import (
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/wowhead"
)

var data = struct {
	Gems              map[int64]*models.Gem
	CastInFight       map[int64]*models.CastInFight
	Consumables       map[int64]*models.Consumable
	Wowhead           Wowhead
	Enchants          map[int64]*models.Enchant
	TemporaryEnchants map[int64]*models.TemporaryEnchant
}{
	Gems:        make(map[int64]*models.Gem),
	CastInFight: make(map[int64]*models.CastInFight),
	Consumables: make(map[int64]*models.Consumable),
	Wowhead: Wowhead{
		Items: make(map[string]*wowhead.Item),
	},
	Enchants:          make(map[int64]*models.Enchant),
	TemporaryEnchants: make(map[int64]*models.TemporaryEnchant),
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
	loadWowhead,
	loadEnchants,
	loadTemporaryEnchants,
}
