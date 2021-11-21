package main

import (
	"github.com/alecthomas/kong"
	"github.com/origin-finkle/logs/internal/analyze"
	"github.com/origin-finkle/logs/internal/extract"
	"github.com/sirupsen/logrus"
)

var CLI struct {
	Extract struct {
		ReportIDs []string `arg:"" optional:"" name:"report-id" help:"Report ID"`
		Folder    string   `name:"folder" help:"Folder to store data in" type:"existingdir"`
	} `cmd:"" help:"Extract reports. If no report ID is given, will try to extract reports from last 14 days"`

	Analyze analyze.Analyze `cmd:"" help:"Analyze reports. If no report ID is given, will analyze every report"`

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
	case "analyze", "analyze <report-id>":
		CLI.Analyze.Run(ctx)
	default:
		ctx.Fatalf("command %s not implemented", ctx.Command())
	}
}
