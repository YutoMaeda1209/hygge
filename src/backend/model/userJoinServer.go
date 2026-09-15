package model

import "time"

type UserJoinServer struct {
	DiscordUserId int
	DiscordUser   DiscordUser
	ServerId      int
	Server        Server
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
