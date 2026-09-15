package model

import "time"

type Function struct {
	ID        uint `gorm:"primarykey"`
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
