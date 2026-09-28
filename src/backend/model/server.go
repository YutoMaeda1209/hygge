package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Server struct {
	Id        uint   `gorm:"primarykey"`
	GuildId   string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// EnsureByGuildId creates the server together with a free subscription contracted by contractorId if missing.
// contractorId is ignored when the server already exists.
func (server *Server) EnsureByGuildId(ctx context.Context, guildId string, contractorId uint) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		newServer := Server{GuildId: guildId}
		err := gorm.G[Server](tx, clause.OnConflict{Columns: []clause.Column{{Name: "guild_id"}}, DoNothing: true}).Create(ctx, &newServer)
		if err != nil {
			return err
		}
		// Id stays zero when the server already existed and nothing was inserted.
		if newServer.Id == 0 {
			return nil
		}

		subscription := Subscription{
			ServerId:        newServer.Id,
			SubscribeTypeId: freeSubscribeTypeId,
			ContractorId:    contractorId,
			StartedAt:       time.Now(),
		}
		return gorm.G[Subscription](tx).Create(ctx, &subscription)
	})
	if err != nil {
		return err
	}

	findServer, err := gorm.G[Server](db).Where("guild_id = ?", guildId).First(ctx)
	if err != nil {
		return err
	}
	*server = findServer
	return nil
}
