#!/usr/bin/env bash
# PG MCP 工具回归验证脚本
# 用法：在 src/postgres 目录下执行 VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
#
# 该脚本会：
#   1. 编译 MCP server 和验证客户端(cmd/verify)到临时目录
#   2. 加载 .env 启动本地 server 进程（127.0.0.1:9000）
#   3. 运行验证客户端，对一批只读(Describe*)接口发起真实调用
#   4. 打印结果，并在脚本退出时自动清理 server 进程
#
# 当前默认的 request-credential 模式下：
#   - 服务端本身不要求预置 MCP_SECRET_ID / MCP_SECRET_KEY
#   - 验证客户端会优先读取 MCP_REQUEST_SECRET_ID / MCP_REQUEST_SECRET_KEY
#   - 若未设置，也会回退读取 MCP_SECRET_ID / MCP_SECRET_KEY
#
# 可选环境变量：
#   VERIFY_REGION=ap-guangzhou
#   VERIFY_INSTANCE_ID=postgres-xxxxxxxx
#   VERIFY_SERVER_PORT=9000
#
# 注意：验证脚本仅调用只读接口，不会对云上资源产生写操作；
#      调用范围受 .env 中 GUARD_PROFILE / SCOPE 配置限制。

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="$(mktemp -d)"
SERVER_BIN="${BIN_DIR}/postgres_server"
VERIFY_BIN="${BIN_DIR}/pg_verify"
LOG_FILE="${BIN_DIR}/server.log"
SERVER_PID=""

cleanup() {
  if [[ -n "${SERVER_PID}" ]] && kill -0 "${SERVER_PID}" 2>/dev/null; then
    echo ""
    echo "==> 停止本地 server 进程 (PID ${SERVER_PID})"
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
  rm -rf "${BIN_DIR}"
}
trap cleanup EXIT

pick_free_port() {
  python3 - <<'PY'
import socket
s = socket.socket()
s.bind(("127.0.0.1", 0))
print(s.getsockname()[1])
s.close()
PY
}

if [[ ! -f "${PROJECT_DIR}/.env" ]]; then
  echo "错误：未找到 ${PROJECT_DIR}/.env，请先根据 .env.example 配置真实密钥后再运行。" >&2
  exit 1
fi

cd "${PROJECT_DIR}"

# 探测本机 OS/ARCH，显式覆盖编译目标，避免继承 shell 里为交叉编译（如 Docker/Linux 部署）
# 导出的 GOOS/GOARCH（例如 linux/amd64），导致产物在本机（如 macOS/arm64）无法执行。
case "$(uname -s)" in
  Darwin) HOST_GOOS="darwin" ;;
  Linux)  HOST_GOOS="linux" ;;
  *) HOST_GOOS="$(uname -s | tr '[:upper:]' '[:lower:]')" ;;
esac
case "$(uname -m)" in
  arm64|aarch64) HOST_GOARCH="arm64" ;;
  x86_64|amd64)  HOST_GOARCH="amd64" ;;
  *) HOST_GOARCH="$(uname -m)" ;;
esac

echo "==> 编译 server 与验证客户端（目标平台: ${HOST_GOOS}/${HOST_GOARCH}）..."
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${SERVER_BIN}" .
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${VERIFY_BIN}" ./cmd/verify

echo "==> 加载 .env 并启动本地 server..."
set -a
# shellcheck disable=SC1091
source .env
set +a

if [[ -z "${VERIFY_SERVER_PORT:-}" ]]; then
  VERIFY_SERVER_PORT="$(pick_free_port)"
else
  VERIFY_SERVER_PORT="${VERIFY_SERVER_PORT}"
fi
VERIFY_HTTP_ENDPOINT="${MCP_SERVER_HTTP_ENDPOINT:-${MCP_SERVER_SSE_ENDPOINT:-/mcp}}"
VERIFY_SERVER_URL="${VERIFY_SERVER_URL:-${VERIFY_SSE_URL:-http://127.0.0.1:${VERIFY_SERVER_PORT}${VERIFY_HTTP_ENDPOINT}}}"
VERIFY_REGION="${VERIFY_REGION:-ap-guangzhou}"

if [[ -z "${VERIFY_INSTANCE_ID:-}" ]]; then
  echo "错误：请设置 VERIFY_INSTANCE_ID=postgres-xxxxxxxx 后再运行。" >&2
  exit 1
fi

MCP_SERVER_PORT="${VERIFY_SERVER_PORT}" "${SERVER_BIN}" > "${LOG_FILE}" 2>&1 &
SERVER_PID=$!

# 等待 server 就绪
READY=0
for _ in $(seq 1 40); do
  if ! kill -0 "${SERVER_PID}" 2>/dev/null; then
    break
  fi
  if curl -fsS "http://127.0.0.1:${VERIFY_SERVER_PORT}/healthz" >/dev/null 2>&1; then
    READY=1
    break
  fi
  sleep 0.5
done

echo "----- server 启动日志 -----"
cat "${LOG_FILE}"
echo "---------------------------"

if [[ "${READY}" != "1" ]] || ! kill -0 "${SERVER_PID}" 2>/dev/null; then
  echo "错误：本地 server 未成功启动，请检查上方日志。" >&2
  exit 1
fi

echo "==> 运行验证客户端（对真实腾讯云测试实例发起只读调用）..."
VERIFY_REGION="${VERIFY_REGION}" \
VERIFY_INSTANCE_ID="${VERIFY_INSTANCE_ID}" \
VERIFY_SERVER_URL="${VERIFY_SERVER_URL}" \
"${VERIFY_BIN}"

echo ""
echo "==> 验证完成。"
