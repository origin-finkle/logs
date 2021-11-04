package config

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/origin-finkle/logs/internal/models"
)

func loadCastInFight(folder string) error {
	file, err := os.Open(folder + "/cast_in_fight.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var m map[string]*models.CastInFight
	if err := json.NewDecoder(file).Decode(&m); err != nil {
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
