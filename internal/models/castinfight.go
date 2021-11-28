package models

type CastInFight struct {
	CommonConfig

	SpellID          int64  `json:"spell_id"`
	Display          bool   `json:"display"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	ItemID           int64  `json:"item_id,omitempty"`
	CooldownID       string `json:"cooldown_id,omitempty"`
	SuggestedSpellID int64  `json:"suggested_spell_id,omitempty"`
}
