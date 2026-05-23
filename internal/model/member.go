package model

import (
	"log/slog"
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

// Login 描述注册接口的请求体。
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

// SelectByUsername 根据用户名查询
func SelectByUsername(db *gorm.DB, username string, tenantCode string) *MemberInfo {
	var member MemberInfo
	db.Where("username = ?", username).Where("tenant_code = ?", tenantCode).First(&member)

	if member.Username != "" {
		logMember := NewMemberView(&member)
		slog.Info("根据用户名查询用户返回数据 参数",
			slog.String("username", username),
			slog.String("tenantCode", tenantCode),
			slog.Any("member", logMember),
		)
		return &member
	}
	slog.Info("根据用户名查询用户不存在 参数",
		slog.String("username", username),
		slog.String("tenantCode", tenantCode),
	)

	return nil
}
