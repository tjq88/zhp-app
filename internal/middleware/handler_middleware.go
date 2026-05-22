package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"zhp-app/pkg/common"
)

// JwtMiddleWare 从请求头提取轻量身份信息。
// 当前它还是占位实现，暂未进行真实 JWT 校验。
func JwtMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("authorization")
		if authorization != "" {
			// 当前先写入一个模拟 userId，便于下游示例代码联调。
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
