package middleware

import (
	"log/slog"
	"strings"
	"time"
	"zhp-app/pkg/common"

	"github.com/gin-gonic/gin"
)

// RequestLogger 在请求处理完成后输出一条结构化访问日志。
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// 未命中路由时 FullPath 为空，这里回退到原始请求路径。
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", c.Writer.Status()),
			slog.String("clientIP", c.ClientIP()),
			slog.Duration("latency", time.Since(start)),
			slog.Int("errorCount", len(c.Errors)),
			slog.Int64("responseSize", int64(c.Writer.Size())),
			slog.String("tenantCode", c.GetString(common.TenantCode)),
			slog.String("username", c.GetString(common.Username)),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", collectErrors(c)))
		}

		// 4xx 和 5xx 提升为 warn/error，方便快速定位异常请求。
		switch {
		case c.Writer.Status() >= 500:
			slog.Error("http_request", attrs...)
		case c.Writer.Status() >= 400:
			slog.Warn("http_request", attrs...)
		default:
			slog.Info("http_request", attrs...)
		}
	}
}

// collectErrors 把 Gin 内部错误列表压平成一段便于写日志的字符串。
func collectErrors(c *gin.Context) string {
	errors := make([]string, 0, len(c.Errors))
	for _, item := range c.Errors {
		errors = append(errors, item.Error())
	}

	return strings.Join(errors, "; ")
}
