package app

import (
	"net/http"
	"zhp-app/internal/handler"
	"zhp-app/internal/middleware"
	"zhp-app/internal/service"

	"github.com/gin-gonic/gin"
)

// NewRouter 创建 Gin 引擎，并注册全局中间件和业务路由。
func NewRouter() *gin.Engine {
	// 服务
	memberService := service.NewMemberService()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(), middleware.JwtMiddleWare())
	memberHandler := handler.NewMemberHandler(memberService)

	// 健康检查保持简单，方便容器编排系统做存活探测。
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "ok"})
	})

	// 会员注册接口。
	router.POST("app-api/member/register", memberHandler.Register)
	router.POST("app-api/member/login", memberHandler.Login)

	return router
}
