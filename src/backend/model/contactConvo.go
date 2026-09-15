package model

import "time"

type ContactConvo struct {
	ID          uint `gorm:"primarykey"`
	ContactID   int
	Contact     Contact
	SendBy      int
	DiscordUser DiscordUser `gorm:"SendBy"`
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
