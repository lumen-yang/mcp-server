# PostgreSQL MCP Server 部署到腾讯云 SCF（Web 函数）

> **当前主分支默认模式是 `request-credential + streamable-http`**：客户端每次请求通过 Header 直传 `SecretId/SecretKey`，并连接单端点 `/mcp`。如果你需要旧的 `issued-token` / Bearer token / SSE 版本，请查看 `variants/token-oauth-authorization-version` 副本。

本文按**当前主分支**说明如何把服务部署到腾讯云 SCF，并让 MCP 客户端直接连接云函数 URL。

---

## 1. 方案结论

当前推荐的 SCF 形态是：

- **SCF 只托管 MCP Server**
- **客户端直接连接 `https://你的函数URL/mcp`**
- **客户端在每次请求 Header 中传自己的腾讯云凭证**
- **SCF 环境变量中不保存用户的 `SecretId/SecretKey`**
- **服务端启用 `MCP_STREAMABLE_HTTP_STATELESS=true`**，避免跨请求依赖进程内 session

这套模式适合你当前主分支，因为：

- 不依赖 SCF 内部的 token store
- 不需要在云函数里保存长期云密钥
- 和当前代码默认的 `request-credential` 鉴权模式一致
- 不再依赖 `SSE / message` 双请求落到同一实例
- 更适合 Web 函数、网关、无状态弹性实例这类运行环境

### 1.1 Header 约定

客户端在访问 `/mcp` 时，按请求附带：

- `X-TencentCloud-Secret-Id`
- `X-TencentCloud-Secret-Key`
- `X-TencentCloud-Session-Token`（可选，仅临时凭证时使用）

> **不要把密钥放到 URL / query 里。** 仅通过 `HTTPS` + Header 传递，并确保网关/日志不会记录这些 Header。

---

## 2. 已准备好的交付物

仓库里已经准备好以下 SCF 文件：

- `deploy/scf/scf_bootstrap`：SCF 包内启动文件
- `deploy/scf/scf.console.startup.sh`：控制台可粘贴的启动命令模板
- `deploy/scf/scf.env.example`：SCF 环境变量模板
- `deploy/scf/scf.console.env.txt`：控制台环境变量最小清单
- `scripts/build_scf_zip.sh`：构建 Linux zip 包脚本

默认打包输出：

```bash
./scripts/build_scf_zip.sh
```

生成：

```bash
dist/postgres-mcp-scf-web-linux-amd64.zip
```

如需 ARM：

```bash
./scripts/build_scf_zip.sh arm64
```

---

## 3. SCF 控制台创建函数

建议按下面方式创建：

- **函数类型**：Web 函数
- **运行环境**：Go 标准运行环境
- **代码上传方式**：本地上传 zip
- **架构**：与 zip 保持一致（默认 `amd64`）

上传 `dist/postgres-mcp-scf-web-linux-amd64.zip` 后，开启函数 URL 公网访问。

### 3.1 启动命令

zip 包已经内置 `scf_bootstrap`，通常直接使用包内启动文件即可。

如果控制台要求手动填写启动命令，请填与 `deploy/scf/scf.console.startup.sh` 相同的内容：

```bash
#!/bin/bash
set -euo pipefail

export PG_MCP_RUNTIME="${PG_MCP_RUNTIME:-scf}"
export PORT="${PORT:-9000}"
export MCP_SERVER_BIND_HOST="${MCP_SERVER_BIND_HOST:-0.0.0.0}"
export MCP_SERVER_PORT="${MCP_SERVER_PORT:-${PORT}}"
export MCP_SERVER_HTTP_ENDPOINT="${MCP_SERVER_HTTP_ENDPOINT:-/mcp}"
export MCP_STREAMABLE_HTTP_STATELESS="${MCP_STREAMABLE_HTTP_STATELESS:-true}"
export MCP_AUTH_MODE="${MCP_AUTH_MODE:-request-credential}"
export MCP_REQUEST_VALIDATE_IDENTITY="${MCP_REQUEST_VALIDATE_IDENTITY:-true}"
export MCP_REQUEST_CREDENTIAL_SCOPES="${MCP_REQUEST_CREDENTIAL_SCOPES:-pg.read}"
export MCP_REQUEST_ALLOWED_REGIONS="${MCP_REQUEST_ALLOWED_REGIONS:-}"
export MCP_STS_REGION="${MCP_STS_REGION:-ap-guangzhou}"
export READ_ONLY="${READ_ONLY:-true}"
export TOKEN_EXCHANGE_ENABLED="${TOKEN_EXCHANGE_ENABLED:-false}"

exec /var/user/postgres-server
```

