# 腾讯云 PG MCP v3 — 将要支持的 API 列表

> 数据来源：`PG_MCP_API支持列表_完整对照.xlsx` → `MCP支持列表(全112API)` sheet
> 范围：仅保留"腾讯云 PG MCP"列为 ✅ 的 API
> Tool 数：**48** 个（v2 为 44 个），覆盖率 **43%**（48/112）
> 适用：MCP 服务端开发 / 工具注册 / 文档与测试用例编写

---

## 目录

1. [实例相关接口（10）](#1-实例相关接口10)
2. [只读实例相关接口（2）](#2-只读实例相关接口2)
3. [备份与恢复相关接口（8）](#3-备份与恢复相关接口8)
4. [参数管理相关接口（5）](#4-参数管理相关接口5)
5. [网络相关接口（2）](#5-网络相关接口2)
6. [安全组相关接口（2）](#6-安全组相关接口2)
7. [性能优化相关接口（3）](#7-性能优化相关接口3)
8. [账号相关接口（6）](#8-账号相关接口6)
9. [数据库相关接口（4）](#9-数据库相关接口4)
10. [规格相关接口（5）](#10-规格相关接口5)
11. [任务相关接口（1）](#11-任务相关接口1)
12. [速览汇总](#速览汇总)
13. [v3 相对 v2 变更](#v3-相对-v2-变更)
14. [设计原则与不覆盖项](#设计原则与不覆盖项)

---

## 1. 实例相关接口（10）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 1 | `CreateInstances` | 创建实例 | 写（高危） |
| 2 | `DescribeDBInstances` | 查询实例列表 | 只读 |
| 3 | `DescribeDBInstanceAttribute` | 查询实例详情 | 只读 |
| 4 | `ModifyDBInstanceName` | 修改实例名字 | 写（低危） |
| 5 | `ModifyDBInstanceSpec` | 修改实例规格 | 写（中危） |
| 6 | `UpgradeDBInstanceKernelVersion` | 升级实例内核版本号 | 写（中危） |
| 7 | `RestartDBInstance` | 重启实例 | 写（中危） |
| 8 | `IsolateDBInstances` | 隔离实例 | 写（中危） |
| 9 | `DisIsolateDBInstances` | 解隔离实例 | 写（低危） |
| 10 | `DescribeDBInstanceSSLConfig` | 查询实例 SSL 配置 | 只读 |

---

## 2. 只读实例相关接口（2）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 11 | `DescribeReadOnlyGroups` | 查询只读组列表 | 只读 |
| 12 | `CreateReadOnlyDBInstance` | 创建只读实例 | 写（中危） |

---

## 3. 备份与恢复相关接口（8）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 13 | `DescribeAvailableRecoveryTime` | 查询实例可恢复的时间范围 | 只读 |
| 14 | `DescribeCloneDBInstanceSpec` | 查询克隆实例可购买的规格 | 只读 |
| 15 | `CloneDBInstance` | 克隆实例 | 写（中危） |
| 16 | `DescribeBackupOverview` | 查询备份概览 | 只读 |
| 17 | `CreateBaseBackup` | 创建实例数据备份 | 写（低危） |
| 18 | `DescribeBaseBackups` | 查询数据备份列表 | 只读 |
| 19 | `DescribeLogBackups` | 查询日志备份列表 | 只读 |
| 20 | `DescribeBackupDownloadURL` | 查询备份集的下载地址 | 只读 |

---

## 4. 参数管理相关接口（5）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 21 | `DescribeParameterTemplates` | 查询参数模板列表 | 只读 |
| 22 | `DescribeParameterTemplateAttributes` | 查询参数模板详情 | 只读 |
| 23 | `DescribeDBInstanceParameters` | 查询实例参数 | 只读 |
| 24 | `DescribeParamsEvent` | 查询参数修改事件 | 只读 |
| 25 | `ModifyDBInstanceParameters` | 修改实例参数 | 写（中危） |

---

## 5. 网络相关接口（2）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 26 | `OpenDBExtranetAccess` | 开通实例公网地址 | 写（中危） |
| 27 | `CloseDBExtranetAccess` | 关闭实例公网地址 | 写（中危） |

---

## 6. 安全组相关接口（2）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 28 | `DescribeDBInstanceSecurityGroups` | 查询实例安全组 | 只读 |
| 29 | `ModifyDBInstanceSecurityGroups` | 修改实例的安全组 | 写（中危） |

---

## 7. 性能优化相关接口（3）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 30 | `DescribeDBErrlogs` | 查询错误日志 | 只读 |
| 31 | `DescribeSlowQueryAnalysis` | 获取慢查询统计分析列表 | 只读 |
| 32 | `DescribeSlowQueryList` | 获取慢查询列表 | 只读 |

---

## 8. 账号相关接口（6）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 33 | `DescribeAccounts` | 查询实例的数据库账号列表 | 只读 |
| 34 | `ResetAccountPassword` | 重置账户密码 | 写（中危） |
| 35 | `CreateAccount` | 创建数据库账号 | 写（中危） |
| 36 | `DescribeAccountPrivileges` | 查询数据库账号的权限信息 | 只读 |
| 37 | `DeleteAccount` | 删除数据库账号 | 写（中危） |
| 38 | `ModifyAccountPrivileges` | 修改数据库账号的权限、类型 | 写（中危） |

---

## 9. 数据库相关接口（4）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 39 | `DescribeDatabases` | 查询实例的数据库列表 | 只读 |
| 40 | `CreateDatabase` | 创建数据库 | 写（中危） |
| 41 | `DescribeDatabaseObjects` | 查询数据库对象列表 | 只读 |
| 42 | `ModifyDatabaseOwner` | 修改数据库所有者 | 写（中危） |

---

## 10. 规格相关接口（5）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 43 | `DescribeRegions` | 查询售卖地域 | 只读 |
| 44 | `DescribeZones` | 查询售卖可用区 | 只读 |
| 45 | `DescribeClasses` | 查询售卖规格 | 只读 |
| 46 | `DescribeDBVersions` | 查询支持的数据库版本 | 只读 |
| 47 | `DescribeProductConfig` | 查询售卖规格配置 | 只读 |

---

## 11. 任务相关接口（1）

| # | API 名称 | 接口描述 | 操作类型 |
|---|---|---|---|
| 48 | `DescribeTasks` | 查询任务列表 | 只读 |

---

## 速览汇总

| 分类 | MCP 支持数 / 总数 | 覆盖率 |
|---|---|---|
| 实例相关接口 | 10 / 26 | 38% |
| 只读实例相关接口 | 2 / 10 | 20% |
| 备份与恢复相关接口 | 8 / 21 | 38% |
| 参数管理相关接口 | 5 / 9 | 56% |
| 网络相关接口 | 2 / 6 | 33% |
| 安全组相关接口 | 2 / 2 | **100%** |
| 性能优化相关接口 | 3 / 3 | **100%** |
| 账号相关接口 | 6 / 12 | 50% |
| 数据库相关接口 | 4 / 4 | **100%** |
| 规格相关接口 | 5 / 9 | 56% |
| 任务相关接口 | 1 / 1 | **100%** |
| 其他接口 | 0 / 1 | 0% |
| 数据库审计相关接口 | 0 / 8 | 0% |
| **合计** | **48 / 112** | **43%** |

### 只读 vs 写操作分布

| 类型 | 数量 | 占比 |
|---|---|---|
| 只读（安全，AI 可直接调用） | 27 | 56% |
| 写操作（需 guard 确认/二次校验） | 21 | 44% |

---

## v3 相对 v2 变更

### ❌ 移除（2 个）

| API 名称 | 移除原因 |
|---|---|
| `LockAccount` | 安全应急操作不应 AI 化——误锁正常账号会导致应用大面积报错，应由安全团队人工执行。 |
| `RestoreDBInstanceObjects` | 库表级恢复直接覆盖线上数据，风险过高。恢复操作应由 DBA 人工确认执行。 |

### ✅ 新增（6 个）

| API 名称 | 新增说明 |
|---|---|
| `DescribeCloneDBInstanceSpec` | 只读查询，克隆实例的前置依赖查询。与已有 `CloneDBInstance` 配对使用。 |
| `DescribeParamsEvent` | 只读查询，参数修改事件审计。排障时 DBA 经常需要查"谁在什么时候改了什么参数"。 |
| `DescribeAccountPrivileges` | 只读查询，查账号权限是修改权限前的必要步骤。与已有 `ModifyAccountPrivileges` 配对。 |
| `DescribeRegions` | 只读查询，建实例前选地域的基础信息查询。 |
| `DescribeZones` | 只读查询，与 `DescribeRegions` 配对——先选地域再选可用区。 |
| `DescribeProductConfig` | 只读查询，一站式规格配置查询。比单独查 Classes/DBVersions 更便捷。 |

### 统计

- Tool 数量：`44 → 48`（-2 +6）
- 覆盖率：`39% → 43%`

---

## 设计原则与不覆盖项

### 核心原则

> **AI 不替用户花钱、不做不可逆高危决策。**

### 不覆盖类别一览

| 不覆盖类别 | API 数量 | 不覆盖原因 |
|---|---|---|
| 计费/售卖/续费/询价 | 9 | 商业决策前置，AI 不应替用户花钱/下单 |
| 销毁实例（`DestroyDBInstance`） | 1 | 最高危不可逆——删了就没了，数据全丢 |
| 大版本升级 | 2 | 数据迁移风险极高，需 DBA 规划停机窗口 |
| HA/主备切换/删除保护 | 5 | 影响业务连续性，需人工决策 |
| 参数模板增删改 | 3 | 一次性配置，做好模板后用查询即可 |
| 审计日志（全部） | 8 | 合规配置不应 AI 化，审计开关由安全团队管理 |
| 备份增删/计划变更 | 7 | 删除备份不可逆；备份策略需人工评估 |
| 只读组管理（增删改/权重） | 8 | 负载均衡策略需人工调整 |
| 网络增删/RO 组网络 | 4 | 一次性配置，非 MCP 运维高频场景 |
| 账号 CAM/备注/解锁 | 6 | CAM 权限管理是安全红线，不应 AI 化 |
| 售卖价格/订单查询 | 5 | 计费关联，非 MCP 运维核心场景 |
| 其他（专属集群/SSL 修改/项目归属等） | 8 | 特定场景受众窄或一次性安全配置 |

### 友商对比下的腾讯云差异化 API

#### 友商均不支持（7 个，腾讯云独有）

| API 名称 | 差异化说明 |
|---|---|
| `DescribeDatabaseObjects` | 腾讯独有——让 AI 回答"这个库里有哪些表/视图/函数"，竞品做不到 |
| `DescribeDBInstanceSSLConfig` | SSL 合规查询，竞品均无 |
| `UpgradeDBInstanceKernelVersion` | 小版本升级为 DBA 例行运维，竞品均不做 |
| `DescribeCloneDBInstanceSpec` | 克隆前置查询，与 `CloneDBInstance` 配对 |
| `DescribeParamsEvent` | 参数变更审计追溯，排障刚需 |
| `DescribeAccountPrivileges` | 权限管理前置查询，与 `ModifyAccountPrivileges` 配对 |
| `DescribeProductConfig` | 一站式选型查询 |

#### 腾讯云+火山有，阿里云缺（4 个）

| API 名称 | 说明 |
|---|---|
| `IsolateDBInstances` | 阿里云缺隔离/解隔离，我们和火山均有 |
| `DisIsolateDBInstances` | 解除隔离=恢复，补救误操作 |
| `CreateBaseBackup` | 主动备份 DBA 刚需，阿里云缺 |
| `CreateReadOnlyDBInstance` | 建只读分流读流量，阿里云缺 |

---

## 友商覆盖度对比

| 产品 | 管控面 Tool 数 | API 总量 | 覆盖率 |
|---|---|---|---|
| **腾讯云 PG MCP（v3 目标）** | **48** | 112 | **43%** |
| 阿里云 RDS PG MCP | 56 | 200+ | ~25% |
| 阿里云 PolarDB PG MCP | 28 | 403 | ~7% |
| 火山 RDS PG MCP | 44 | ~250 | ~18% |
| AWS postgres-mcp | 7 | — | — |
| 腾讯云 PG MCP（v2 现状） | 9 | 112 | 8% |

---

> **给 codeBuddy 的提示**：
> 1. 本文档为 v3 MCP tool 注册清单的"权威来源"，请按 11 个分类逐组实现。
> 2. 写操作类（共 21 个）需在 MCP 层加 guard 二次确认（特别是 `CreateInstances`/`IsolateDBInstances`/`RestartDBInstance`/`UpgradeDBInstanceKernelVersion`/`CloneDBInstance`）。
> 3. 只读类（共 27 个）可直接放行。
> 4. 每个 tool 命名建议统一前缀，如 `pg_mcp_<api_snake_case>`，例如 `pg_mcp_describe_db_instances`。
