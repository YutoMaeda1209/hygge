package model

import "time"

type Function struct {
	Id        uint `gorm:"primarykey"`
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
