package models

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/origin-finkle/logs/internal/wowhead"
)

type Gear struct {
	ID               int64      `json:"id"`
	Quality          int64      `json:"quality"`
	Icon             string     `json:"icon"`
	ItemLevel        int64      `json:"itemLevel"`
	PermanentEnchant *int64     `json:"permanentEnchant,omitempty"`
	TemporaryEnchant *int64     `json:"temporaryEnchant,omitempty"`
	Gems             []*GearGem `json:"gems,omitempty"`
	WowheadAttr      string     `json:"wowhead_attr"`
	UUID             string     `json:"uuid"`

	WowheadData *wowhead.Item `json:"-"`
}

var (
	slotsToEnchant = map[string]bool{
		"Tête":        true,
		"Épaule":      true,
		"Torse":       true,
		"Jambes":      true,
		"Pieds":       true,
		"Poignets":    true,
		"Mains":       true,
		"Main droite": true,
		"Deux mains":  true,
		"Dos":         true,
		"À une main":  true,
	}
	slotsWithTemporaryEnchant = map[string]bool{
		"Main droite": true,
		"Deux mains":  true,
		"À une main":  true,
	}
)

func (g Gear) ShouldBeEnchanted() bool {
	return slotsToEnchant[g.WowheadData.Slot]
}

func (g Gear) ShouldHaveTemporaryEnchant() bool {
	return slotsWithTemporaryEnchant[g.WowheadData.Slot]
}

func (g Gear) CountMissingGems() int {
	return int(g.WowheadData.Sockets) - len(g.Gems)
}

func (g *Gear) ComputeUUID() {
	g.UUID = g.WowheadAttr
	if g.TemporaryEnchant != nil {
		g.UUID += ":" + strconv.FormatInt(*g.TemporaryEnchant, 10)
	} else {
		g.UUID += ":None"
	}
}

func (g *Gear) ComputeWowheadAttr() {
	v := url.Values{}
	v.Add("domain", "fr.tbc")
	v.Add("item", strconv.FormatInt(g.ID, 10))
	if len(g.Gems) > 0 {
		gems := []string{}
		for _, gem := range g.Gems {
			gems = append(gems, strconv.FormatInt(gem.ID, 10))
		}
		v.Add("gems", strings.Join(gems, ":"))
	}
	if g.PermanentEnchant != nil {
		v.Add("ench", strconv.FormatInt(*g.PermanentEnchant, 10))
	}
	g.WowheadAttr = v.Encode()
}

func (g Gear) CountGems(color string) int64 {
	count := int64(0)
	for _, gem := range g.Gems {
		for _, gemColor := range gemColorToRealColors[gem.Color] {
			if gemColor == color {
				count++
			}
		}
	}
	return count
}

var (
	gemColorToRealColors = map[string][]string{
		"purple": {"red", "blue"},
		"red":    {"red"},
		"blue":   {"blue"},
		"yellow": {"yellow"},
		"green":  {"blue", "yellow"},
		"orange": {"red", "yellow"},
	}
)

func (g Gear) GetGems(color string) []*GearGem {
	result := make([]*GearGem, 0)
	for _, gem := range g.Gems {
		if gem.Color == color {
			result = append(result, gem)
		}
	}
	return result
}
