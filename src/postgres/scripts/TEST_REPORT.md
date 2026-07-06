# PG MCP 工具集 — 综合测试报告

| 项目 | 内容 |
|---|---|
| **报告时间** | 2026-07-06 17:52 CST |
| **Server 版本** | `mcp-server-postgres` v1.0.0 |
| **测试实例** | `postgres-1lbqykq6`（`ap-chengdu`，`pg.it.medium4` / 2C4G / 130GB / PG 18.4） |
| **Guard 画像** | `test`（`READ_ONLY=false`，`ScopeEnabled=true`；当前代码仅保留地域范围限制，实例范围限制已移除） |
| **测试方式** | 本地编译 MCP server，分别执行参数对齐校验、真实云 API 只读调用验证、MCP 协议级 smoke test |
| **执行命令** | `go build ./...` / `./scripts/run_openapi_param_check.sh` / `./scripts/run_verify.sh` / `./scripts/run_mcp_smoke.sh` |

---

## 一、总体结论

| 维度 | 结果 |
|---|---|
| **工程构建** | ✅ `go build ./...` 通过 |
| **OpenAPI 参数对齐** | ✅ **48/48** 通过（`passed=48 failed=0 total=48`） |
| **真实云 API 只读调用** | ✅ **27/27** 通过 |
| **MCP 协议联调** | ✅ `initialize / ping / tools/list / tools/call` 全链路通过 |
| **客户端工具发现** | ✅ 成功识别 **48 个工具** |
| **写工具保护机制** | ✅ 写工具可被识别，且 `confirm=false` 时返回 `require_confirm=true`，未触发真实写操作 |

**结论**：当前 `postgres` MCP server 在 **参数契约、真实只读接口调用、MCP 客户端识别与协议交互** 三个层面均已通过验证；本次删除“实例访问范围限制”后，服务行为与当前设计保持一致。

---

## 二、本次重点回归项

### 1. 实例访问范围限制已移除

本次回归重点验证：此前为了安全测试加入的“实例访问范围限制”代码已删除，服务端不再自动注入 `DBInstanceId`，也不再对查询结果做实例级二次裁剪。

**验证结果：**
- `OpenAPI` 参数对齐全量通过，说明删除该限制后，工具参数结构未偏离腾讯云 SDK / OpenAPI。
- `MCP smoke test` 中，`postgres-DescribeDBInstances` 返回 `TotalCount=3`，不再像旧行为那样被强制收敛到单实例。
- `DescribeBaseBackups`、`DescribeLogBackups`、`DescribeReadOnlyGroups` 等相关查询路径在代码层已移除实例 scope 裁剪逻辑，且构建与联调均通过。

### 2. MCP 客户端可识别并可调用工具

这次新增了 `./scripts/run_mcp_smoke.sh` 与 `cmd/mcp_smoke`，用于从客户端视角验证：
- 能否连接本地 SSE MCP server
- 能否完成 `initialize`
- 能否列出工具
- 能否实际调用只读工具
- 写工具是否正确暴露 `confirm` 参数与保护逻辑

**验证结果：通过。**

---

## 三、OpenAPI 参数对齐结果

执行 `./scripts/run_openapi_param_check.sh` 后结果如下：

| 指标 | 数值 |
|---|---|
| 校验工具总数 | **48** |
| 通过 | **48** |
| 失败 | **0** |
| 通过率 | **100%** |

**说明：**
- 该校验只验证参数与腾讯云 `OpenAPI Request` / SDK 请求结构是否一致。
- 不会触发真实实例写操作。
- 可用于快速发现字段名、类型、必填项、归一化逻辑与 SDK 结构体不一致的问题。

---

## 四、真实云 API 调用验证结果

执行 `./scripts/run_verify.sh`，对一批只读 `Describe*` 接口发起真实云 API 调用。

| 指标 | 数值 |
|---|---|
| 注册工具总数 | **48** |
| 本次验证覆盖 | **27**（仅限只读 `Describe*` 接口） |
| 未覆盖 | **21**（写操作/管理类接口，需 `confirm=true`，本次有意跳过） |
| 通过 | **27** |
| 失败（500/403） | **0** |
| 通过率 | **100%（27/27）** |

