# PostgreSQL Skills 当前现状参考文档

## 1. 文档目的

这是一份面向**后续对话复用**的 PostgreSQL skill 现状快照，用于帮助其他对话快速理解：

- 当前支持哪些场景入口
- 每个场景入口负责什么、不负责什么
- 需要哪些输入槽位
- 遇到缺参、缺地域、缺凭证时应该如何回复
- 最近已经完成了哪些关键能力收敛

本文档基于当前 `src/postgres/skills/` 目录下的 `README.md`、`USAGE_GUIDE.md`、各 skill 的 `SKILL.md` 与 `references/` 整理。

## 2. 当前场景总览

当前 PostgreSQL skill 包已收敛为 **5 个场景**：

- **`tencent-pg-management`**：管控面统一入口
- **`tencent-pg-mem0-deploy`**：一句话开通 mem0 服务
- **`tencent-pg-rest-deploy`**：一句话开通 / 调用 / 排障 REST / PostgREST 服务
- **`tencent-pg-inspection`**：运维面巡检
- **`tencent-pg-slowquery-diagnosis`**：运维面慢 SQL 查询

可按能力分层理解为：

- **管控面**：`tencent-pg-management`
- **扩展服务开通**：`tencent-pg-mem0-deploy`、`tencent-pg-rest-deploy`
- **运维面观察**：`tencent-pg-inspection`、`tencent-pg-slowquery-diagnosis`

## 3. 全局共识与公共规则

### 3.1 统一输入前提

所有 skill 都优先依赖运行时环境中的腾讯云凭证，不应要求用户把密钥写进仓库或直接发到聊天中。

推荐环境变量最小集：

```bash
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
# 临时凭证场景再补 TENCENTCLOUD_SESSION_TOKEN
```

统一变量名要求：

- `TENCENTCLOUD_SECRET_ID`
- `TENCENTCLOUD_SECRET_KEY`
- `TENCENTCLOUD_REGION`
- 临时凭证场景再补 `TENCENTCLOUD_SESSION_TOKEN`
- 如果宿主使用自定义变量名，需要在触发 skill 前先映射到上述标准变量

### 3.2 地域归一化规则

当前公共地域规则已经固定：

- 优先使用标准地域码，例如 `ap-guangzhou`
- 接受常见中文别名：`广州`、`上海`、`成都`、`北京`
- 若运行时已有 `TENCENTCLOUD_REGION`，可作为默认地域
- 无法安全确认时必须停止，不能猜测性修正

当前沉淀的官方地域参考链接：

- PostgreSQL 产品地域页：`https://cloud.tencent.com/document/product/409/113656`
- 腾讯云通用地域页：`https://cloud.tencent.com/document/api/238/7520`
- PostgreSQL 售卖地域参考：`https://cloud.tencent.com/document/product/409/16768`

### 3.3 统一错误处理模板

公共错误处理已经沉淀到 `references/common/error_handling.md`，重点包括：

- **`missing-credentials`**：缺少 `SecretId` / `SecretKey` / `Region`
- **`invalid-region`**：地域无法安全归一化
- **`missing-target-scope`**：巡检/诊断类请求缺少地域或实例 ID
- **`missing-sdk`**：需要 SDK 但当前环境未检测到

这些模板的共同要求是：

- 先明确指出**缺什么**
- 直接给出**下一步修复动作**
- 能给官方链接时就给官方链接
- 对 **可安全代查的非密钥参数**，同时补一个“我可以代你查”的选项
- 对 **密钥 / Token / API Key** 这类敏感值，只给官方入口和配置方式，不承诺代取
- 涉及安装或修改环境时，不自动执行高风险动作

### 3.4 安全与执行边界

当前 skill 包遵循以下统一原则：

- **先查后动**：先做只读取证，再决定是否进入写动作
- **凭证只读运行时**：不从仓库文件、URL、查询参数读取密钥
- **范围不唯一不执行写动作**：目标实例不唯一时必须先收敛范围
- **高风险动作要明确确认**：尤其是管控面修改、费用类动作、安全影响动作
- **运维面场景当前保持只读**：不直接执行修复、调优或管理动作

## 4. 各场景入口详细现状

### 4.1 `tencent-pg-management`

#### 定位

这是 PostgreSQL **管控面统一入口**，负责把一句话自然语言任务路由到 4 条主 lane：

- `overview`
- `instance-change`
- `backup-recovery`
- `access-security`

#### 适用场景

适合以下类型的管理类请求：

