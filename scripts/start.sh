#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# 使用项目内缓存目录，避免写入系统级 Go 缓存目录。
export GOCACHE="${GOCACHE:-$ROOT_DIR/.gocache}"

# HTTP
export APP_PORT="${APP_PORT:-:8080}"
export LOG_LEVEL="${LOG_LEVEL:-info}"

# MySQL
export MYSQL_DSN="${MYSQL_DSN:-root:root123456@tcp(127.0.0.1:3306)/zpxc?charset=utf8mb4&parseTime=True&loc=Local}"

# Redis
export REDIS_ADDR="${REDIS_ADDR:-127.0.0.1:6379}"
export REDIS_PASSWORD="${REDIS_PASSWORD:-}"
export REDIS_DB="${REDIS_DB:-0}"
export REDIS_MASTER_NAME="${REDIS_MASTER_NAME:-}"
export REDIS_WORKER_ID_MIN="${REDIS_WORKER_ID_MIN:-0}"
export REDIS_WORKER_ID_MAX="${REDIS_WORKER_ID_MAX:-1023}"
export REDIS_WORKER_ID_LIFE_SECONDS="${REDIS_WORKER_ID_LIFE_SECONDS:-15}"

# 安全
export PWD_KEY="${PWD_KEY:-53b8e2d890c5535a574f8f19eea8ef4451ec0f43e8b0d5a0d616f1da9578d1b4}"

echo "启动 zhp-app..."
echo "APP_PORT=$APP_PORT"
echo "MYSQL_DSN=$MYSQL_DSN"
echo "REDIS_ADDR=$REDIS_ADDR"
echo "REDIS_DB=$REDIS_DB"
echo "LOG_LEVEL=$LOG_LEVEL"

go run ./cmd/app
