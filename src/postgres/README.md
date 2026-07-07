# 云数据库 TencentDB for PostgreSQL

> 腾讯云数据库 PostgreSQL（TencentDB for PostgreSQL，云 API 使用 `postgres` 作为简称）能够让您在云端轻松设置、操作和扩展强大的开源数据库 PostgreSQL。腾讯云负责软件安装、存储管理、高可用复制与备份恢复等基础设施能力，让您更专注于业务开发。

---

## Tools

当前默认注册 **48 个工具**，覆盖以下能力：

### 1. 实例管理
- `DescribeDBInstances`：查询实例列表
- `DescribeDBInstanceAttribute`：查询实例详情
- `DescribeClasses`：查询售卖规格
- `DescribeDBVersions`：查询数据库版本
- `DescribeTasks`：查询实例任务
- `DescribeRegions` / `DescribeZones` / `DescribeProductConfig`
- `UpgradeDBInstanceKernelVersion` / `CreateInstances` / `ModifyDBInstanceName` / `ModifyDBInstanceSpec` / `RestartDBInstance` / `IsolateDBInstances` / `DisIsolateDBInstances`

### 2. 账号与数据库
- `DescribeAccounts` / `DescribeAccountPrivileges`
- `CreateAccount` / `DeleteAccount` / `ModifyAccountPrivileges` / `ResetAccountPassword`
- `DescribeDatabases` / `DescribeDatabaseObjects`
- `CreateDatabase` / `ModifyDatabaseOwner`

### 3. 参数、备份、监控、网络与只读实例
- `DescribeDBInstanceParameters` / `DescribeParameterTemplates` / `DescribeParameterTemplateAttributes` / `DescribeParamsEvent` / `ModifyDBInstanceParameters`
- `DescribeBackupOverview` / `DescribeBaseBackups` / `DescribeLogBackups` / `DescribeAvailableRecoveryTime` / `DescribeCloneDBInstanceSpec` / `DescribeBackupDownloadURL` / `CreateBaseBackup` / `CloneDBInstance`
- `DescribeSlowQueryList` / `DescribeSlowQueryAnalysis` / `DescribeDBErrlogs`
- `DescribeDBInstanceSecurityGroups` / `ModifyDBInstanceSecurityGroups` / `OpenDBExtranetAccess` / `CloseDBExtranetAccess`
- `DescribeReadOnlyGroups` / `CreateReadOnlyDBInstance`
- `DescribeDBInstanceSSLConfig`

---

## 快速开始

### 方式一：本地直接运行

1. 复制配置模板
   ```bash
   cp .env.example .env
   ```
2. 填写 `.env` 中的 **`MCP_SECRET_ID`** / **`MCP_SECRET_KEY`**
3. 如需团队共享或远程访问，建议同时设置 **`MCP_API_TOKEN`**
4. 启动服务
   ```bash
   ./scripts/run_server.sh
   ```

### 方式二：Docker Compose

```bash
cp .env.example .env
```

```bash
docker compose up -d --build
```

### MCP 客户端配置示例

默认本机运行时，可将以下配置填入支持 MCP 的客户端：

```json
{
 "mcpServers": {
  "mcp-server-postgres": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

如果你启用了 **`MCP_API_TOKEN`**，则客户端还需要携带鉴权头：

```json
{
 "mcpServers": {
  "mcp-server-postgres": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse",
   "headers": {
    "Authorization": "Bearer <MCP_API_TOKEN>"
   }
  }
 }
}
```

如果你部署在远程主机或反向代理后，请将地址替换为 `MCP_SERVER_PUBLIC_URL` 对应的实际可访问地址。

---

## 安全部署建议

建议初次部署时使用以下保守配置：

- **`READ_ONLY=true`**：先只开放只读能力，确认流程后再按需放开写操作
- **`MCP_SERVER_BIND_HOST=127.0.0.1`**：默认仅允许本机访问
- **`MCP_API_TOKEN=<高强度随机串>`**：远程或共享部署时开启入口鉴权
- **`SCOPE_ENABLED=true` + `REGION_SCOPE=ap-xxx`**：如果只管理单一地域，建议启用地域范围限制
- **不要直接暴露到公网**：如需远程访问，请放在反向代理、访问控制、VPN 或零信任网络后面

当前 `Guard` 支持通过 `GUARD_PROFILE` 切换不同安全画像；例如同一份 `.env` 中可并存 `GUARD_DEV_*` 与 `GUARD_TEST_*` 两套限制策略，切换时只需修改一个变量并重启服务。

---

## 鉴权与凭证说明

- **云 API 凭证**：服务端启动时优先读取 `MCP_SECRET_ID` / `MCP_SECRET_KEY`
- **旧变量兼容**：仍兼容 `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY`，但新部署不再推荐
- **MCP 服务入口鉴权**：设置 `MCP_API_TOKEN` 后，`/sse` 和 `/message` 都会校验 `Authorization: Bearer <token>`，同时兼容 `X-MCP-API-Token`
- **本地验证脚本**：`cmd/mcp_smoke`、`cmd/verify`、`cmd/write_test` 会自动读取 `MCP_API_TOKEN` 并附带请求头

---

## 常用验证命令

### 1. OpenAPI 参数对齐校验

```bash
./scripts/run_openapi_param_check.sh
```

### 2. MCP 协议联调

```bash
./scripts/run_mcp_smoke.sh
```

### 3. 真实只读接口验证

```bash
VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
```

---

## 部署与验证文档

- `DEPLOY.md`：本地 / Docker / 远程部署说明
- `scripts/TEST_REPORT.md`：公开仓库可用的验证记录模板
- `scripts/full_test_plan.yaml`：写操作验证配置模板
- `scripts/full_test_plan.observe.yaml`：分步观察式写验证模板

---

## 参考链接

- 产品详情：`https://cloud.tencent.com/product/postgres`
- API 文档：`https://cloud.tencent.com/document/product/409/16761`
- 地域列表：`https://cloud.tencent.com/document/product/1596/77930`

---

## 许可证

MIT
