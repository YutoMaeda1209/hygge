package main

import (
	"github.com/YutoMaeda1209/hygge/api"
	"github.com/YutoMaeda1209/hygge/discord"
)

func Main() {
	go api.Api()
	go discord.Discord()
}