---

## 4. SCF 环境变量怎么填

建议直接参考 `deploy/scf/scf.console.env.txt`，最小可用配置如下：

```env
MCP_AUTH_MODE=request-credential
TOKEN_EXCHANGE_ENABLED=false
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
READ_ONLY=true
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
```

推荐补充：

```env
MCP_SERVER_PUBLIC_URL=https://你的函数URL/mcp
MCP_STS_REGION=ap-guangzhou
MCP_REQUEST_ALLOWED_REGIONS=ap-guangzhou
FEATURES=instance,account,database,parameter,backup,monitoring,network,readonly
```

### 4.1 这些变量不要放进 SCF

以下变量**不要**配置到云函数环境变量中：

- `MCP_SECRET_ID`
- `MCP_SECRET_KEY`
- `MCP_REQUEST_SECRET_ID`
- `MCP_REQUEST_SECRET_KEY`
- `MCP_API_TOKEN`
- `MCP_ACCESS_TOKEN`

原因很简单：

- 这些都属于**用户侧或客户端侧凭据**
- 当前模式要求**按请求传递**，而不是提前固化在服务端环境中
- 放进 SCF 环境变量会扩大泄露面，不符合这个模式的目标

---

## 5. 客户端如何连接云函数

部署成功后，MCP 客户端直接连：

```text
https://你的函数URL/mcp
```

并在 Header 中携带自己的腾讯云凭证。

### 5.1 MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "https://你的函数URL/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "你的SecretId",
        "X-TencentCloud-Secret-Key": "你的SecretKey"
      }
    }
  }
}
```

如果你用的是临时凭证，再额外带上：

```json
{
  "X-TencentCloud-Session-Token": "你的SessionToken"
}
```

> **更推荐用临时凭证而不是长期 AK/SK。**

---

## 6. 部署后可访问的地址

函数 URL 开通公网访问后，主要地址如下：

- **健康检查**：`https://你的函数URL/healthz`
- **就绪检查**：`https://你的函数URL/readyz`
- **MCP streamable-http**：`https://你的函数URL/mcp`

当前主分支默认不依赖 `/auth/token-exchange/tencentcloud` 作为主链路。

---

## 7. 本地联调 / 远程验收

### 7.1 健康检查

先验证函数活着：

```text
https://你的函数URL/healthz
```

预期返回 `200 OK`。

### 7.2 MCP 协议冒烟

仓库里的冒烟工具会优先读取 `MCP_REQUEST_SECRET_ID` / `MCP_REQUEST_SECRET_KEY` 并自动带 Header：

```bash
MCP_REQUEST_SECRET_ID=你的SecretId \
MCP_REQUEST_SECRET_KEY=你的SecretKey \
go run ./cmd/mcp_smoke --url https://你的函数URL/mcp --region ap-guangzhou
```

### 7.3 真实只读能力验证

```bash
MCP_REQUEST_SECRET_ID=你的SecretId \
MCP_REQUEST_SECRET_KEY=你的SecretKey \
go run ./cmd/verify --url https://你的函数URL/mcp --region ap-guangzhou --instance-id postgres-xxxxxxxx
```

如果你使用临时凭证，再补：

```bash
MCP_REQUEST_SESSION_TOKEN=你的SessionToken
```

---

## 8. 安全建议

- **只用 `HTTPS` 暴露函数 URL**
- **不要在网关、CDN、日志平台记录上述鉴权 Header**
- **不要把密钥拼到 URL、query、日志、报错回显中**
- **优先用临时凭证**，不要长期复用主账号密钥
- **先以 `pg.read` + `READ_ONLY=true` 起步**，确认链路没问题后再放开能力
- **无状态环境建议保持 `MCP_STREAMABLE_HTTP_STATELESS=true`**
- **如果要做公网多租户生产化**，建议后续再引入更短期 token 或外部鉴权网关
