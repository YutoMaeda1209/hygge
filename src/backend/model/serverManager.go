package model

import "time"

type ServerManager struct {
	ServerId  int
	Server    Server
	AccountId int
	Account   Account
	CreatedAt time.Time
	UpdatedAt time.Time
}
