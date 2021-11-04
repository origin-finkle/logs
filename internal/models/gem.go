package models

type Gem struct {
	CommonConfig

	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Quality int64  `json:"quality"`
	Color   string `json:"color"`
}
