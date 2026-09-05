package main

import (
	"log/slog"
	"os"

	"github.com/YutoMaeda1209/hygge/api"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	api.Api()
}
