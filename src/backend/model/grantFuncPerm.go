package model

import "time"

type GrantFuncPerm struct {
	SubscribeTypeId int
	SubscribeType   SubscribeType
	FunctionId      int
	Function        Function
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
