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
		slog.Error("Failed to load environment variables.", "err", err)
		panic("Failed to load environment variables.")
	}

	api.RunApiEngine()
}
