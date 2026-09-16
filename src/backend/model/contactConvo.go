package model

import "time"

type ContactConvo struct {
	Id          uint `gorm:"primarykey"`
	ContactId   uint
	Contact     Contact
	SendBy      uint
	DiscordUser DiscordUser `gorm:"foreignKey:SendBy"`
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
