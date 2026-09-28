package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const freeSubscribeTypeLabel = "free"

var freeSubscribeTypeId uint

type SubscribeType struct {
	Id        uint `gorm:"primarykey"`
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func loadFreeSubscribeTypeId(ctx context.Context) error {
	subscribeType, err := gorm.G[SubscribeType](db).Where("label = ?", freeSubscribeTypeLabel).First(ctx)
	if err != nil {
		return err
	}
	freeSubscribeTypeId = subscribeType.Id
	return nil
}
