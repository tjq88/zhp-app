package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"zhp-app/internal/model"
	"zhp-app/internal/service"
	"zhp-app/pkg/common"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService *service.MemberService
}

// NewMemberHandler 创建会员相关接口的 HTTP 处理器。
func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

// Register 处理会员注册请求：
// 绑定入参、补充上下文、调用 service，并返回安全的响应对象。
func (h *MemberHandler) Register(c *gin.Context) {
	var register model.Register
	if err := c.ShouldBindJSON(&register); err != nil {
		slog.Error("member_register_bind_failed", slog.String("err", err.Error()))
		fail(c, http.StatusBadRequest, "1", "invalid request")
		return
	}

	// 租户信息由中间件从请求头注入到上下文。
	register.TenantCode = c.GetString(common.TenantCode)
	slog.Info("member_register_requested",
		slog.String("username", register.Username),
		slog.String("tenantCode", register.TenantCode),
		slog.Bool("hasVerificationCode", register.VerificationCode != ""),
	)

	member, err := h.memberService.Register(&register)
	if err != nil {
		slog.Error("member_register_failed",
			slog.String("username", register.Username),
			slog.String("tenantCode", register.TenantCode),
			slog.String("err", err.Error()),
		)
		fail(c, http.StatusInternalServerError, "1", "register failed")
		return
	}

	slog.Info("member_register_succeeded",
		slog.Int64("memberID", member.Id),
		slog.String("username", member.Username),
		slog.String("tenantCode", member.TenantCode),
	)

	// 不直接返回持久化模型，避免把敏感字段暴露给调用方。
	success(c, model.NewMemberView(member))
}

func (h *MemberHandler) Login(c *gin.Context) {
	var login model.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		slog.Error("login_bind_failed", slog.String("err", err.Error()))
		fail(c, http.StatusBadRequest, "1", "invalid request")
		return
	}
	login.TenantCode = c.GetString(common.TenantCode)
	memberInfo, i := h.memberService.Login(&login)
	if memberInfo == nil {
		fail(c, http.StatusInternalServerError, strconv.Itoa(i), "login failed")
		return
	}
	success(c, model.NewMemberView(memberInfo))
}
