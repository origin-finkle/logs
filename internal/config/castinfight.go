package config

import (
	"fmt"

	"github.com/origin-finkle/logs/internal/models"
)

func GetCastInFight(id int64) (*models.CastInFight, error) {
	if v, ok := data.CastInFight[id]; ok {
		return v, nil
	}
	return nil, ErrCastInFightNotFound
}

func SetCastInFight(v *models.CastInFight) {
	data.CastInFight[v.SpellID] = v
}

var (
	ErrCastInFightNotFound = fmt.Errorf("cast in fight not found")
)
