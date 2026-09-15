package model

import "time"

type ServerManager struct {
	ServerID  int
	Server    Server
	AccountID int
	Account   Account
	CreatedAt time.Time
	UpdatedAt time.Time
}