### 实例组（8/15 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeDBInstanceAttribute` | ✅ | 返回实例详情（running / PG 18.4 / 2C4G / 130GB） |
| `DescribeDBInstances` | ✅ | 返回云侧真实实例列表；当前已不再受实例 scope 限制 |
| `DescribeClasses` | ✅ | 返回 `ap-chengdu-1` + PG 18 的规格列表 |
| `DescribeDBVersions` | ✅ | 返回可用 PG 版本列表 |
| `DescribeTasks` | ✅ | 返回异步任务列表（当前为空） |
| `DescribeRegions` | ✅ | 返回售卖地域列表 |
| `DescribeZones` | ✅ | 返回售卖可用区列表 |
| `DescribeProductConfig` | ✅ | 返回一站式规格配置 |

> 未覆盖：`UpgradeDBInstanceKernelVersion`、`CreateInstances`、`ModifyDBInstanceName`、`ModifyDBInstanceSpec`、`RestartDBInstance`、`IsolateDBInstances`、`DisIsolateDBInstances`

### 参数组（4/5 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeDBInstanceParameters` | ✅ | 返回实例运行参数列表 |
| `DescribeParamsEvent` | ✅ | 返回参数修改事件（当前 0 条） |
| `DescribeParameterTemplates` | ✅ | 返回参数模板列表 |
| `DescribeParameterTemplateAttributes` | ✅ | 返回指定模板详情 |

> 未覆盖：`ModifyDBInstanceParameters`

### SSL 组（1/1 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeDBInstanceSSLConfig` | ✅ | 返回 SSL 配置状态 |

### 账号组（2/6 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeAccounts` | ✅ | 返回实例账号列表 |
| `DescribeAccountPrivileges` | ✅ | 从 `DescribeAccounts` 提取 `UserName` 后查询 `postgres` 库权限 |

> 未覆盖：`CreateAccount`、`DeleteAccount`、`ModifyAccountPassword`、`ResetAccountPassword`

### 网络组（1/4 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeDBInstanceSecurityGroups` | ✅ | 返回安全组信息 |

> 未覆盖：`ModifyDBInstanceSecurityGroups`、`OpenDBExtranetAccess`、`CloseDBExtranetAccess`

### 监控组（3/3 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeSlowQueryList` | ✅ | 返回近 24h 慢查询列表（当前 0 条） |
| `DescribeSlowQueryAnalysis` | ✅ | 返回慢查询分析（当前 0 条） |
| `DescribeDBErrlogs` | ✅ | 返回错误日志（当前 0 条） |

### 数据库组（2/4 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeDatabases` | ✅ | 返回数据库列表 |
| `DescribeDatabaseObjects` | ✅ | 从 `DescribeDatabases` 提取 `DatabaseName` 后查询对象列表 |

> 未覆盖：`ModifyDatabaseOwner`、`CreateDatabase`

### 备份组（5/8 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeBackupOverview` | ✅ | 返回备份概览 |
| `DescribeBaseBackups` | ✅ | 返回基础备份集列表 |
| `DescribeLogBackups` | ✅ | 返回日志备份列表 |
| `DescribeAvailableRecoveryTime` | ✅ | 返回可恢复时间范围 |
| `DescribeCloneDBInstanceSpec` | ✅ | 从基础备份结果提取 `BackupSetId` 后返回克隆推荐规格 |

> 未覆盖：`CreateBaseBackup`、`DeleteBaseBackup`、`DescribeBackupDownloadURL`

### 只读实例组（1/2 已验证）

| 接口 | 结果 | 说明 |
|---|---|---|
| `DescribeReadOnlyGroups` | ✅ | 返回只读组列表（当前为空） |

> 未覆盖：`CreateReadOnlyDBInstance`

---

## 五、MCP 客户端协议级联调结果

