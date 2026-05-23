package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"zhp-app/internal/model"
	"zhp-app/internal/service"
	"zhp-app/pkg/common"
)

type GameHandler struct {
	platformService *service.GamePlatformService
}

// NewGameHandler 创建游戏相关接口处理器。
func NewGameHandler(platformService *service.GamePlatformService) *GameHandler {
	return &GameHandler{
		platformService: platformService,
	}
}

// StartGame 处理开始游戏请求：
// 绑定参数、补齐上下文、校验平台状态，并返回开始游戏结果。
func (h *GameHandler) StartGame(c *gin.Context) {
	var gameStartReq model.GameStartReq
	if err := c.ShouldBindJSON(&gameStartReq); err != nil {
		slog.Error("game_start_bind_failed", slog.String("err", err.Error()))
		fail(c, http.StatusBadRequest, "1", "invalid request")
		return
	}

	// 租户和用户信息优先从中间件上下文补齐，避免信任客户端直接传入。
	gameStartReq.TenantCode = c.GetString(common.TenantCode)
	gameStartReq.Username = c.GetString(common.Username)
	gameStartReq.UserId = fmt.Sprint(c.Get(common.UserId))
	gameStartReq.Ip = c.ClientIP()

	slog.Info("game_start_requested",
		slog.String("platformCode", gameStartReq.PlatformCode),
		slog.Int("gameType", int(gameStartReq.GameType)),
		slog.String("tenantCode", gameStartReq.TenantCode),
		slog.String("username", gameStartReq.Username),
	)

	resp, err := h.platformService.StartGame(&gameStartReq)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGamePlatformNotFound):
			fail(c, http.StatusNotFound, "2001", "game platform not found")
		case errors.Is(err, service.ErrGamePlatformDisabled):
			fail(c, http.StatusForbidden, "2002", "game platform disabled")
		case errors.Is(err, service.ErrGamePlatformMaintaining):
			fail(c, http.StatusForbidden, "2003", "game platform maintaining")
		default:
			slog.Error("game_start_failed",
				slog.String("platformCode", gameStartReq.PlatformCode),
				slog.Int("gameType", int(gameStartReq.GameType)),
				slog.String("tenantCode", gameStartReq.TenantCode),
				slog.String("err", err.Error()),
			)
			fail(c, http.StatusInternalServerError, "1", "start game failed")
		}
		return
	}

	slog.Info("game_start_succeeded",
		slog.String("platformCode", resp.PlatformCode),
		slog.String("platformName", resp.PlatformName),
		slog.String("tenantCode", gameStartReq.TenantCode),
		slog.String("username", gameStartReq.Username),
	)
	success(c, resp)
}
