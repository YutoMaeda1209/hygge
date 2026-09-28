package model

import (
	"github.com/YutoMaeda1209/hygge/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDb connects to the database at config.Conf.DatabaseUrl and applies pending migrations.
// It must be called after config.LoadConf and before any other function in this package.
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
