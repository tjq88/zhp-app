package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
	"zhp-app/internal/model"
	"zhp-app/pkg/common"
	"zhp-app/pkg/config"
	"zhp-app/pkg/utils"

	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

type MemberService struct {
	db     *gorm.DB
	pwdKey string
}

var (
	// ErrMemberNotFound 表示会员不存在。
	ErrMemberNotFound = errors.New("member not found")
	// ErrInvalidPassword 表示登录密码不正确。
	ErrInvalidPassword = errors.New("invalid password")
)

// NewMemberService 基于共享基础设施状态创建会员领域服务。
func NewMemberService() *MemberService {
	return &MemberService{
		db:     common.Db,
		pwdKey: config.AppConfig.PasswordKey,
	}
}

// Register 执行会员注册业务流程：
// 校验依赖、组装持久化模型、加密密码并落库。
func (s *MemberService) Register(register *model.Register) (*model.MemberInfo, error) {
	if register == nil {
		return nil, fmt.Errorf("register payload is nil")
	}

	if s.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	// 统一使用同一个时间戳，避免审计字段和登录时间出现细小偏差。
	now := time.Now()
	member := &model.MemberInfo{
		Username: register.Username,
		Password: utils.HmacMd5(s.pwdKey, register.Username+register.Password),
		TenantBaseEntity: model.TenantBaseEntity{
			BaseEntity: model.BaseEntity{
				Id:         idgen.NextId(),
				CreateTime: now,
				UpdateTime: now,
				CreateUser: register.Username,
				UpdateUser: register.Username,
			},
			TenantCode: register.TenantCode,
		},
		LastLoginTime: now,
	}

	// 记录脱敏后的副本，避免原始密码摘要离开 service 边界。
	logMember := *member
	logMember.Password = maskPassword(logMember.Password)
	slog.Info("member_register_member_built",
		slog.Any("member", logMember),
	)

	if err := model.CreateMember(s.db, member); err != nil {
		slog.Error("member_create_persist_failed",
			slog.Int64("memberID", member.Id),
			slog.String("username", member.Username),
			slog.String("tenantCode", member.TenantCode),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	slog.Info("member_create_persisted",
		slog.Int64("memberID", member.Id),
		slog.String("username", member.Username),
		slog.String("tenantCode", member.TenantCode),
	)

	return member, nil
}

// Login 执行会员登录校验流程。
func (s *MemberService) Login(login *model.Login) (*model.MemberInfo, error) {
	if login == nil {
		return nil, fmt.Errorf("login payload is nil")
	}

	if s.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	memberInfo, err := model.FindMemberByUsername(s.db, login.Username, login.TenantCode)
	if err != nil {
		slog.Error("member_login_query_failed",
			slog.String("username", login.Username),
			slog.String("tenantCode", login.TenantCode),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	if memberInfo == nil {
		slog.Warn("member_login_member_not_found",
			slog.String("username", login.Username),
			slog.String("tenantCode", login.TenantCode),
		)
		return nil, ErrMemberNotFound
	}

	password := utils.HmacMd5(s.pwdKey, login.Username+login.Password)
	if memberInfo.Password != password {
		slog.Warn("member_login_invalid_password",
			slog.String("username", login.Username),
			slog.String("tenantCode", login.TenantCode),
		)
		return nil, ErrInvalidPassword
	}

	slog.Info("member_login_succeeded",
		slog.Int64("memberID", memberInfo.Id),
		slog.String("username", memberInfo.Username),
		slog.String("tenantCode", memberInfo.TenantCode),
	)

	return memberInfo, nil
}

// maskPassword 用于避免敏感密码内容写入日志。
func maskPassword(password string) string {
	if password == "" {
		return ""
	}

	return "******"
}
