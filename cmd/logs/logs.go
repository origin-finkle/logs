package main

import (
	"github.com/alecthomas/kong"
	"github.com/origin-finkle/logs/internal/extract"
	"github.com/sirupsen/logrus"
)

var CLI struct {
	Extract struct {
		ReportIDs []string `arg:"" optional:"" name:"report-id" help:"Report ID"`
		Folder    string   `name:"folder" help:"Folder to store data in" type:"existingdir"`
	} `cmd:"" help:"Extract reports. If no report ID is given, will try to extract reports from last 14 days"`

	Analyze struct {
		ReportID       []string `arg:"" optional:"" name:"report-id" help:"ID of the report to extract"`
		ReportsFolder  string   `name:"reports-folder" help:"Folder where reports are stored"`
		AnalysisFolder string   `name:"analysis-folder" help:"Folder where analysis are stored" type:"existingdir"`
	} `cmd:"" help:"Analyze reports. If no report ID is given, will analyze every report"`

	Verbose bool `optional:"" name:"verbose" help:"Activate debug logs"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if CLI.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	switch ctx.Command() {
	case "extract", "extract <report-id>":
		extract.Extract(ctx, CLI.Extract.ReportIDs, CLI.Extract.Folder)
	default:
		ctx.Fatalf("command %s not implemented", ctx.Command())
	}
}