- 查看实例状态、规格、版本、任务状态、只读组现状
- 评估或准备重启、升降配、隔离、解隔离、创建实例、创建只读实例
- 查看备份概览、恢复时间窗，或评估下载备份、克隆恢复、创建基础备份
- 查看账号权限、数据库 owner、公网访问、安全组、SSL 配置

#### 核心行为

- 从用户一句话中先识别**真实目标**，而不是按表面关键词硬匹配
- 提取最小必要槽位：地域、实例 ID、账号、数据库、恢复时间窗、目标规格等
- 优先走**最小 action 集合**，避免无关扩查
- 先执行只读 API，返回当前事实、阻塞项和可行性判断
- 即使只是查询恢复时间窗，也会明确提醒后续恢复 / 克隆 / 下载备份链接属于高风险或敏感动作，仍为 **`待确认`**
- 对任何仍需确认的写类、费用类、高风险动作，都会明确说明 **在等确认**、确认后会执行什么、这个动作的意义以及主要风险
- 写类、费用类、高风险动作统一标记为 **`待确认`**

#### 当前支持的主要路由 lane

- **`overview`**：实例状态、实例列表、任务、只读组现状
- **`instance-change`**：改名、重启、升降配、隔离/解隔离、内核升级、建实例
- **`backup-recovery`**：备份概览、恢复时间窗、克隆恢复、备份下载
- **`access-security`**：账号、权限、数据库、公网访问、安全组、SSL

#### 当前边界

- 不把巡检/慢 SQL 请求强行塞入管控面执行链路
- 不会在未确认时直接执行重启、升配、建库、改权限、开公网等动作
- 不会把多个高风险 lane 混在一个隐式执行流里

#### 典型输出

- 目标范围
- 识别意图和路由原因
- 提取槽位
- 已调用的 OpenAPI
- 当前事实 / 阻塞项 / 风险
- 明确的等待确认话术
- `待确认` 动作的意义与主要风险
- 下一步建议

### 4.2 `tencent-pg-mem0-deploy`

#### 定位

这是 **一句话开通 mem0 服务** 的专用 skill，目标是尽量一轮把服务推进到可用态。

#### 适用场景

典型请求如：

- `帮我给广州 postgres-abc12345 开通 mem0`
- `给当前实例一键部署 mem0，AgenticBaseId 用 ab-xxxxx`

#### 核心槽位

- `Region`
- `DBInstanceId`
- `AgenticBaseId`
- `LLMModel`（默认可用 `auto`）
- `EmbeddingApiKey`（只允许从运行时读取）

#### 核心行为

- 自动补齐非密钥槽位
- 优先使用显式实例 ID；必要时才做实例发现
- 先做只读预检查：实例存在性、实例适配性、当前 mem0 状态
- 当用户明确表达“开通/部署/启用 mem0”且目标唯一时，可直接执行 `OpenMem0Service`
- 轮询 `DescribeMem0Service`，直到 ready 或达到上限

#### 最近完成的重要增强

当前 `mem0` skill 已经加强了**缺参引导**，不再只给模糊提醒，而是会直接提供可点击的网站入口、最短点击路径、运行时配置方式和明确下一步：

- `Region`：给 PostgreSQL 控制台，并说明“右上角切地域 / 实例列表看所属地域”
- `AgenticBaseId`：给 PostgreSQL 控制台，并说明 `AI 应用 → AgenticBase → 详情复制 / 新建`
- `EmbeddingApiKey`：给混元控制台，并说明 `立即接入管理 → API Key 管理 → 创建 API KEY`
- `LLMModel`：默认推荐 `auto`，只有用户要求固定模型时才引导去混元控制台确认模型

#### 当前边界

- 不要求用户把 `EmbeddingApiKey` 贴到聊天里
- 不会在实例目标不唯一时执行 `OpenMem0Service`
- 不会对已运行 mem0 服务做静默重建或改配
- 不会自动调用 `CloseMem0Service` 作为回滚

#### 典型输出

- 目标范围
- 槽位值及其来源（用户输入 / 运行时默认 / 内置默认）
- 预检查事实
- 已执行动作
- 当前 mem0 状态
- `InnerAddress` 或等效可用地址
- 下一步最小使用示例

### 4.3 `tencent-pg-rest-deploy`

#### 定位

这是 **一句话开通 / 调用 / 排障 REST / PostgREST 服务** 的专用 skill，目标是在一轮内尽量完成服务开关、只读调用，或把 502 一类阻塞收敛到可执行的下一步。

#### 适用场景

典型请求如：

- `帮我给广州 postgres-abc12345 开通 REST 服务`
- `给当前实例一键部署 PostgREST`

#### 核心槽位

- `Region`
- `DBInstanceId`

#### 核心行为

