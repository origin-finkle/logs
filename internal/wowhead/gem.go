package wowhead

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type Gem struct {
	Quality int64
	Name    string
	Color   string
}

func GetGem(ctx context.Context, gemID int64) (*Gem, error) {
	item, err := getWowheadItem(ctx, gemID)
	if err != nil {
		return nil, err
	}
	quality, err := strconv.ParseInt(item.Item.Quality.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	gem := Gem{
		Quality: quality,
		Name:    item.Item.Name,
	}
	for name, color := range nameToColor {
		if strings.Contains(strings.ToLower(gem.Name), strings.ToLower(name)) {
			gem.Color = color
			break
		}
	}
	if gem.Color == "" {
		return nil, fmt.Errorf("gem %d: did not find color for %s", gemID, gem.Name)
	}
	return &gem, nil
}

var (
	nameToColor = map[string]string{
		"Rubis vivant":         "red",
		"Pierre d'aube":        "yellow",
		"Diamant tonneterre":   "meta",
		"Tanzanite":            "purple",
		"Opale de feu":         "orange",
		"Spessarite de flamme": "orange",
		"Diamant brûleciel":    "meta",
		"Draénite dorée":       "yellow",
		"Tourmaline":           "red",
	}
)
