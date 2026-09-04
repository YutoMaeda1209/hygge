package main

import (
	"github.com/YutoMaeda1209/hygge/api"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	api.Api()
}
