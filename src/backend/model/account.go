package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	Id              uint   `gorm:"primarykey"`
	DiscordId       string `gorm:"unique"`
	EmailAddress    string
	SubscribeTypeId uint `gorm:"not null"`
	SubscribeType   SubscribeType
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (account *Account) Find(ctx context.Context, query any, args ...any) error {
	findAccount, err := gorm.G[Account](db).Where(query, args...).First(ctx)
	if err != nil {
		return err
	}
	*account = findAccount
	return nil
}

// EnsureByDiscordId creates the account if missing, tolerating concurrent logins of the same user.
func (account *Account) EnsureByDiscordId(ctx context.Context, discordId string) error {
	newAccount := Account{DiscordId: discordId, SubscribeTypeId: freeSubscribeTypeId}
	err := gorm.G[Account](db, clause.OnConflict{Columns: []clause.Column{{Name: "discord_id"}}, DoNothing: true}).Create(ctx, &newAccount)
	if err != nil {
		return err
	}
	return account.Find(ctx, "discord_id = ?", discordId)
}
