package model

import "time"

type ServerManager struct {
	ServerId  uint `gorm:"primarykey"`
	Server    Server
	AccountId uint `gorm:"primarykey"`
	Account   Account
	CreatedAt time.Time
	UpdatedAt time.Time
}
