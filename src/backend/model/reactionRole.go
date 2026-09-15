package model

import "time"

type ReactionRole struct {
	ID         uint `gorm:"primarykey"`
	ServerID   int
	Server     Server
	MessageID  string
	ReactionID string
	RoleID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
