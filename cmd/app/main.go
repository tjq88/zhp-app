package main

import (
	"log/slog"
	"zhp-app/internal/app"
	"zhp-app/pkg/common"
	"zhp-app/pkg/config"
	"zhp-app/pkg/idgenx"
	"zhp-app/pkg/logger"
)

// main 是进程唯一入口，负责串联配置、基础设施、业务服务和 HTTP 路由。
func main() {
	// 先加载配置，再初始化日志，这样启动失败也能被记录下来。
	cfg, err := config.Load()
	if err != nil {
		logger.Init("info")
		slog.Error("load_config_failed", slog.String("err", err.Error()))
		return
	}

	// 配置可用后，按配置切换日志级别。
	logger.Init(cfg.LogLevel)
	slog.Info("config_loaded",
		slog.String("port", cfg.Port),
		slog.Uint64("configuredWorkerID", uint64(cfg.WorkerID)),
		slog.String("logLevel", cfg.LogLevel),
	)

	// Redis 同时用于业务访问和雪花算法 workerId 注册。
	redis, err := common.InitRedis(cfg.RedisConfig)
	if err != nil {
		slog.Error("init_redis_failed", slog.String("err", err.Error()))
		return
	}
	slog.Info("redis_initialized")
	defer redis.Close()

	// 通过 Redis 注册雪花算法 workerId，保证多实例部署时不会冲突。
	workerID, err := idgenx.InitFromRedis(cfg.RedisConfig)
	if err != nil {
		slog.Error("id_generator_init_failed", slog.String("err", err.Error()))
		return
	}
	defer idgenx.Unregister()

	slog.Info("id_generator_initialized",
		slog.Uint64("workerID", uint64(workerID)),
	)

	// 启动阶段主动校验数据库连通性，尽早失败，避免服务启动后再暴露问题。
	_, err = common.InitDB(cfg.MySQLDSN)
	if err != nil {
		slog.Error("init_db_failed", slog.String("err", err.Error()))
		return
	}
	slog.Info("db_initialized")

	// 基础设施就绪后再组装业务服务和 HTTP 路由。
	router := app.NewRouter()

	slog.Info("server_starting", slog.String("port", cfg.Port))
	if err := router.Run(cfg.Port); err != nil {
		slog.Error("server_start_failed", slog.String("port", cfg.Port), slog.String("err", err.Error()))
	}
}
