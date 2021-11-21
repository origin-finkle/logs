package models

type Enchant struct {
	CommonConfig

	ID      int64  `json:"id"`
	Name    string `json:"name"`
	SpellID int64  `json:"spell_id"` // TODO:implement
}
