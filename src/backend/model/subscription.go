package model

import "time"

type Subscription struct {
	ServerId        uint `gorm:"primarykey;autoIncrement:false"`
	Server          Server
	SubscribeTypeId uint
	SubscribeType   SubscribeType
	ContractorId    uint
	Contractor      Account
	StartedAt       time.Time
	ExpiresAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
