package model

import "time"

type ReactionRole struct {
	Id         uint `gorm:"primarykey"`
	ServerId   uint
	Server     Server
	MessageId  string
	ReactionId string
	RoleId     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
