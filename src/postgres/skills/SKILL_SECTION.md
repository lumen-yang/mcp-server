# PostgreSQL Skill

PostgreSQL Skill 是面向云数据库 PostgreSQL 的统一技能入口，当前已经收敛为 **5 个场景**：**管控面统一入口**、**mem0 一句话部署**、**REST 一句话开通 / 调用 / 排障**、**运维面巡检**、**运维面慢 SQL 查询**。

其中：

- `tencent-pg-management` 负责管控面统一入口
- `tencent-pg-mem0-deploy` 负责 mem0 的一句话开通
- `tencent-pg-rest-deploy` 负责 REST 的一句话开通 / 调用 / 排障
- `tencent-pg-inspection` 与 `tencent-pg-slowquery-diagnosis` 属于运维面场景，当前保持只读

## 1. 使用场景与获取地址

| 场景 | 说明 |
|---|---|
| 管控面统一入口 | 根据自然语言自动识别实例总览、实例变更、备份恢复、访问安全意图，并路由到最匹配的 OpenAPI |
| mem0 一句话部署 | 自动补齐部署槽位，先查状态，再把 mem0 服务开通到可用 |
| REST 一句话开通 / 调用 / 排障 | 自动补齐地域与实例，先做只读取证，再完成 REST 服务开关、只读调用或 502 排障 |
| 运维面巡检 | 查看固定监控指标的巡检结果 |
| 运维面慢 SQL 查询 | 查看固定时间窗内的慢 SQL 基础信息 |

### 1.2 获取地址

| 入口 | 说明 | 地址 |
|---|---|---|
| GitHub 仓库 | 查看源码与文档 | [前往 GitHub 仓库](https://github.com/TencentCloudCommunity/mcp-server/tree/feat/postgres-stdio-npx-support/src/postgres/skills) |
| ClawHub | 通过市场页面查看与分发 | [打开 ClawHub 页面](https://clawhub.ai/tencent-adm/tencentdb-postgresql-skill) |
| SkillHub | 通过技能市场查看详情 | [打开 SkillHub 页面](https://skillhub.cn/skills/tencentdb-postgresql-skill) |

## 2. 功能介绍

| 功能 | 说明 | 输出内容 |
|---|---|---|
| 管控面统一入口 | 覆盖实例状态查看、实例变更、备份恢复、访问安全治理，并在 skill 内按意图路由最小 API 集合 | 当前事实、识别意图、提取槽位、调用的 API、风险与下一步；需要确认时的明确等待话术、待确认动作的意义与主要风险 |
| mem0 一句话部署 | 覆盖目标实例识别、只读预检查、`OpenMem0Service` 执行和轮询直到可用 | 目标范围、槽位来源、当前 mem0 状态、可直接使用的地址、缺参时的官方获取链接与可代查选项 |
| REST 一句话开通 / 调用 / 排障 | 覆盖目标实例识别、只读预检查、`OpenPostgRESTService` / `ClosePostgRESTService`、`QueryPostgRESTService` 只读代执行，以及 502 安全组排障 | 目标范围、槽位来源、当前 REST 状态、只读 HTTP 结果、可直接使用的访问地址、502 排障结论、需要确认时的明确等待话术与动作风险 |
| 运维面巡检 | 基于腾讯云监控接口查看固定指标集合 | 巡检摘要、健康快照、指标明细、风险与人工复核项 |
| 运维面慢 SQL 查询 | 基于只读慢查询接口查看慢 SQL 信息 | 查询摘要、核心发现、慢 SQL 明细、人工复核项 |

## 3. 使用方法

推荐输入信息：

```text
地域：如 ap-guangzhou
实例 ID：如 postgres-abc12345
自然语言任务：如 查看实例状态并评估是否适合升级规格
可选范围：如 最近 1 小时 / 最近 24 小时 / 指定数据库 / 目标规格 / 是否公网 / AgenticBaseId
```

## 4. 安全边界

| 边界 | 说明 |
|---|---|
| 先识别意图再调 API | `tencent-pg-management` 先判断任务属于哪一类管控意图，再选择最小 API 集合 |
| 专项部署 skill 先查后开 | `tencent-pg-mem0-deploy` 与 `tencent-pg-rest-deploy` 都会先做只读检查，再决定是否执行开通 |
| 高风险或目标不明确时不盲写 | 目标不唯一、条件不满足或阻塞未解除时不会直接执行写动作 |
| 运维面场景当前保持只读 | `PG巡检` 与 `慢SQL查询` 不直接执行管理动作 |
| 密钥仅在运行时使用 | 不将凭证写入仓库文件或 URL |
