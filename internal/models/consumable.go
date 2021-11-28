package models

type Consumable struct {
	CommonConfig

	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Description string   `json:"description,omitempty"`
}

func (c *Consumable) IsBattleElixir() bool {
	return c.Is("battle_elixir")
}

func (c *Consumable) Is(consumableType string) bool {
	for _, t := range c.Types {
		if t == consumableType {
			return true
		}
	}
	return false
}

func (c Consumable) IsGuardianElixir() bool {
	return c.Is("guardian_elixir")
}

func (c Consumable) IsFood() bool {
	return c.Is("food")
}
