package model

import "time"

// BaseEntity 定义持久化模型通用的审计字段。
type BaseEntity struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
	CreateUser string
	UpdateUser string
}

// TenantBaseEntity 在基础审计字段上补充租户归属信息。
type TenantBaseEntity struct {
	BaseEntity
	TenantCode string
}
