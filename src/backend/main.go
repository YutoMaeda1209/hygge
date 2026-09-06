package main

import (
	"log/slog"
	"os"

	"github.com/YutoMaeda1209/hygge/api"
	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file is not loaded.")
	}

	api.Api()
}
