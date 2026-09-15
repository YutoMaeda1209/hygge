package model

import "time"

type VoiceChannel struct {
	ID             uint `gorm:"primarykey"`
	ServerID       int
	Server         Server
	ChannelID      string
	NewChannelName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
