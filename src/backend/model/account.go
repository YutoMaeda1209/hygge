package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Id              uint   `gorm:"primarykey"`
	DiscordId       string `gorm:"unique"`
	EmailAddress    string
	SubscribeTypeId uint
	SubscribeType   SubscribeType
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (account *Account) Find(ctx context.Context, query interface{}, args ...interface{}) error {
	findAccount, err := gorm.G[Account](db).Where(query, args).First(ctx)
	account = &findAccount
	return err
}

func (account *Account) Create(ctx context.Context) error {
	return gorm.G[Account](db).Create(ctx, account)
}
