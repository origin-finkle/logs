package extract

import (
	"encoding/json"

	"github.com/hasura/go-graphql-client"
	"github.com/origin-finkle/logs/internal/models"
)

type GetReport struct {
	ReportData struct {
		Report struct {
			Code      graphql.String
			EndTime   graphql.Float
			StartTime graphql.Float
			Title     graphql.String
			Zone      struct {
				ID graphql.Int
			}
			MasterData struct {
				Lang      graphql.String
				Abilities []struct {
					GameID graphql.Int `graphql:"gameID"`
					Type   graphql.String
				}
				Actors []struct {
					GameID   graphql.Int `graphql:"gameID"`
					Name     graphql.String
					PetOwner graphql.Int
					Type     graphql.String
					SubType  graphql.String
				}
			}
			Fights []struct {
				EncounterID graphql.Int `graphql:"encounterID"`
				EndTime     graphql.Float
				StartTime   graphql.Float
				Name        graphql.String
				EnemyNPCs   []struct {
					GameID        graphql.Int `graphql:"gameID"`
					ID            graphql.Int
					InstanceCount graphql.Int
					GroupCount    graphql.Int
				} `graphql:"enemyNPCs"`
				FriendlyPlayers []graphql.Int
				FightPercentage graphql.Float
				Kill            graphql.Boolean
			} `graphql:"fights(killType: Encounters)"`
		} `graphql:"report(code: $code)"`
	}
}

func (r GetReport) toReport() *models.Report {
	report := &models.Report{
		Code:      string(r.ReportData.Report.Code),
		Title:     string(r.ReportData.Report.Title),
		StartTime: int64(r.ReportData.Report.StartTime),
		EndTime:   int64(r.ReportData.Report.EndTime),
		Zone: models.Zone{
			ID: int64(r.ReportData.Report.Zone.ID),
		},
		Fights: make([]*models.Fight, 0, len(r.ReportData.Report.Fights)),
		MasterData: &models.MasterData{
			Lang:      string(r.ReportData.Report.MasterData.Lang),
			Actors:    make([]*models.Actor, 0, len(r.ReportData.Report.MasterData.Actors)),
			Abilities: make([]*models.Ability, 0, len(r.ReportData.Report.MasterData.Abilities)),
		},
	}
	for _, fight := range r.ReportData.Report.Fights {
		f := &models.Fight{
			EncounterID:     int64(fight.EncounterID),
			EndTime:         int64(fight.EndTime),
			StartTime:       int64(fight.StartTime),
			Name:            string(fight.Name),
			InternalName:    string(fight.Name),
			EnemyNPCs:       make([]models.EnemyNPC, 0, len(fight.EnemyNPCs)),
			FriendlyPlayers: make([]int64, 0, len(fight.FriendlyPlayers)),
			FightPercentage: float64(fight.FightPercentage),
			Kill:            bool(fight.Kill),
			Events:          make([]json.RawMessage, 0),
		}
		for _, id := range fight.FriendlyPlayers {
			f.FriendlyPlayers = append(f.FriendlyPlayers, int64(id))
		}
		for _, enpc := range fight.EnemyNPCs {
			f.EnemyNPCs = append(f.EnemyNPCs, models.EnemyNPC{
				GameID:        int64(enpc.GameID),
				ID:            int64(enpc.ID),
				InstanceCount: int64(enpc.InstanceCount),
				GroupCount:    int64(enpc.GroupCount),
			})
		}
		report.Fights = append(report.Fights, f)
	}
	for _, actor := range r.ReportData.Report.MasterData.Actors {
		a := &models.Actor{
			GameID:  int64(actor.GameID),
			Name:    string(actor.Name),
			Type:    string(actor.Type),
			SubType: string(actor.SubType),
		}
		if actor.PetOwner != 0 {
			po := int64(actor.PetOwner)
			a.PetOwner = &po
		}
		report.MasterData.Actors = append(report.MasterData.Actors, a)
	}
	for _, ability := range r.ReportData.Report.MasterData.Abilities {
		report.MasterData.Abilities = append(report.MasterData.Abilities, &models.Ability{
			GameID: int64(ability.GameID),
			Type:   string(ability.Type),
		})
	}
	return report
}

type GetReportEvents struct {
	ReportData struct {
		Report struct {
			Events struct {
				Data              []graphql.String
				NextPageTimestamp graphql.Float
			} `graphql:"events(startTime: $startTime, endTime: $endTime, limit: 10000)"`
		} `graphql:"report(code: $code)"`
	}
}

type ListReports struct {
	ReportData struct {
		Reports struct {
			Data []struct {
				Code      graphql.String `graphql:"code"`
				StartTime graphql.Float  `graphql:"startTime"`
				EndTime   graphql.Float  `graphql:"endTime"`
			} `graphql:"data"`
		} `graphql:"reports(guildID: $guildID, startTime: $startTime, endTime: $endTime)"`
	} `graphql:"reportData"`
}
