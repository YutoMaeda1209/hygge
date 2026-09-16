package model

import "time"

type GrantFuncPerm struct {
	SubscribeTypeId uint `gorm:"primarykey"`
	SubscribeType   SubscribeType
	FunctionId      uint `gorm:"primarykey"`
	Function        Function
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
