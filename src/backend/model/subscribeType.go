package model

import "time"

type SubscribeType struct {
	Id        uint `gorm:"primarykey"`
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
