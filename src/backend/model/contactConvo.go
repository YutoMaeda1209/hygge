package model

import "time"

type ContactConvo struct {
	Id          uint `gorm:"primarykey"`
	ContactId   int
	Contact     Contact
	SendBy      int
	DiscordUser DiscordUser `gorm:"SendBy"`
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
