# 云数据库 TencentDB for PostgreSQL

> 腾讯云数据库 PostgreSQL（TencentDB for PostgreSQL，云 API 使用 `postgres` 作为简称）能够让您在云端轻松设置、操作和扩展强大的开源数据库 PostgreSQL。

---

## Tools

当前默认注册 **48 个工具**，覆盖实例、账号、数据库、参数、备份、监控、网络、只读实例与 SSL 配置等能力。

---

## 当前版本定位

当前主分支已切换为 **Hosted URL + 请求头直传 `SecretId/SecretKey` + `streamable-http`** 的接入模型：

1. MCP 客户端直接连接 `/mcp`
2. 每次请求都带：
   - `X-TencentCloud-Secret-Id`
   - `X-TencentCloud-Secret-Key`
   - `X-TencentCloud-Session-Token`（可选，临时凭证时使用）
3. 服务端按请求解析凭据，并在默认配置下调用 `STS GetCallerIdentity` 做身份确认
4. 工具执行时直接使用本次请求携带的腾讯云凭证创建 SDK Client
5. 默认建议启用 `MCP_STREAMABLE_HTTP_STATELESS=true`，避免跨请求依赖进程内 session

> 旧的 **托管 URL + token/OAuth/Authorization + SSE** 版本已按原样保存在：`variants/token-oauth-authorization-version`

---

## 快速开始

### 1. 准备配置

```bash
cp .env.example .env
```

最小建议配置：

```env
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_STS_REGION=ap-guangzhou
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
READ_ONLY=true
```

然后启动：

```bash
./scripts/run_server.sh
```

### 2. 在 MCP 客户端中配置 Hosted URL

服务启动后，终端会输出一份可复制的配置。默认形态如下：

```json
{
 "mcpServers": {
  "mcp-server-postgres": {
   "type": "streamable-http",
   "url": "http://127.0.0.1:9000/mcp",
   "headers": {
    "X-TencentCloud-Secret-Id": "<TENCENTCLOUD_SECRET_ID>",
    "X-TencentCloud-Secret-Key": "<TENCENTCLOUD_SECRET_KEY>",
    "X-TencentCloud-Session-Token": "<TENCENTCLOUD_SESSION_TOKEN_OPTIONAL>"
   }
  }
 }
}
```

### 3. 直接调用示例

下面是一个最小 `initialize` 请求示例：

```bash
curl http://127.0.0.1:9000/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -H 'X-TencentCloud-Secret-Id: 你的SecretId' \
  -H 'X-TencentCloud-Secret-Key: 你的SecretKey' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","clientInfo":{"name":"demo-client","version":"1.0.0"},"capabilities":{}}}'
```

如果你使用的是 STS 临时凭证，再额外加上：

```bash
-H 'X-TencentCloud-Session-Token: 你的SessionToken'
```

---

## 本仓库自带客户端如何连当前服务

`cmd/mcp_smoke`、`cmd/verify`、`cmd/write_test` 已自动支持从环境变量读取请求头：

```bash
export MCP_REQUEST_SECRET_ID=你的SecretId
export MCP_REQUEST_SECRET_KEY=你的SecretKey
export MCP_REQUEST_SESSION_TOKEN=你的SessionToken # 可选
```

也兼容直接复用：

```bash
export MCP_SECRET_ID=你的SecretId
export MCP_SECRET_KEY=你的SecretKey
```

然后运行：

```bash
VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
```

---

## 当前鉴权与授权行为

### 数据面鉴权

默认 `MCP_AUTH_MODE=request-credential`：

- 服务端从请求头提取腾讯云凭证
- 默认执行 `STS GetCallerIdentity`
- 认证成功后构造 `Principal`
- `tools/registry.go` 仍沿用现有 `Principal + Guard` 授权链路

### Scope 与地域

请求凭据模式下，Principal 的默认权限由环境变量控制：

```env
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read,pg.write
MCP_REQUEST_ALLOWED_REGIONS=ap-guangzhou
```

说明：

- `pg.read`：允许只读工具
- `pg.write`：允许写类工具（仍继续受 `Guard` 和 `confirm=true` 约束）
- `MCP_REQUEST_ALLOWED_REGIONS` 留空表示不额外做地域限制

---

## 与旧 token 版的关系

当前仓库保留了兼容代码分支能力，但默认入口已不再推荐 `issued-token`。

如果你需要继续使用下面这类模型：

- `POST /auth/token-exchange/tencentcloud`
- `POST /auth/bootstrap/tencentcloud`
- `Authorization: Bearer <MCP_ACCESS_TOKEN>`
- `/admin/tokens`
- `type: sse`

请直接参考保留副本：`variants/token-oauth-authorization-version`

---

## 兼容模式

### `shared-token`

```env
MCP_AUTH_MODE=shared-token
MCP_API_TOKEN=请替换为高强度随机串
MCP_SECRET_ID=你的SecretId
MCP_SECRET_KEY=你的SecretKey
```

适用于受控环境下的共享入口，但所有客户端共用一个访问 token。

### `none`

```env
MCP_AUTH_MODE=none
MCP_SECRET_ID=你的SecretId
MCP_SECRET_KEY=你的SecretKey
```

仅建议本机临时调试，**不要**用于共享环境或远程部署。

### `issued-token`

当前主仓库仍保留兼容实现，但默认已不再推荐。若需要完整旧版使用说明，请查看：

- `variants/token-oauth-authorization-version/README.md`

---

## 安全部署建议

- **生产环境必须放在 HTTPS / 反向代理之后**
- **不要把 `SecretId` / `SecretKey` 放进 URL 或 query 参数**
- **禁止在日志、trace、错误回显中输出凭据明文**
- **优先使用 STS 临时凭证，而不是长期 AK/SK**
- **保持 `READ_ONLY=true` 起步**，确认流程后再按需开放写操作
- **对外暴露前务必加 IP 白名单、VPN 或零信任访问控制**
- **部署在 SCF / API 网关等无状态环境时，建议保持 `MCP_STREAMABLE_HTTP_STATELESS=true`**

---

## 回归验证

本地快速验活：

```bash
./scripts/run_server.sh
```

另开一个终端：

```bash
export MCP_REQUEST_SECRET_ID=你的SecretId
export MCP_REQUEST_SECRET_KEY=你的SecretKey
VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
```

如果要做更完整的协议验证，可继续使用：

- `cmd/mcp_smoke`
- `cmd/verify`
- `cmd/write_test`
