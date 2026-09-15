package model

import "time"

type UserJoinServer struct {
	DiscordUserID int
	DiscordUser   DiscordUser
	ServerID      int
	Server        Server
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
