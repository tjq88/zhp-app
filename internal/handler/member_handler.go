package handler

import (
	"log/slog"
	"net/http"
	"zhp-app/internal/model"
	"zhp-app/internal/service"
	"zhp-app/pkg/common"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService *service.MemberService
}

func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

func (h *MemberHandler) Register(c *gin.Context) {
	var register model.Register
	if err := c.ShouldBindJSON(&register); err != nil {
		slog.Error("member_register_bind_failed", slog.String("err", err.Error()))
		fail(c, http.StatusBadRequest, "1", "invalid request")
		return
	}

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
	success(c, model.NewMemberView(member))
}
