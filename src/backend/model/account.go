package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	Id           uint `gorm:"primarykey"`
	DiscordId    string
	EmailAddress string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Find loads the first account matching query into account.
// It returns gorm.ErrRecordNotFound if no account matches.
func (account *Account) Find(ctx context.Context, query any, args ...any) error {
	findAccount, err := gorm.G[Account](db).Where(query, args...).First(ctx)
	if err != nil {
		return err
	}
	*account = findAccount
	return nil
}

// EnsureByDiscordId loads the account with discordId into account, creating it first if it does not exist.
// It is safe to call concurrently with the same discordId.
func (account *Account) EnsureByDiscordId(ctx context.Context, discordId string) error {
	newAccount := Account{DiscordId: discordId}
	err := gorm.G[Account](db, clause.OnConflict{Columns: []clause.Column{{Name: "discord_id"}}, DoNothing: true}).Create(ctx, &newAccount)
	if err != nil {
		return err
	}
	return account.Find(ctx, "discord_id = ?", discordId)
}