执行 `./scripts/run_mcp_smoke.sh`，实际启动本地 SSE MCP server，并使用真实 `mcp-go` 客户端完成协议级 smoke test。

### 1. 服务端启动结果

| 检查项 | 结果 |
|---|---|
| Server 启动 | ✅ |
| 注册工具数 | ✅ `48` |
| SSE 监听 | ✅ `http://127.0.0.1:9000/sse` |

### 2. 客户端交互结果

| 步骤 | 结果 | 说明 |
|---|---|---|
| `initialize` | ✅ | 成功识别服务端 `腾讯云 Postgres MCP 1.0.0` |
| `ping` | ✅ | 返回 `status: ok` |
| `tools/list` | ✅ | 成功列出 **48 个工具** |
| schema spot check | ✅ | `CreateInstances`、`CreateReadOnlyDBInstance` 均暴露 `confirm` 参数 |
| `tools/call`（只读） | ✅ | `DescribeRegions`、`DescribeDBVersions`、`DescribeDBInstances`、`DescribeDBInstanceAttribute` 调用成功 |
| `tools/call`（写工具保护） | ✅ | `CreateInstances` 在 `confirm=false` 时返回 `require_confirm=true` |

### 3. 样例日志摘要

- **`initialize`**：返回服务端名称 `腾讯云 Postgres MCP`、版本 `1.0.0`
- **`tools/list`**：返回 `tool_count: 48`
- **`DescribeDBInstances`**：返回 `TotalCount=3`，说明实例列表按云侧真实数据返回
- **`CreateInstances`**：返回
  `{"code":403,"warning":"Tool 'CreateInstances' will incur costs. Set confirm=true to proceed.","require_confirm":true}`

**结论：** 从通用 MCP SSE 客户端视角看，当前 server 已具备正常的 **工具发现、工具 schema 暴露、只读工具调用、写工具保护提示** 能力。

---

## 六、验证过程中修正的问题

验证过程中发现 3 个问题属于**验证脚本传参不完整**，已修复，非 `tools/` 工具实现缺陷：

| 接口 | 问题 | 修复方式 |
|---|---|---|
| `DescribeClasses` | 缺必填参数 `Zone`、`DBMajorVersion` | 补充 `Zone=ap-chengdu-1`、`DBMajorVersion=18` |
| `DescribeAccountPrivileges` | 缺必填参数 `UserName`、`DatabaseObjectSet` | 串联 `DescribeAccounts` 结果提取 `UserName` |
| `DescribeCloneDBInstanceSpec` | 缺 `BackupSetId` / `RecoveryTargetTime` | 串联 `DescribeBaseBackups` 结果提取 `BackupSetId` |

此外，本次补充新增：
- `cmd/mcp_smoke`：MCP 协议级 smoke client
- `scripts/run_mcp_smoke.sh`：一键联调脚本，便于后续 review 与回归

---

## 七、边界说明

- 本报告已经验证 **MCP 协议层面** 的客户端兼容性，但**尚未覆盖某个具体 AI 桌面客户端/IDE 产品的 UI 行为差异**。
- 本报告未执行任何 `confirm=true` 的真实写操作，因此不会对云上资源产生实际修改。
- `.env` 中即使仍存在旧的 `INSTANCE_SCOPE` 配置项，当前代码已不再读取该字段；实际生效的范围限制仅剩 `REGION_SCOPE`。

---

## 八、最终结论

- **当前 48 个 MCP 工具均可正常构建、注册并暴露给客户端。**
- **工具参数与腾讯云 OpenAPI / SDK 请求结构保持一致（48/48 通过）。**
- **真实只读云 API 调用验证通过（27/27 通过）。**
- **MCP 客户端协议级联调通过，客户端可正常识别工具、列出工具并执行只读调用。**
- **写工具不会误执行，`confirm` 保护逻辑正常。**
- **实例访问范围限制已成功移除，当前行为符合预期。**

**综合判断：该 MCP server 当前已具备继续交付 review / 接入具体 AI 客户端进行最终体验验证的条件。**
