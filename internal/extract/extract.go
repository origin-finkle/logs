package extract

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/hasura/go-graphql-client"
	"github.com/origin-finkle/logs/internal/extract/models"
	"github.com/origin-finkle/logs/internal/wcl"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Extract(app *kong.Context, reportIDs []string, folder string) {
	ctx := context.Background()
	if len(reportIDs) == 0 {
		end := time.Now()
		start := end.AddDate(0, 0, -14)
		var q ListReports
		err := wcl.Query(ctx, &q, map[string]interface{}{
			"guildID":   graphql.Int(516114),
			"startTime": graphql.Float(start.Unix() * 1e3),
			"endTime":   graphql.Float(end.Unix() * 1e3),
		})
		if err != nil {
			app.Fatalf("error while requesting reports: %s", err)
		}
		for _, report := range q.ReportData.Reports.Data {
			reportIDs = append(reportIDs, string(report.Code))
		}
	}

	var wg sync.WaitGroup
	for _, code := range reportIDs {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			location := fmt.Sprintf("%s/%s.json", folder, code)
			logrus.Infof("will write to file %s", location)
			file, err := os.Create(location)
			if err != nil {
				app.Fatalf("could not create file %s: %s", location, err)
			}
			defer file.Close()
			if err := doReport(file, code); err != nil {
				app.Fatalf("error while loading report %s: %s", code, err)
			}
		}(code)
	}
	wg.Wait()

}

func doReport(file io.Writer, code string) error {
	ctx := context.Background()
	var q GetReport
	err := wcl.Query(ctx, &q, map[string]interface{}{
		"code": graphql.String(code),
	})
	if err != nil {
		return err
	}

	report := q.toReport()
	fightIndex := map[string]int{}
	var g errgroup.Group
	for _, fight := range report.Fights {
		fightIndex[fight.Name]++
		if !fight.Kill {
			fight.InternalName = fmt.Sprintf("%s - Wipe %d (%.1f%%)", fight.InternalName, fightIndex[fight.Name], fight.FightPercentage)
		}
		f := fight // redefine for closure
		g.Go(func() error {
			return doReportEvents(report, f)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetIndent(" ", "    ")
	return enc.Encode(report)
}

func doReportEvents(report *models.Report, fight *models.Fight) error {
	ctx := context.Background()
	startTime := fight.StartTime
	log := logrus.WithFields(logrus.Fields{
		"report": report.Code,
		"fight":  fight.InternalName,
	})
	for {
		if startTime == 0 {
			break
		}
		var q GetReportEvents
		log.Infof("fetching events with page %d", startTime)
		data, err := wcl.QueryRaw(ctx, &q, map[string]interface{}{
			"code":      graphql.String(report.Code),
			"startTime": graphql.Float(startTime),
			"endTime":   graphql.Float(fight.EndTime),
		})
		if err != nil {
			return err
		}
		var eventsFetched struct {
			ReportData struct {
				Report struct {
					Events struct {
						Data              []json.RawMessage `json:"data"`
						NextPageTimestamp int64             `json:"nextPageTimestamp"`
					} `json:"events"`
				} `json:"report"`
			} `json:"reportData"`
		}
		if err := json.Unmarshal(*data, &eventsFetched); err != nil {
			return err
		}
		for _, event := range eventsFetched.ReportData.Report.Events.Data {
			var ef eventFilter
			if err := json.Unmarshal(event, &ef); err != nil {
				return err
			}
			switch ef.Type {
			case "combatantinfo", "applybuff", "cast":
				fight.Events = append(fight.Events, event)
			}
		}
		startTime = eventsFetched.ReportData.Report.Events.NextPageTimestamp
	}

	return nil
}

type eventFilter struct {
	Type string `json:"type"`
}
