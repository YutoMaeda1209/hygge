package model

import "time"

// Must match the id seeded in migrations/000002_seed_free_subscribe_type.up.sql.
const freeSubscribeTypeId uint = 1

type SubscribeType struct {
	Id        uint `gorm:"primarykey"`
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
