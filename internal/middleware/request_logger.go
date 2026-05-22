package middleware

import (
	"log/slog"
	"strings"
	"time"
	"zhp-app/pkg/common"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

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

func collectErrors(c *gin.Context) string {
	errors := make([]string, 0, len(c.Errors))
	for _, item := range c.Errors {
		errors = append(errors, item.Error())
	}

	return strings.Join(errors, "; ")
}
