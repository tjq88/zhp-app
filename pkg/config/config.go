package config

import (
	"fmt"
	"os"
	"strconv"
)

// AppConf 保存当前进程已加载的运行时配置，供共享访问。
var AppConf App

// App 汇总应用进程运行所需的配置项。
type App struct {
	Port      string
	MySqlDsn  string
	PwdKey    string
	WorkerID  uint16
	LogLevel  string
	RedisConf RedisConf
}

// RedisConf 同时包含 Redis 连接信息和 workerId 注册区间配置。
type RedisConf struct {
	Addr                string
	Password            string
	DB                  int
	MasterName          string
	WorkerIDMin         int32
	WorkerIDMax         int32
	WorkerIDLifeSeconds int32
}

// Load 从环境变量读取配置，并填充适合开发环境的默认值。
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

	redisDB, _ := getEnvInt("REDIS_DB", 0)
	//workerIDMin, _ := getEnvInt32("REDIS_WORKER_ID_MIN", 0)
	//workerIDMax, _ := getEnvInt32("REDIS_WORKER_ID_MAX", 1023)
	//workerIDLifeSeconds, _ := getEnvInt32("REDIS_WORKER_ID_LIFE_SECONDS", 15)
	cfg.RedisConf = RedisConf{
		Addr:                getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Password:            getEnv("REDIS_PASSWORD", ""),
		DB:                  redisDB,
		MasterName:          getEnv("REDIS_MASTER_NAME", ""),
		WorkerIDMin:         0,
		WorkerIDMax:         1023,
		WorkerIDLifeSeconds: 15,
	}

	// PwdKey 是密码加密必需配置，不能为空。
	if cfg.PwdKey == "" {
		return App{}, fmt.Errorf("PWD_KEY cannot be empty")
	}

	AppConf = cfg

	return cfg, nil
}

// getEnv 读取字符串环境变量，未设置时返回默认值。
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvUint16 解析无符号整型配置。
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

// getEnvInt 解析整型配置。
func getEnvInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return 0, err
	}

	return int(parsed), nil
}

// getEnvInt32 解析 32 位整型配置。
func getEnvInt32(key string, fallback int32) (int32, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(parsed), nil
}
