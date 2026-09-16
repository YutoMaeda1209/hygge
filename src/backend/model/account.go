package model

import "time"

type Account struct {
	Id              uint `gorm:"primarykey"`
	DiscordId       string
	EmailAddress    string
	SubscribeTypeId uint
	SubscribeType   SubscribeType
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
