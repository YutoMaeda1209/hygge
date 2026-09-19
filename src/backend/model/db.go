package model

import (
	"log/slog"

	"github.com/YutoMaeda1209/hygge/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() {
	if postgresDb, err := gorm.Open(postgres.Open(config.Conf.DatabaseUrl), &gorm.Config{}); err != nil {
		slog.Error("Failed to connect to the database.", "err", err)
		panic("Failed to connect to the database.")
	} else {
		db = postgresDb
	}
}
