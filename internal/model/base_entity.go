package model

import "time"

type BaseEntity struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
	CreateUser string
	UpdateUser string
}

type TenantBaseEntity struct {
	BaseEntity
	TenantCode string
}
