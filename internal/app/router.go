package app

import (
	"net/http"
	"zhp-app/internal/handler"
	"zhp-app/internal/middleware"
	"zhp-app/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(memberService *service.MemberService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(), middleware.JwtMiddleWare())
	memberHandler := handler.NewMemberHandler(memberService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "ok"})
	})

	//注册
	router.POST("app-api/member/register", memberHandler.Register)

	return router
}