- 自动补齐地域和实例默认值
- 先做只读预检查：实例存在性、实例适配性、当前 REST 状态
- 能识别 `开通/关闭服务`、`只读调用 REST 路径`、`排查 502` 三类主意图，并为每类任务选择最小工具集合
- 对安全的只读 GET 调用，可直接执行 `QueryPostgRESTService`
- 当下一步仍需用户确认时，会明确输出 `等你确认` 一类等待话术，并解释确认后动作的意义和主要风险
- 当用户明确表达开通意图且目标唯一时，可直接执行 `OpenPostgRESTService`；明确表达关闭意图时可执行 `ClosePostgRESTService`
- 针对 502 一类问题，会联动 `DescribeDBInstanceSecurityGroups` 做网络侧取证；若只是需要改现有安全组规则，会明确告知需手动修改，若只是切换绑定安全组，才会在确认后进入 `ModifyDBInstanceSecurityGroups`
- 轮询 `DescribePostgRESTService`，直到拿到可用访问地址、确认服务已关闭，或返回明确阻塞原因
- 已按更保守口径重跑一次全地域实扫（`2026-07-20T06:40:56Z`）：当前 `DescribeRegions` 返回的 16 个地域里，只有 **1 个地域**仍作为 skill 的默认自动候选且已被真实实例验证可查询 REST 状态：`ap-shanghai`；另有 **10 个 `AVAILABLE` 地域**目前仅达到“占位探针接受”层级，说明请求能打到实例校验，但**还不能直接写成已确认支持**：`ap-beijing`、`ap-guangzhou`、`ap-hongkong`、`ap-seoul`、`ap-shanghai-fsi`、`ap-shenzhen`、`ap-singapore`、`ap-tianjin`、`eu-frankfurt`、`na-siliconvalley`；还有 **4 个 `UNAVAILABLE` 地域**不应作为默认候选：`ap-guangzhou-open`、`ap-shenzhen-fsi`、`na-ashburn`、`na-toronto`。另外，`ap-chengdu` 虽已被真实实例复核到 `DescribePostgRESTService=Status:not_open`，但在补充 `OpenPostgRESTService` 级别验证前，当前 skill 记录**暂按不支持 / 不纳入自动候选处理**。原始扫描记录保存在 `tmp/postgrest-region-scan.json`

#### 当前边界

- 不要求用户把云密钥发到聊天里
- 缺少 `Region` / `DBInstanceId` 这类非密钥参数时，可以提供代查选项；但敏感值仍不会代取
- 不会在目标不唯一时直接开通
- 不会自动执行 `ClosePostgRESTService`
- 不会无限轮询

#### 典型输出

- 目标范围
- 槽位来源
- 预检查事实
- 已执行动作
- 最终 REST 状态
- 访问地址 / endpoint
- 地域不支持时的候选地域列表
- 需要确认时的明确等待话术
- 待确认动作的意义与主要风险
- 下一步使用方式

### 4.4 `tencent-pg-inspection`

#### 定位

这是 PostgreSQL **运维面巡检** skill，当前是**固定范围、只读、报告式输出**。

#### 适用场景

适合对单个实例做基础健康巡检，例如查看：

- CPU
- 内存
- 存储
- 连接数
- I/O
- 复制延迟

#### 核心行为

- 先确认地域、实例 ID、可选时间窗
- 用监控接口发现支持指标，再拉取固定指标集合
- 不扩张到备份、参数、安全组、账号、SSL、慢 SQL 等其它领域
- 输出结构已调整为更接近运维报告，而不是松散指标堆砌

#### 当前输出结构

- 巡检摘要
- 巡检对象
- 健康快照
- 指标明细
- 风险与人工复核项
- 数据说明

#### 当前边界

- 不做修复建议
- 不做管理动作
- 不做根因分析
- 不杜撰阈值或不可得指标

### 4.5 `tencent-pg-slowquery-diagnosis`

#### 定位

这是 PostgreSQL **运维面慢 SQL 查询** skill，当前也是**固定范围、只读、报告式输出**。

#### 适用场景

适合查看某个实例在指定时间窗内的慢 SQL 基础事实：

- SQL 文本 / 归一化 SQL
- 数据库名
- 用户名
- 客户端地址
- 执行耗时
- 执行次数
- 总耗时
- 执行时间

#### 核心行为

- 先确认地域、实例 ID、时间窗
- 只使用 `DescribeSlowQueryList` 与 `DescribeSlowQueryAnalysis`
- 输出结构化慢 SQL 结果，不进入调优或排障链路

#### 当前输出结构

- 查询摘要 / 查询范围
- 核心慢 SQL 发现
- 详细慢 SQL 列表
- 人工复核项
- 数据说明

