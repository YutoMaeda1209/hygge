package main

import (
	"log/slog"
	"os"

	"github.com/YutoMaeda1209/hygge/api"
	"github.com/YutoMaeda1209/hygge/config"
	"github.com/YutoMaeda1209/hygge/model"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	err := config.LoadConf()
	if err != nil {
		slog.Error("Failed to load environment variables.", "err", err)
		panic("Failed to load environment variables.")
	}

	err = model.InitDb()
	if err != nil {
		slog.Error("Failed to initialize the database.", "err", err)
		panic("Failed to initialize the database.")
	}

	err = api.RunApiEngine()
	if err != nil {
		slog.Error("Failed to run the api engine.", "err", err)
		panic("Failed to run the api engine.")
	}
}
