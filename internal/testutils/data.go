package testutils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func LoadJSONData(t *testing.T, filename string) json.RawMessage {
	require := td.Require(t)
	file, err := os.Open(filename)
	require.CmpNoError(err)
	defer file.Close()
	var v json.RawMessage
	err = json.NewDecoder(file).Decode(&v)
	require.CmpNoError(err)
	return v
}
