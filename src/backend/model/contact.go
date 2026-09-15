package model

import "time"

type Contact struct {
	Id          uint `gorm:"primarykey"`
	IssuedBy    int
	DiscordUser DiscordUser `gorm:"IssuedBy"`
	ServerId    int
	Server      Server
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
