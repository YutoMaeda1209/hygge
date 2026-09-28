package model

import "time"

// Every server has exactly one subscription, and every subscription has a contractor, even on the free plan.
type Subscription struct {
	ServerId        uint `gorm:"primarykey;autoIncrement:false"`
	Server          Server
	SubscribeTypeId uint `gorm:"not null"`
	SubscribeType   SubscribeType
	ContractorId    uint `gorm:"not null"`
	Contractor      Account
	StartedAt       time.Time `gorm:"not null"`
	ExpiresAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
