package extract

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/hasura/go-graphql-client"
	"github.com/origin-finkle/logs/internal/common"
	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/wcl"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func filenameForReportCode(folder, code string) string {
	return fmt.Sprintf("%s/%s.json", folder, code)
}

type Extract struct {
	ReportIDs     []string `arg:"" optional:"" name:"report-id" help:"Report ID"`
	Folder        string   `name:"folder" help:"Folder to store data in" type:"existingdir"`
	CheckOnRemote bool     `name:"check-on-remote" help:"Instead of checking locally, check on remote"`
}

func (e *Extract) Extract(app *kong.Context) {
	ctx := context.Background()
	if e.Folder == "" {
		app.Fatalf("--folder is mandatory")
	}
	if len(e.ReportIDs) == 0 {
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
			ctx := logger.ContextWithLogger(ctx, logrus.WithField("report_code", report.Code))
			logrus.Debugf("Checking report %s", report.Code)
			lastActivityAgo := time.Since(time.Unix(int64(report.EndTime/1000), 0))
			if lastActivityAgo < 30*time.Minute {
				logger.FromContext(ctx).Infof("report cannot be processed, finished %s ago", lastActivityAgo)
			}
			if e.shouldExtractReport(ctx, string(report.Code)) {
				logger.FromContext(ctx).Infof("will process %s", report.Code)
				e.ReportIDs = append(e.ReportIDs, string(report.Code))
			}
		}
	}

	var wg sync.WaitGroup
	for _, code := range e.ReportIDs {
		wg.Add(1)
		go func(code string) {
			ctx := logger.ContextWithLogger(ctx, logrus.WithField("report_code", code))
			defer wg.Done()
			location := filenameForReportCode(e.Folder, code)
			logger.FromContext(ctx).Infof("will write to file %s", location)
			file, err := os.Create(location)
			if err != nil {
				app.Fatalf("could not create file %s: %s", location, err)
			}
			defer file.Close()
			if err := doReport(ctx, file, code); err != nil {
				app.Fatalf("error while loading report %s: %s", code, err)
			}
		}(code)
	}
	wg.Wait()

}

func (e *Extract) shouldExtractReport(ctx context.Context, reportCode string) bool {
	if !e.CheckOnRemote {
		location := filenameForReportCode(e.Folder, reportCode)
		if _, err := os.Stat(location); err == nil {
			// file exists, skip
			logger.FromContext(ctx).Infof("skipping %s, as it has already been processed", reportCode)
			return false
		}
		return true
	}
	url := fmt.Sprintf("https://raw.githubusercontent.com/origin-finkle/wcl-origin/master/raid-data/%s.json", reportCode)
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Warn("could not check report existency")
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Warn("could not check report existency")
		return false
	}
	logger.FromContext(ctx).Debugf("HEAD %s returned %d", url, resp.StatusCode)
	return resp.StatusCode == 404 // 404 means we don't have the file on remote
}

func doReport(ctx context.Context, file io.Writer, code string) error {
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
			ctx := logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithField("fight_name", f.InternalName))
			return doReportEvents(ctx, report, f)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	common.SetupJSONEncoder(enc)
	return enc.Encode(report)
}

func doReportEvents(ctx context.Context, report *models.Report, fight *models.Fight) error {
	startTime := fight.StartTime
	endTime := fight.EndTime
	for {
		if startTime == 0 {
			break
		}
		if startTime == fight.EndTime {
			endTime += 1
		}
		var q GetReportEvents
		logger.FromContext(ctx).Infof("fetching events with page %d", startTime)
		data, err := wcl.QueryRaw(ctx, &q, map[string]interface{}{
			"code":      graphql.String(report.Code),
			"startTime": graphql.Float(startTime),
			"endTime":   graphql.Float(endTime),
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
