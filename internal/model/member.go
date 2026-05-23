package model

import (
	"time"

	"gorm.io/gorm"
)

// MemberInfo 是会员表对应的持久化模型。
type MemberInfo struct {
	TenantBaseEntity
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}

// Register 描述注册接口的请求体。
type Register struct {
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"verificationCode" binding:"required"`
	TenantCode       string `json:"tenantCode"`
}

// Login 描述登录接口的请求体。
type Login struct {
	Register
}

// TableName 指定模型对应的数据库表名。
func (m MemberInfo) TableName() string {
	return "zp_member_info"
}

// CreateMember 使用给定的 GORM 连接写入会员记录。
func CreateMember(db *gorm.DB, member *MemberInfo) error {
	return db.Create(member).Error
}

// FindMemberByUsername 根据用户名和租户编号查询会员。
// 未查询到记录时返回 (nil, nil)。
func FindMemberByUsername(db *gorm.DB, username string, tenantCode string) (*MemberInfo, error) {
	var member MemberInfo
	err := db.Where("username = ?", username).
		Where("tenant_code = ?", tenantCode).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
}
