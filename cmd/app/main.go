package main

import (
	"log/slog"
	"zhp-app/internal/app"
	"zhp-app/internal/service"
	"zhp-app/pkg/common"
	"zhp-app/pkg/config"
	"zhp-app/pkg/idgenx"
	"zhp-app/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Init("info")
		slog.Error("load_config_failed", slog.String("err", err.Error()))
		return
	}

	logger.Init(cfg.LogLevel)
	slog.Info("config_loaded",
		slog.String("port", cfg.Port),
		slog.Uint64("workerID", uint64(cfg.WorkerID)),
		slog.String("logLevel", cfg.LogLevel),
	)

	idgenx.Init(cfg.WorkerID)
	slog.Info("id_generator_initialized",
		slog.Uint64("workerID", uint64(cfg.WorkerID)),
	)

	db, err := common.InitDB(cfg.MySqlDsn)
	if err != nil {
		slog.Error("init_db_failed", slog.String("err", err.Error()))
		return
	}
	slog.Info("db_initialized")

	memberService := service.NewMemberService(db, cfg.PwdKey)
	router := app.NewRouter(memberService)
	slog.Info("server_starting", slog.String("port", cfg.Port))
	if err := router.Run(cfg.Port); err != nil {
		slog.Error("server_start_failed", slog.String("port", cfg.Port), slog.String("err", err.Error()))
	}
}
