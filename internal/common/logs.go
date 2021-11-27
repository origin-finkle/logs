package common

import (
	"encoding/json"
	"os"
	"path"

	"github.com/origin-finkle/logs/internal/models"
)

func CreateAnalyzisFolder(folder, reportID string) error {
	if err := os.Mkdir(path.Join(folder, reportID), 0755); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func SaveLogs(folder string, logs *models.Logs) error {
	if err := CreateAnalyzisFolder(folder, logs.Code); err != nil {
		return err
	}
	file, err := os.OpenFile(path.Join(folder, logs.Code, "logs.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	SetupJSONEncoder(enc)
	return enc.Encode(logs)
}

func SaveAnalysis(folder string, reportCode string, analysis *models.Analysis) error {
	if err := CreateAnalyzisFolder(folder, reportCode); err != nil {
		return err
	}
	file, err := os.OpenFile(path.Join(folder, reportCode, "analysis.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	SetupJSONEncoder(enc)
	return enc.Encode(analysis)
}
