package main

import (
	"log/slog"
	"os"

	"github.com/YutoMaeda1209/hygge/api"
	"github.com/YutoMaeda1209/hygge/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	err := config.LoadConf()
	if err != nil {
		panic(err)
	}

	api.Api()
}
