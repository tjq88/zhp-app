package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
)

var Db *gorm.DB

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

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db failed: %w", err)
	}

	Db = db
	slog.Info("db_connected")
	return db, nil
}
