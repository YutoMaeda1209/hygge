package model

import "time"

type GrantFuncPerm struct {
	SubscribeTypeID int
	SubscribeType   SubscribeType
	FunctionID      int
	Function        Function
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
