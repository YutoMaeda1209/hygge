package model

import "time"

type DiscordUser struct {
	Id        uint `gorm:"primarykey"`
	DiscordId string
	CreatedAt time.Time
	UpdatedAt time.Time
}
