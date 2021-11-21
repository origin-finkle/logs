package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/origin-finkle/logs/internal/models"
)

func ReportFilenameFromReportID(folder, reportID string) string {
	return path.Join(folder, fmt.Sprintf("%s.json", reportID))
}

func LoadReport(folder, reportID string) (*models.Report, error) {
	file, err := os.Open(ReportFilenameFromReportID(folder, reportID))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var report models.Report
	if err := json.NewDecoder(file).Decode(&report); err != nil {
		return nil, err
	}
	return &report, nil
}
