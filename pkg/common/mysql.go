package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
)

// Db 保存当前进程共享使用的 GORM 连接。
var Db *gorm.DB

// InitDB 创建 GORM 连接，校验可用性，并保存为共享实例。
func InitDB(dsn string) (*gorm.DB, error) {
	slog.Info("db_connecting")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 启动阶段先 Ping，确保服务对外提供能力前数据库已可用。
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db failed: %w", err)
	}

	Db = db
	slog.Info("db_connected")
	return db, nil
}
