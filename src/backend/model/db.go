package model

import (
	"context"

	"github.com/YutoMaeda1209/hygge/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() error {
	postgresDb, err := gorm.Open(postgres.Open(config.Conf.DatabaseUrl), &gorm.Config{})
	if err != nil {
		return err
	}
	db = postgresDb

	err = db.AutoMigrate(
		&SubscribeType{},
		&Function{},
		&GrantFuncPerm{},
		&Account{},
		&DiscordUser{},
		&Server{},
		&ServerManager{},
		&UserJoinServer{},
		&Contact{},
		&ContactConvo{},
		&ReactionRole{},
		&VoiceChannel{},
		&Session{},
	)
	if err != nil {
		return err
	}

	return ensureFreeSubscribeType(context.Background())
}