#### 当前边界

- 不做根因判断
- 不做原因排序
- 不输出调优建议、参数修改建议、升配建议
- 不扩展到错误日志、安全组、备份等无关模块

## 5. 当前能力收敛重点

截至当前状态，skill 包已经完成以下关键收敛：

### 5.1 从“多入口碎片化”收敛到“5 个明确场景”

对话里不需要再先手动拆很多小能力，当前可直接按 5 个场景的定位做路由。

### 5.2 管控面统一到 `tencent-pg-management`

大多数实例状态、实例变更、备份恢复、访问安全问题，都应优先考虑走 `tencent-pg-management`，而不是继续拆新 skill。

### 5.3 运维面输出升级为“报告体”

`PG巡检` 与 `慢SQL查询` 已不再追求简单列表，而是要求输出更接近运维报告：有摘要、有范围、有明细、有数据说明，但仍保持保守、不做过度推断。

### 5.4 mem0 缺参引导已显著增强

`AgenticBaseId`、`EmbeddingApiKey`、`LLMModel` 已支持直接控制台入口引导；地域缺失时也会优先给 PostgreSQL 控制台，而不是先丢文档页链接。

### 5.5 地域缺参 / 非法地域处理已统一

遇到地域问题时，不应只说“region 不对”或“自己去控制台找”，而应给：

- 用户原始输入回显
- 合法示例
- PostgreSQL 控制台入口
- “右上角切地域 / 实例列表查所属地域”的最短操作
- 一个可直接复制的 `export TENCENTCLOUD_REGION="ap-guangzhou"`

## 6. 其它对话可直接复用的路由建议

### 6.1 看到这些诉求，优先用哪个 skill

- **查看实例状态 / 评估规格 / 看恢复时间窗 / 看权限 / 看 SSL / 看公网**
  - 优先：`tencent-pg-management`
- **明确要求开通 mem0 / 长期记忆服务**
  - 优先：`tencent-pg-mem0-deploy`
- **明确要求开通 REST / PostgREST 服务**
  - 优先：`tencent-pg-rest-deploy`
- **查看实例基础健康 / 监控指标 / 资源情况**
  - 优先：`tencent-pg-inspection`
- **查看慢 SQL / 慢查询列表 / Top 慢 SQL**
  - 优先：`tencent-pg-slowquery-diagnosis`

### 6.2 其它对话里的常见回复策略

- **缺地域**：直接给地域参考链接和标准地域码示例
- **缺实例 ID**：引导用户去 PostgreSQL 控制台复制 `postgres-xxxxxxxx`
- **缺凭证**：给最小环境变量模板，不要求用户把密钥发到聊天中
- **目标不唯一**：停止执行，要求用户补 `region + instance ID`
- **请求属于运维观察，不属于管理动作**：从 `tencent-pg-management` 分流到运维 skill
- **请求属于高风险写动作**：即便走 `tencent-pg-management`，也必须先只读取证，再标记 `待确认`

## 7. 打包与分发现状

当前 `skills/package.json` 已提供标准打包命令：

```bash
npm run verify
npm run release
```

其中：

- `verify`：校验 skill 结构
- `release`：执行清理、校验和整包打包

预期产物包括：

- `tencent-pg-management-v1.0.3.zip`
- `tencent-pg-mem0-deploy-v1.0.3.zip`
- `tencent-pg-rest-deploy-v1.0.3.zip`
- `tencent-pg-inspection-v1.0.3.zip`
- `tencent-pg-slowquery-diagnosis-v1.0.3.zip`
- `tencentdb-postgresql-skill-v1.0.3.zip`

## 8. 后续对话使用这份文档时的建议

如果其他对话需要快速接手 PostgreSQL skill 相关工作，建议优先从这几个问题入手：

1. 当前请求属于 **管控面**、**扩展服务开通**，还是 **运维观察**？
2. 用户是否已经提供 **地域 + 实例 ID**？
3. 当前请求是**只读查询**，还是带有**写动作 / 费用影响 / 安全影响**？
4. 是否存在现成的**官方获取链接**可直接给用户，而不是让用户自己搜索？
5. 是否应该保持**事实型、保守型输出**，避免根因推断或过度建议？

## 9. 一句话总结

当前 PostgreSQL skill 体系已经形成了一个比较清晰的结构：**1 个管控面统一入口 + 2 个专项开通 skill + 2 个只读运维 skill**。对后续对话而言，最重要的是按场景正确路由、优先收敛地域与实例范围、严格遵守只读/确认边界，并在缺参时直接返回官方链接和可复制模板。
