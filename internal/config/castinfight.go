package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
)

func loadCastInFight(folder string) error {
	file, err := os.Open(folder + "/cast_in_fight.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var m map[string]*models.CastInFight
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		logger.FromContext(context.TODO()).WithError(err).Warn("could not load cast_in_fight data")
		return err
	}
	data.CastInFight = make(map[int64]*models.CastInFight)
	for k, v := range m {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return err
		}
		data.CastInFight[id] = v
		v.SpellID = id
	}
	return nil
}

func GetCastInFight(id int64) (*models.CastInFight, error) {
	if v, ok := data.CastInFight[id]; ok {
		return v, nil
	}
	return nil, ErrCastInFightNotFound
}

var (
	ErrCastInFightNotFound = fmt.Errorf("cast in fight not found")
)
