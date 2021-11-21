package models

type Logs struct {
	StartTime int64               `json:"startTime"`
	EndTime   int64               `json:"endTime"`
	Title     string              `json:"title"`
	Actors    []string            `json:"actors"`
	ZoneID    int64               `json:"zoneID"`
	Fights    map[string]LogFight `json:"fights"`
	Code      string              `json:"code"`
}

type LogFight struct {
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Name      string `json:"name"`
}
