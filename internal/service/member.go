package service

import (
	"fmt"
	"log/slog"
	"time"
	"zhp-app/internal/model"
	"zhp-app/pkg/utils"

	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

type MemberService struct {
	db     *gorm.DB
	pwdKey string
}

func NewMemberService(db *gorm.DB, pwdKey string) *MemberService {
	return &MemberService{
		db:     db,
		pwdKey: pwdKey,
	}
}

func (s *MemberService) Register(register *model.Register) (*model.MemberInfo, error) {
	if register == nil {
		return nil, fmt.Errorf("register payload is nil")
	}

	if s.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

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

func maskPassword(password string) string {
	if password == "" {
		return ""
	}

	return "******"
}
