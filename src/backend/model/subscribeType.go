package model

import "time"

type SubscribeType struct {
	ID        uint `gorm:"primarykey"`
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
