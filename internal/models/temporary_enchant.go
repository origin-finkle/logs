package models

type TemporaryEnchant struct {
	CommonConfig

	ID          int64  `json:"id"`
	Description string `json:"description"`
	SpellID     int64  `json:"spellID"`
}
