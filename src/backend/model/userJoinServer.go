package model

import "time"

type UserJoinServer struct {
	DiscordUserId uint `gorm:"primarykey"`
	DiscordUser   DiscordUser
	ServerId      uint `gorm:"primarykey"`
	Server        Server
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
