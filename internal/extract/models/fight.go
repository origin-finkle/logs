package models

import "encoding/json"

type Fight struct {
	EncounterID     int64             `json:"encounterID"`
	StartTime       int64             `json:"startTime"`
	EndTime         int64             `json:"endTime"`
	Name            string            `json:"name"`
	EnemyNPCs       []EnemyNPC        `json:"enemyNPCs"`
	FriendlyPlayers []int64           `json:"friendlyPlayers"`
	FightPercentage float64           `json:"fightPercentage"`
	Kill            bool              `json:"kill"`
	InternalName    string            `json:"internalName"`
	Events          []json.RawMessage `json:"events"`
}

type EnemyNPC struct {
	GameID        int64 `json:"gameID"`
	ID            int64 `json:"id"`
	InstanceCount int64 `json:"instanceCount"`
	GroupCount    int64 `json:"groupCount"`
}
