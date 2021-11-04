package models

type Consumable struct {
	CommonConfig

	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Description string   `json:"description"`
}
