package config

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"strconv"

	"github.com/origin-finkle/logs/internal/logger"
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
	for _, loader := range cfg {
		if err := loader.Load(folder); err != nil {
			return err
		}
	}
	return nil
}

func Teardown(folder string) error {
	for _, loader := range cfg {
		if err := loader.Teardown(folder); err != nil {
			return err
		}
	}
	return nil
}

type loader struct {
	Name     string
	Filename string
	Decode   func(*json.Decoder) error
	Encode   func(*json.Encoder) error
}

func (l loader) Load(folder string) error {
	file, err := os.Open(path.Join(folder, l.Filename))
	if err != nil {
		return err
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := l.Decode(dec); err != nil {
		logger.FromContext(context.TODO()).WithError(err).Warn("could not load %s", l.Name)
		return err
	}
	return nil
}

func (l loader) Teardown(folder string) error {
	if l.Encode == nil {
		// no need to save
		return nil
	}
	file, err := os.OpenFile(path.Join(folder, l.Filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent(" ", "    ")
	if err := l.Encode(enc); err != nil {
		logger.FromContext(context.TODO()).WithError(err).Warnf("could not save %s", l.Name)
		return err
	}
	return nil
}

var cfg = []loader{
	{
		Name:     "gems",
		Filename: "gems.json",
		Decode: func(dec *json.Decoder) error {
			var gems []*models.Gem
			if err := dec.Decode(&gems); err != nil {
				return err
			}
			for _, gem := range gems {
				SetGem(gem)
			}
			return nil
		},
		Encode: func(enc *json.Encoder) error {
			gemMu.Lock()
			defer gemMu.Unlock()

			gems := make([]*models.Gem, 0, len(data.Gems))
			for _, gem := range data.Gems {
				gems = append(gems, gem)
			}
			return enc.Encode(gems)
		},
	},
	{
		Name:     "cast_in_fight",
		Filename: "cast_in_fight.json",
		Decode: func(dec *json.Decoder) error {
			var m map[string]*models.CastInFight
			if err := dec.Decode(&m); err != nil {
				return err
			}
			for k, v := range m {
				id, err := strconv.ParseInt(k, 10, 64)
				if err != nil {
					return err
				}
				v.SpellID = id
				SetCastInFight(v)
			}
			return nil
		},
	},
	{
		Name:     "consumables",
		Filename: "consumables.json",
		Decode: func(dec *json.Decoder) error {
			var m []*models.Consumable
			if err := dec.Decode(&m); err != nil {
				return err
			}
			for _, consumable := range m {
				SetConsumable(consumable)
			}
			return nil
		},
	},
	{
		Name:     "wowhead",
		Filename: "wowhead.json",
		Decode: func(dec *json.Decoder) error {
			return dec.Decode(&data.Wowhead)
		},
		Encode: func(enc *json.Encoder) error {
			data.Wowhead.mu.Lock()
			defer data.Wowhead.mu.Unlock()

			return enc.Encode(data.Wowhead) //nolint:copylocks
		},
	},
	{
		Name:     "enchants",
		Filename: "enchants.json",
		Decode: func(dec *json.Decoder) error {
			var m []*models.Enchant
			if err := dec.Decode(&m); err != nil {
				return err
			}
			for _, enchant := range m {
				SetEnchant(enchant)
			}
			return nil
		},
	},
	{
		Name:     "temporary_enchants",
		Filename: "temporary_enchants.json",
		Decode: func(dec *json.Decoder) error {
			var m []*models.TemporaryEnchant
			if err := dec.Decode(&m); err != nil {
				return err
			}
			for _, enchant := range m {
				SetTemporaryEnchant(enchant)
			}
			return nil
		},
	},
}
