package model

import "time"

type Account struct {
	Id              uint `gorm:"primarykey"`
	DiscordId       string
	EmailAddress    string
	SubscribeTypeId int
	SubscribeType   SubscribeType
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
