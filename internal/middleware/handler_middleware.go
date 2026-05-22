package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"zhp-app/pkg/common"
)

func JwtMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("authorization")
		if authorization != "" {
			c.Set(common.Username, authorization)
			c.Set(common.UserId, (int64(100)))
		}

		tenantCode := c.Request.Header.Get("tenantCode")
		if tenantCode != "" {
			c.Set(common.TenantCode, tenantCode)
		}

		if authorization != "" || tenantCode != "" {
			slog.Debug("jwt_context_loaded",
				slog.Bool("hasAuthorization", authorization != ""),
				slog.String("tenantCode", tenantCode),
			)
		}
		c.Next()
	}
}
