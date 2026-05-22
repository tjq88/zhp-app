package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init 初始化进程级结构化日志。
func Init(level string) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	})

	slog.SetDefault(slog.New(handler))
}

// parseLevel 把环境变量中的字符串日志级别转换成 slog 级别。
func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
