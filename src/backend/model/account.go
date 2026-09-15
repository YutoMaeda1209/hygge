package model

import "time"

type Account struct {
	ID              uint `gorm:"primarykey"`
	DiscordID       string
	EmailAddress    string
	SubscribeTypeID int
	SubscribeType   SubscribeType
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
