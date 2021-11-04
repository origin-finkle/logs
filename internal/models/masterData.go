package models

type MasterData struct {
	Lang      string     `json:"lang"`
	Abilities []*Ability `json:"abilities"`
	Actors    []*Actor   `json:"actors"`
}

type Ability struct {
	GameID int64  `json:"gameID"`
	Type   string `json:"type"`
}

type Actor struct {
	GameID   int64  `json:"gameID"`
	Name     string `json:"name"`
	PetOwner *int64 `json:"petOwner"`
	Type     string `json:"type"`
	SubType  string `json:"subType"`
}
