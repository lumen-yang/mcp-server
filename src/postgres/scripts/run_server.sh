#!/usr/bin/env bash
# 本地启动 PostgreSQL MCP Server：自动加载 .env、编译并启动。

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${PROJECT_DIR}/.bin"
SERVER_BIN="${BIN_DIR}/postgres-server"

if [[ ! -f "${PROJECT_DIR}/.env" ]]; then
  echo "错误：未找到 ${PROJECT_DIR}/.env，请先执行 'cp .env.example .env' 并填写密钥。" >&2
  exit 1
fi

cd "${PROJECT_DIR}"
set -a
# shellcheck disable=SC1091
source .env
set +a

if [[ -z "${MCP_SECRET_ID:-${TENCENTCLOUD_SECRET_ID:-}}" || -z "${MCP_SECRET_KEY:-${TENCENTCLOUD_SECRET_KEY:-}}" ]]; then
  echo "错误：.env 中缺少 MCP_SECRET_ID/MCP_SECRET_KEY。" >&2
  echo "      兼容过渡期仍可读取 TENCENTCLOUD_SECRET_ID/TENCENTCLOUD_SECRET_KEY。" >&2
  exit 1
fi

if [[ -n "${MCP_API_TOKEN:-}" ]]; then
  echo "==> 已启用 MCP_API_TOKEN 入口鉴权"
fi

mkdir -p "${BIN_DIR}"
echo "==> 编译本地可执行文件..."
env -u GOOS -u GOARCH go build -o "${SERVER_BIN}" .

echo "==> 启动 PostgreSQL MCP Server..."
exec "${SERVER_BIN}"
