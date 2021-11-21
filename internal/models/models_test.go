package models_test

import (
	"os"
	"testing"

	"github.com/origin-finkle/logs/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.Init()
	os.Exit(m.Run())
}
