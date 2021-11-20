package models

type Enchant struct {
	CommonConfig

	ID   int64  `json:"id"`
	Name string `json:"name"`
}
