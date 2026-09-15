package model

import "time"

type VoiceChannel struct {
	Id             uint `gorm:"primarykey"`
	ServerId       int
	Server         Server
	ChannelId      string
	NewChannelName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
