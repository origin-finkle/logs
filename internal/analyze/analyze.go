package analyze

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/kong"
	"github.com/origin-finkle/logs/internal/config"
	"github.com/origin-finkle/logs/internal/models"
	"golang.org/x/sync/errgroup"
)

func Analyze(app *kong.Context, reportIDs []string, analysisFolder, reportFolder, configFolder string) {
	if analysisFolder == "" {
		app.Fatalf("--analysis-folder is mandatory")
	}
	if reportFolder == "" {
		app.Fatalf("--report-folder is mandatory")
	}
	if configFolder == "" {
		app.Fatalf("--config-folder is mandatory")
	}
	if err := config.Init(configFolder); err != nil {
		app.Fatalf("could not initialize config: %s", err)
	}
	if len(reportIDs) == 0 {
		// load all reports
		app.Fatalf("no report id is not implemented")
	}
	var g errgroup.Group
	for _, reportID := range reportIDs {
		reportID := reportID // prevent closure issues
		g.Go(func() error {
			return doReport(reportID, reportFolder)
		})
	}
	if err := g.Wait(); err != nil {
		app.Fatalf("failed to process reports%s", err)
	}
}

func doReport(reportID string, reportFolder string) error {
	file, err := os.Open(reportFolder + "/" + reportID + ".json")
	if err != nil {
		return err
	}
	defer file.Close()
	var report models.Report
	if err := json.NewDecoder(file).Decode(&report); err != nil {
		return err
	}
	players := map[int64]*models.Player{}
	for _, fight := range report.Fights {
		if fight.Name == "Chess Event" {
			continue // nothing to report
		}
		for _, playerID := range fight.FriendlyPlayers {
			if _, ok := players[playerID]; !ok {
				players[playerID] = &models.Player{Actor: report.MasterData.Actors[playerID]}
			}
		}
		for _, event := range fight.Events {
			event = event
		}
	}
	return nil
}
