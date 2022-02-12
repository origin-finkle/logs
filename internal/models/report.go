package models

type Report struct {
	AppVersion string      `json:"app_version"`
	Code       string      `json:"code"`
	EndTime    int64       `json:"endTime"`
	StartTime  int64       `json:"startTime"`
	Title      string      `json:"title"`
	Zone       Zone        `json:"zone"`
	MasterData *MasterData `json:"masterData"`
	Fights     []*Fight    `json:"fights"`
}

type Zone struct {
	ID int64 `json:"id"`
}
