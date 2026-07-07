# PostgreSQL MCP Server 部署指南

本文面向希望**自行部署并接入 MCP 客户端**的用户，覆盖本地运行、Docker 部署、环境变量说明与安全建议。

---

## 一、部署前准备

### 1. 需要准备的内容

- **腾讯云账号 API 密钥**：`MCP_SECRET_ID` / `MCP_SECRET_KEY`
- **运行环境**（二选一）
  - 本地运行：Go 1.24+
  - 容器运行：Docker 或 Docker Compose
- **MCP 客户端**：支持 `SSE` 类型 MCP Server 的客户端

### 2. 建议的最小安全配置

首次部署建议使用以下策略：

- **`READ_ONLY=true`**：先只开放只读能力
- **`MCP_SERVER_BIND_HOST=127.0.0.1`**：默认仅允许本机访问
- **`MCP_API_TOKEN=<高强度随机串>`**：远程或多人共享时开启入口鉴权
- **`REGION_SCOPE=ap-xxx`**：如果只管理单一地域，建议进一步限制地域范围
- **不要直接把 MCP 端口暴露到公网**：如需远程访问，请放在反向代理、访问控制或 VPN 后面

---

## 二、方式一：本地直接运行

### 1. 复制配置模板

```bash
cp .env.example .env
```

### 2. 编辑 `.env`

至少填好以下两项：

```env
MCP_SECRET_ID=你的SecretId
MCP_SECRET_KEY=你的SecretKey
```

如果需要共享访问或远程接入，建议再配置：

```env
MCP_API_TOKEN=请替换为高强度随机串
```

推荐保留以下默认值：

```env
MCP_SERVER_BIND_HOST=127.0.0.1
MCP_SERVER_PORT=9000
READ_ONLY=true
SCOPE_ENABLED=false
```

### 3. 启动服务

```bash
./scripts/run_server.sh
```

启动成功后，服务会在终端打印 MCP 客户端配置；如果启用了 `MCP_API_TOKEN`，打印结果中也会提示需要携带 `Authorization` 头。

---

## 三、方式二：使用 Docker Compose

### 1. 准备环境变量

```bash
cp .env.example .env
```

填写密钥后，直接启动：

```bash
docker compose up -d --build
```

默认配置下：

- 容器内服务监听 `0.0.0.0`
- 宿主机仅在 `127.0.0.1:${MCP_SERVER_PORT:-9000}` 暴露端口
- 外部网络无法直接访问该端口，更适合桌面客户端本机接入

### 2. 查看日志

```bash
docker compose logs -f
```

### 3. 停止服务

```bash
docker compose down
```

---

## 四、方式三：只用 Docker 命令

```bash
docker build -t mcp-server-postgres:latest .
```

```bash
docker run --rm -it \
  --name postgres-mcp \
  --env-file .env \
  -e MCP_SERVER_BIND_HOST=0.0.0.0 \
  -p 127.0.0.1:9000:9000 \
  mcp-server-postgres:latest
```

---

## 五、远程主机部署说明

如果你要把它部署到云主机、内网服务器，或通过域名提供给团队使用，建议：

1. 服务端监听：

```env
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_API_TOKEN=请替换为高强度随机串
```

2. 通过反向代理暴露，例如 `https://mcp.example.com/postgres/sse`

3. 将 `.env` 中的公开地址改为客户端实际访问地址：

```env
MCP_SERVER_PUBLIC_URL=https://mcp.example.com/postgres/sse
```

4. 在客户端配置 `Authorization: Bearer <MCP_API_TOKEN>`；若客户端不支持 Bearer 配置，也兼容 `X-MCP-API-Token`

> `MCP_SERVER_PUBLIC_URL` 只影响终端里打印给用户的 MCP 配置，不影响服务实际监听地址。

---

## 六、环境变量说明

| 变量 | 是否必填 | 说明 |
|---|---|---|
| `MCP_SECRET_ID` | 是 | 腾讯云 API SecretId |
| `MCP_SECRET_KEY` | 是 | 腾讯云 API SecretKey |
| `TENCENTCLOUD_SECRET_ID` | 否 | 旧变量名，当前仅做历史兼容 |
| `TENCENTCLOUD_SECRET_KEY` | 否 | 旧变量名，当前仅做历史兼容 |
| `MCP_API_TOKEN` | 否 | MCP 服务入口鉴权 token；设置后会校验 `Authorization: Bearer` 或 `X-MCP-API-Token` |
| `MCP_SERVER_BIND_HOST` | 否 | 服务监听地址，默认 `127.0.0.1` |
| `MCP_SERVER_PORT` | 否 | 服务端口，默认 `9000` |
| `MCP_SERVER_SSE_PORT` | 否 | 旧端口变量名，仍兼容 |
| `MCP_SERVER_SSE_ENDPOINT` | 否 | SSE 路径，默认 `/sse` |
| `MCP_SERVER_MESSAGE_ENDPOINT` | 否 | Message 路径，默认 `/message` |
| `MCP_SERVER_PUBLIC_URL` | 否 | 打印给客户端使用的完整 SSE URL |
| `FEATURES` | 否 | 功能组白名单，默认使用全部默认工具 |
| `GUARD_PROFILE` | 否 | 当前启用的 Guard 画像名 |
| `GUARD_<PROFILE>_*` | 否 | 指定画像下的 Guard 变量 |
| `READ_ONLY` | 否 | 全局只读保护，默认建议 `true` |
| `SCOPE_ENABLED` | 否 | 是否启用地域范围限制 |
| `REGION_SCOPE` | 否 | 限定允许操作的地域 |

---

## 七、推荐部署策略

### 个人本机使用

```env
MCP_SERVER_BIND_HOST=127.0.0.1
READ_ONLY=true
SCOPE_ENABLED=false
```

### 团队共享测试环境

```env
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_API_TOKEN=请替换为高强度随机串
READ_ONLY=true
SCOPE_ENABLED=true
REGION_SCOPE=ap-xxx
MCP_SERVER_PUBLIC_URL=https://mcp.example.com/postgres/sse
```

### 需要开放写操作时

只有在你非常清楚影响范围时，才建议关闭只读：

```env
READ_ONLY=false
```

即使关闭 `READ_ONLY`，涉及高风险或有成本的工具仍然需要显式传入 `confirm=true`。

---

## 八、部署后验证

### 1. 仅验证 MCP 协议连通性

```bash
./scripts/run_mcp_smoke.sh
```

### 2. 验证参数是否与 OpenAPI 对齐

```bash
./scripts/run_openapi_param_check.sh
```

### 3. 验证真实只读云 API 调用

```bash
VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
```

> `cmd/mcp_smoke`、`cmd/verify`、`cmd/write_test` 会自动读取环境变量 `MCP_API_TOKEN`，因此使用脚本启动时无需手工重复传 header。

---

## 九、相关文档

- 快速入口：`README.md`
- 本地启动脚本：`scripts/run_server.sh`
- 公开验证记录模板：`scripts/TEST_REPORT.md`
- 腾讯云 PostgreSQL API 文档：`https://cloud.tencent.com/document/product/409/16761`
