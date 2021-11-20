package analyze

import (
	"context"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/alecthomas/kong"
	"github.com/origin-finkle/logs/internal/common"
	"github.com/origin-finkle/logs/internal/config"
	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/events"
	"github.com/origin-finkle/logs/internal/models/remark"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Analyze struct {
	ReportIDs      []string `arg:"" optional:"" name:"report-id" help:"ID of the report to extract"`
	ReportsFolder  string   `name:"reports-folder" help:"Folder where reports are stored"`
	AnalysisFolder string   `name:"analysis-folder" help:"Folder where analysis are stored" type:"existingdir"`
	ConfigFolder   string   `name:"config-folder" help:"Folder where configuration are stored" type:"existingdir"`
}

func (an *Analyze) Run(app *kong.Context) {
	if an.AnalysisFolder == "" {
		app.Fatalf("--analysis-folder is mandatory")
	}
	if an.ReportsFolder == "" {
		app.Fatalf("--report-folder is mandatory")
	}
	if an.ConfigFolder == "" {
		app.Fatalf("--config-folder is mandatory")
	}
	if err := config.Init(an.ConfigFolder); err != nil {
		app.Fatalf("could not initialize config: %s", err)
	}
	if len(an.ReportIDs) == 0 {
		an.ReportIDs = make([]string, 0)
		// load all reports
		err := filepath.Walk(an.ReportsFolder, func(path string, info fs.FileInfo, err error) error {
			if strings.HasSuffix(path, ".json") {
				an.ReportIDs = append(an.ReportIDs, strings.Split(info.Name(), ".")[0])
			}
			return nil
		})
		if err != nil {
			app.Fatalf("could not find report IDs: %s", err)
		}
	}
	var g errgroup.Group
	sem := semaphore.NewWeighted(8)
	for _, reportID := range an.ReportIDs {
		reportID := reportID // prevent closure issues
		g.Go(func() error {
			if err := sem.Acquire(context.TODO(), 1); err != nil {
				return err
			}
			defer sem.Release(1)
			return an.doReport(reportID)
		})
	}
	if err := g.Wait(); err != nil {
		app.Fatalf("failed to process reports: %s", err)
	}
}

func (an *Analyze) doReport(reportID string) error {
	ctx := logger.ContextWithLogger(context.Background(), logger.FromContext(context.TODO()).WithField("report_code", reportID))
	logger.FromContext(ctx).Info("starting analysis")
	report, err := common.LoadReport(an.ReportsFolder, reportID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Debug("could not load report")
		return err
	}
	if report.Zone.ID == 0 {
		logger.FromContext(ctx).Info("logs did not contain any fight")
		return nil
	}
	reportAnalyzer := &ReportAnalyzer{
		Report:  report,
		Players: make(map[int64]*models.Player),
		Cmd:     an,
	}
	if err := reportAnalyzer.Analyze(ctx); err != nil {
		return err
	}
	logger.FromContext(ctx).Info("analysis done")
	return nil
}

type ReportAnalyzer struct {
	Report    *models.Report
	Players   map[int64]*models.Player
	playersMu sync.Mutex
	Cmd       *Analyze
	Analysis  *models.Analysis
}

func (ra *ReportAnalyzer) Analyze(ctx context.Context) error {
	logs := &models.Logs{
		Title:     ra.Report.Title,
		Code:      ra.Report.Code,
		StartTime: ra.Report.StartTime,
		EndTime:   ra.Report.EndTime,
		Fights:    map[string]models.LogFight{},
		ZoneID:    ra.Report.Zone.ID,
		Actors:    []string{},
	}
	ra.Analysis = &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	var g errgroup.Group
	for _, fight := range ra.Report.Fights {
		logs.Fights[fight.InternalName] = models.LogFight{
			StartTime: fight.StartTime,
			EndTime:   fight.EndTime,
			Name:      fight.InternalName,
		}
		f := fight // for closure
		g.Go(func() error {
			return ra.doFight(ctx, f)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	var g2 errgroup.Group
	for _, playerAnalysis := range ra.Analysis.Data {
		playerAnalysis := playerAnalysis
		ctx := logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithFields(logrus.Fields{
			"player": playerAnalysis.Actor.Name,
		})) // for closure
		g2.Go(func() error {
			var g3 errgroup.Group
			for _, fightAnalysis := range playerAnalysis.Fights {
				fightAnalysis := fightAnalysis // for closure
				ctx := logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithFields(logrus.Fields{
					"fight_name": fightAnalysis.Name,
				})) // for closure
				g3.Go(func() error {
					return postProcessAnalysis(ctx, fightAnalysis)
				})
			}
			if err := g3.Wait(); err != nil {
				return err
			}
			return nil
		})
		if err := g2.Wait(); err != nil {
			return err
		}
		playerAnalysis.AggregateRemarks()
	}
	for _, player := range ra.Players {
		logs.Actors = append(logs.Actors, player.Name)
	}
	sort.SliceStable(logs.Actors, func(x, y int) bool { return logs.Actors[x] < logs.Actors[y] })

	if err := common.SaveLogs(ra.Cmd.AnalysisFolder, logs); err != nil {
		return err
	}
	if err := common.SaveAnalysis(ra.Cmd.AnalysisFolder, ra.Report.Code, ra.Analysis); err != nil {
		return err
	}

	return nil
}

func (ra *ReportAnalyzer) setPlayerIfNeeded(ctx context.Context, playerID int64) {
	ra.playersMu.Lock()
	defer ra.playersMu.Unlock()

	if _, ok := ra.Players[playerID]; !ok {
		ra.Players[playerID] = &models.Player{Actor: *ra.Report.MasterData.Actors[playerID]}
		logger.FromContext(ctx).Debugf("found new player %d: %s", playerID, ra.Players[playerID].Name)
		ra.Analysis.SetPlayerAnalysis(playerID, &models.PlayerAnalysis{
			Actor:  *ra.Report.MasterData.Actors[playerID],
			Fights: make(map[string]*models.FightAnalysis),
		})
	}
}

func (ra *ReportAnalyzer) doFight(ctx context.Context, fight *models.Fight) error {
	ctx = logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithField("fight_name", fight.Name))
	if fight.Name == "Chess Event" {
		logger.FromContext(ctx).Debug("skipping analysis of Chess Event")
		return nil
	}
	logger.FromContext(ctx).Debugf("handling fight %s", fight.Name)
	for _, playerID := range fight.FriendlyPlayers {
		ra.setPlayerIfNeeded(ctx, playerID)
		ra.Analysis.GetPlayerAnalysis(playerID).SetFight(&models.FightAnalysis{
			Name:    fight.InternalName,
			Remarks: make([]*remark.Remark, 0),
			Casts:   make(map[int64]int64),
			Analysis: &models.TrueFightAnalysis{
				Items:       make([]models.FightCast, 0),
				Spells:      make([]models.FightCast, 0),
				Unknown:     make([]models.FightCast, 0),
				Consumables: make([]models.FightCast, 0),
			},
		})
	}
	for _, event := range fight.Events {
		if err := events.Process(ctx, event, ra.Analysis, fight.InternalName); err != nil {
			return err
		}
	}
	return nil
}
