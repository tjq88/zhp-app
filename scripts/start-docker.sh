#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "启动 Docker 部署..."
docker compose -f deployments/docker-compose.yml up --build -d

echo "容器已启动，可使用以下命令查看状态："
echo "docker compose -f deployments/docker-compose.yml ps"
