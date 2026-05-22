package config

import (
	"fmt"
	"os"
	"strconv"
)

type App struct {
	Port     string
	MySqlDsn string
	PwdKey   string
	WorkerID uint16
	LogLevel string
}

func Load() (App, error) {
	workerID, err := getEnvUint16("WORKER_ID", 0)
	if err != nil {
		return App{}, fmt.Errorf("invalid WORKER_ID: %w", err)
	}

	cfg := App{
		Port:     getEnv("APP_PORT", ":8080"),
		MySqlDsn: getEnv("MYSQL_DSN", "root:root123456@tcp(127.0.0.1:3306)/zpxc?charset=utf8mb4&parseTime=True&loc=Local"),
		PwdKey:   getEnv("PWD_KEY", "53b8e2d890c5535a574f8f19eea8ef4451ec0f43e8b0d5a0d616f1da9578d1b4"),
		WorkerID: workerID,
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	if cfg.PwdKey == "" {
		return App{}, fmt.Errorf("PWD_KEY cannot be empty")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvUint16(key string, fallback uint16) (uint16, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(parsed), nil
}
