package testutils

import (
	"os"

	"github.com/origin-finkle/logs/internal/config"
)

func Init() {
	err := config.Init(os.Getenv("CONFIG_FOLDER"))
	if err != nil {
		panic(err)
	}
}
