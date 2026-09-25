package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const freeSubscribeTypeLabel = "free"

var freeSubscribeTypeId uint

type SubscribeType struct {
	Id        uint   `gorm:"primarykey"`
	Label     string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ensureFreeSubscribeType(ctx context.Context) error {
	newSubscribeType := SubscribeType{Label: freeSubscribeTypeLabel}
	err := gorm.G[SubscribeType](db, clause.OnConflict{Columns: []clause.Column{{Name: "label"}}, DoNothing: true}).Create(ctx, &newSubscribeType)
	if err != nil {
		return err
	}

	subscribeType, err := gorm.G[SubscribeType](db).Where("label = ?", freeSubscribeTypeLabel).First(ctx)
	if err != nil {
		return err
	}
	freeSubscribeTypeId = subscribeType.Id
	return nil
}
