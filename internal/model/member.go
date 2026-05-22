package model

import (
	"time"

	"gorm.io/gorm"
)

type MemberInfo struct {
	TenantBaseEntity
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}

// 注册对象
type Register struct {
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"verificationCode" binding:"required"`
	TenantCode       string `json:"tenantCode"`
}

func (m MemberInfo) TableName() string {
	return "zp_member_info"
}

func CreateMember(db *gorm.DB, member *MemberInfo) error {
	return db.Create(member).Error
}
