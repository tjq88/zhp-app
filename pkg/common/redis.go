package common

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"zhp-app/pkg/config"
)

// redisClient 保存当前进程共享使用的 Redis 客户端。
var RedisClient *redis.Client

// InitRedis 初始化 Redis 客户端并校验连通性。
func InitRedis(conf config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("ping redis failed: %w", err)
	}
	slog.Info("redis_connected")

	RedisClient = client

	return client, nil
}
