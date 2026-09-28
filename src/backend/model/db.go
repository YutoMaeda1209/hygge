package model

import (
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

	if err := runMigrations(); err != nil {
		return err
	}

	return nil
}
