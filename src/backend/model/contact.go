package model

import "time"

type Contact struct {
	ID          uint `gorm:"primarykey"`
	IssuedBy    int
	DiscordUser DiscordUser `gorm:"IssuedBy"`
	ServerID    int
	Server      Server
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
