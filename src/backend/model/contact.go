package model

import "time"

type Contact struct {
	Id          uint `gorm:"primarykey"`
	IssuedBy    uint
	DiscordUser DiscordUser `gorm:"foreignKey:IssuedBy"`
	ServerId    uint
	Server      Server
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
