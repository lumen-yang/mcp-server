# PostgreSQL 独立技能包 — 使用指南

> English Version: [USAGE_GUIDE_EN.md](./USAGE_GUIDE_EN.md)

本文档介绍当前 PostgreSQL 技能包的安装与使用方式。当前技能包已经收敛为 **5 个场景**：

- **管控面场景**：`tencent-pg-management`
- **扩展服务场景**：`tencent-pg-mem0-deploy`、`tencent-pg-rest-deploy`
- **运维面场景**：`tencent-pg-inspection`、`tencent-pg-slowquery-diagnosis`

其中：

- `tencent-pg-management` 负责把一句话自然语言任务自动路由到实例总览、实例变更、备份恢复、访问安全 4 类 OpenAPI 能力
- `tencent-pg-mem0-deploy` 负责把一句话 mem0 开通 / 关闭请求补齐为可执行链路，并返回最终服务状态
- `tencent-pg-rest-deploy` 负责把一句话 REST / PostgREST 开通、只读调用、502 排障请求补齐为可执行链路，并返回最终服务状态或排障结论
- 两个运维面场景负责运维观察类任务，当前保持只读，不承接复杂排障或执行链路

---

## 1. 安装前准备

### 1.1 获取腾讯云 API 密钥

1. 打开 [腾讯云 API 密钥管理](https://console.cloud.tencent.com/cam/capi)
2. 如果还没有可用密钥，点击 `新建密钥` / `创建密钥`
3. 推荐优先使用 **最小权限 CAM 子账号**，不要长期复用高权限主账号密钥
4. 创建后立即保存 `SecretId` 与 `SecretKey`
   - `SecretId` 可后续查看
   - `SecretKey` 通常只在创建时完整展示一次，丢失后需要重新创建并轮换旧密钥
5. 再打开 [PostgreSQL 控制台](https://console.cloud.tencent.com/postgres) 确认目标地域，优先记下标准地域码，例如 `ap-guangzhou`

> 不要把 `SecretKey`、`SessionToken` 或其他敏感值发到聊天、代码仓库、截图或工单里。

### 1.2 设置环境变量

下面这几项是统一变量名：

- `TENCENTCLOUD_SECRET_ID`
- `TENCENTCLOUD_SECRET_KEY`
- `TENCENTCLOUD_REGION`
- 临时凭证场景再补 `TENCENTCLOUD_SESSION_TOKEN`

#### 先判断你属于哪类使用方式

| 你当前怎么用 | 常见客户端 / 场景 | 推荐方式 |
|---|---|---|
| 在终端直接启动宿主进程、CLI、本地调试命令 | `node`、`npm run ...`、本地脚本 | 方式 A |
| macOS 图形客户端，从图标或 Finder / Dock 启动 | `WorkBuddy`、`CodeBuddy`、`Claude Desktop`、`Cherry Studio`、`Chatbox` 等 | 方式 B |
| IDE / 编辑器类客户端 | `Cursor`、`VS Code`、`Windsurf`、`Trae` 等 | 从终端启动用方式 A；从图标启动用方式 B；仅集成终端跑命令时也可用方式 A |
| Windows 桌面客户端 | Windows 上的桌面 AI 客户端、IDE、Electron App | 方式 D |

#### 方式 A：CLI / 命令行启动

适合你准备 **从命令行直接启动宿主进程、CLI 或本地调试命令** 的场景。

1. 在你接下来要执行启动命令的终端里，先运行：

```bash
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
# 临时凭证场景再补：
export TENCENTCLOUD_SESSION_TOKEN="你的 SessionToken"
```

2. 紧接着验证：

```bash
echo $TENCENTCLOUD_SECRET_ID
echo $TENCENTCLOUD_REGION
```

3. 验证通过后，直接继续执行你的启动命令

> 这种方式只适合 **命令行启动**。如果你使用的是桌面客户端，不要把“临时 `export` 一次”当成主要方案。

#### 方式 B：macOS 图形客户端（WorkBuddy / CodeBuddy / Claude Desktop / Cherry Studio / Chatbox 等）

适合你通过 macOS 桌面客户端使用这些 skill，而不是从 CLI 直接启动。

1. 打开终端，执行：

```bash
launchctl setenv TENCENTCLOUD_SECRET_ID "你的 SecretId"
launchctl setenv TENCENTCLOUD_SECRET_KEY "你的 SecretKey"
launchctl setenv TENCENTCLOUD_REGION "ap-guangzhou"
# 临时凭证场景再补：
launchctl setenv TENCENTCLOUD_SESSION_TOKEN "你的 SessionToken"
```

2. 执行下面命令验证是否已经写入当前登录会话环境：

```bash
launchctl getenv TENCENTCLOUD_SECRET_ID
launchctl getenv TENCENTCLOUD_REGION
```

3. **完全退出** 客户端后再重新打开
4. 重新进入 skill，再继续你的操作

> 对桌面客户端来说，单独在某个终端里执行一次 `export`，通常不会自动传给已经打开或从桌面图标启动的客户端进程；`launchctl setenv` 更适合作为 macOS 图形客户端的配置方式。

#### 方式 C：IDE / 编辑器客户端（Cursor / VS Code / Windsurf / Trae 等）

这类客户端要先分清楚：**到底是 IDE 自己触发 skill，还是你在 IDE 集成终端里启动宿主进程**。

##### 场景 1：你在 IDE 集成终端里自己启动命令

这种情况本质上还是 **方式 A**，直接在那个终端里设置变量即可，例如：

```bash
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
```

然后继续在同一个终端里执行你的命令。

##### 场景 2：你希望 IDE 主进程及其插件都能读到变量

如果客户端是从终端启动的，可以这样做：

```bash
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"

# 按你实际使用的客户端选择其一
cursor .
# code .
# windsurf .
# trae .
```

如果客户端平时是从 Dock / Finder / 桌面图标启动的，则回到 **方式 B**，先用 `launchctl setenv` 注入登录会话环境，再彻底重启客户端。

> 简单理解：**谁启动，变量就要给谁**。IDE 集成终端跑命令，就给终端；IDE 主进程和插件要读，就给 IDE 主进程所在会话。

#### 方式 D：Windows 客户端

适合你在 Windows 上使用桌面客户端、IDE 或自定义宿主。

##### 方案 1：写入用户级环境变量（推荐）

在 PowerShell 中执行：

```powershell
setx TENCENTCLOUD_SECRET_ID "你的 SecretId"
setx TENCENTCLOUD_SECRET_KEY "你的 SecretKey"
setx TENCENTCLOUD_REGION "ap-guangzhou"
# 临时凭证不建议长期 setx；如确有需要，可改用当前会话注入
```

执行后：

1. 关闭并重新打开你的终端 / 客户端
2. 再执行：

```powershell
echo $env:TENCENTCLOUD_SECRET_ID
echo $env:TENCENTCLOUD_REGION
```

##### 方案 2：仅给当前会话临时注入

```powershell
$env:TENCENTCLOUD_SECRET_ID = "你的 SecretId"
$env:TENCENTCLOUD_SECRET_KEY = "你的 SecretKey"
$env:TENCENTCLOUD_REGION = "ap-guangzhou"
$env:TENCENTCLOUD_SESSION_TOKEN = "你的 SessionToken"
```

这更适合本次调试或临时会话。

> `setx` 不会自动改写已经打开的当前窗口；配完后要重开终端或客户端。

#### 其他宿主说明

对于 Docker、CI、自建宿主或宿主已有自定义变量名的场景，当前不再单列独立方式。

统一只遵循一个原则：**在真正启动宿主进程前，把标准 `TENCENTCLOUD_*` 变量注入到实际触发 skill 的进程环境中**。如果你内部已经有其他变量名，也请先在启动链路里映射到标准变量，再继续启动客户端或宿主。

#### 常见误区

- **误区 1**：在一个终端里 `export` 过，桌面客户端就一定能读到  
  不一定。很多从图标启动的 GUI 客户端并不会继承那个终端里的环境。
- **误区 2**：把 `SecretKey` 直接写进客户端配置文件或仓库  
  不建议。优先通过系统环境变量、平台 Secret 或运行时注入来提供。
- **误区 3**：`setx` 或 `launchctl setenv` 配完后不重启客户端  
  很多客户端只有在重新启动后才会重新读取环境。
- **误区 4**：把临时 `SessionToken` 当成长期凭证  
  临时凭证会过期，过期后需要重新注入。

---

## 2. 安装技能

### 2.1 下载

从 Release 下载以下 zip：

- `tencent-pg-management-v1.0.3.zip`
- `tencent-pg-mem0-deploy-v1.0.3.zip`
- `tencent-pg-rest-deploy-v1.0.3.zip`
- `tencent-pg-inspection-v1.0.3.zip`
- `tencent-pg-slowquery-diagnosis-v1.0.3.zip`
- 或 `tencentdb-postgresql-skill-v1.0.3.zip`

### 2.2 导入到 CodeBuddy / WorkBuddy

1. 打开技能管理
2. 导入目标 `.zip`
3. 启用技能

### 2.3 验证安装

安装完成后，可直接输入：

- `查看实例状态并评估是否适合升级规格` → 应触发 `tencent-pg-management`
- `帮我给广州 postgres-abc12345 开通 mem0` → 应触发 `tencent-pg-mem0-deploy`
- `帮我关闭广州 postgres-abc12345 的 mem0` → 应触发 `tencent-pg-mem0-deploy`
- `帮我给广州 postgres-abc12345 开通 REST 服务` → 应触发 `tencent-pg-rest-deploy`
- `帮我关闭广州 postgres-abc12345 的 REST 服务` → 应触发 `tencent-pg-rest-deploy`
- `帮我调用这个 REST 路径` → 应触发 `tencent-pg-rest-deploy`
- `帮我排查 REST 502` → 应触发 `tencent-pg-rest-deploy`
- `PG巡检` → 应触发 `tencent-pg-inspection`
- `慢SQL查询` → 应触发 `tencent-pg-slowquery-diagnosis`

---

## 3. 管控面统一入口

### 3.1 使用场景

适合通过一句话自然语言直接发起以下管理任务：

- 查看实例状态、规格、版本、任务状态、只读组现状
- 评估或准备重启、升降配、隔离、解隔离、创建实例、创建只读实例
- 查看备份概览、恢复时间窗，或评估创建备份、下载备份、克隆恢复
- 查看账号权限、数据库 owner、公网访问、安全组、SSL 配置

### 3.2 示例输入

- `查看广州 postgres-abc12345 的实例状态`
- `帮我评估 ap-guangzhou postgres-abc12345 是否适合升级规格`
- `帮我看上海 postgres-abc12345 的恢复时间窗`
- `检查北京 postgres-abc12345 的账号权限和 SSL 配置`
- `重置 postgres-abc12345 某个账号密码前先帮我看下当前权限`

### 3.3 行为说明

该 skill 会：

1. 提取主目标与关键槽位，例如地域、实例 ID、账号、数据库、恢复时间窗、目标规格等
2. 归一化地域
3. 根据真实任务目标把请求路由到实例总览、实例变更、备份恢复、访问安全 4 条主 lane 之一
4. 为当前 lane 构造最小 API 计划，并先执行只读取证阶段
5. 先返回当前事实、阻塞项和可行性判断，再决定是否进入待确认动作
6. 对写类、费用类或高风险动作明确标记 `待确认`；如果下一步仍需用户确认，会明确说在等确认，并解释确认后动作的意义和主要风险

该 skill **不会**：

- 强制用户先手动选择“实例总览 / 实例变更 / 备份恢复 / 访问安全”中的某一个具体场景
- 在缺少必要范围信息时盲目调用 API
- 在没有明确确认时直接执行重启、升配、建库、重置密码、开公网等动作
- 把运维面巡检或运维面慢 SQL 观察请求硬塞进管控面执行链路

### 3.4 输出格式

```text
一、目标范围：地域 / 实例 / 可选对象
二、识别意图：实例总览 / 实例变更 / 备份恢复 / 访问安全
三、提取槽位：影响本次路由和 API 选择的关键字段
四、调用 API：本次用于取证或执行的 OpenAPI
五、当前事实 / 结果：与任务直接相关的结构化结论
六、风险与确认项：待确认动作、费用影响、安全影响；若下一步仍需确认，要明确说在等确认、确认后会执行什么、该动作的意义和主要风险；若本次只是查询恢复时间窗，也要明确提醒后续恢复 / 克隆 / 下载备份链接仍属高风险待确认动作
七、下一步建议
```

### 3.5 关联专用 skill：一句话开通 / 关闭 mem0 服务

适合这种目标非常明确的请求：

- `帮我给广州 postgres-abc12345 开通 mem0`
- `给当前实例一键部署 mem0，AgenticBaseId 用 ab-xxxxx`
- `帮我关闭广州 postgres-abc12345 的 mem0`

该 skill 会：

1. 自动识别这是 `开通` 还是 `关闭` 请求，并补齐地域、实例，以及开通路径所需的 `AgenticBaseId`、`LLMModel` 等槽位
2. 优先读取运行时环境中的默认值与密钥，不要求把密钥发到聊天里
3. 缺少关键参数时，会直接返回可点击的控制台 / 产品网站入口、最短点击路径、运行时环境配置示例，以及“配好后你下一句可以怎么说”的明确继续话术；像 `Region`、实例 ID 这类非密钥参数还会补充“我可以代你查”的选项，`EmbeddingApiKey` 这类敏感值仍只提供官方入口与运行时注入方式
4. 先做 `DescribeDBInstanceAttribute` / `DescribeMem0Service` 只读检查
5. 如果目标动作是开通，就执行 `OpenMem0Service`；如果目标动作是关闭，就执行 `CloseMem0Service`
6. 轮询服务状态，直到拿到可直接使用的 mem0 地址、确认服务已关闭，或返回明确阻塞原因

该 skill 的目标不是“给方案”，而是 **尽量一轮把 mem0 服务开通到可用，或关闭到最终状态**。

#### 参数获取入口

- `AgenticBaseId`：直接打开 PostgreSQL 控制台 `https://console.cloud.tencent.com/postgres`
  - 点击路径：`AI 应用` → 选择地域 → `AgenticBase`
  - 如果已有 Base：进入详情后直接复制 `AgenticBaseId`
  - 如果还没有：点击 `新建`
- `EmbeddingApiKey`：直接打开混元控制台 `https://console.cloud.tencent.com/hunyuan`
  - 点击路径：`立即接入管理` → `API Key 管理` → `创建 API KEY`
  - 创建后把 Key 放到运行时环境，不要发到聊天里
- `LLMModel`：如无特殊要求，可直接使用默认 `auto`
  - 只有你想固定模型时，再打开混元控制台 `https://console.cloud.tencent.com/hunyuan` 查看当前可接入模型

### 3.6 关联专用 skill：一句话开通 / 关闭 REST 服务

适合这种目标非常明确的请求：

- `帮我给广州 postgres-abc12345 开通 REST 服务`
- `给当前实例一键部署 PostgREST`
- `帮我关闭广州 postgres-abc12345 的 REST 服务`

该 skill 会：

1. 自动识别这是 `开通` 还是 `关闭` 请求，并补齐地域与实例 ID，优先读取运行时环境中的默认值
2. 先做 `DescribeDBInstanceAttribute` / `DescribePostgRESTService` 只读检查；如果实例所在地域不支持 REST，会先返回可继续尝试的地域列表
3. 开通时默认执行 `OpenPostgRESTService(DBInstanceId, EnableWanNet=false)`，也就是**先开 REST 服务，但默认不打开公网 / 外网访问**；关闭请求则执行 `ClosePostgRESTService`
4. 轮询服务状态，直到确认 REST 已可用、服务已关闭，或返回明确阻塞原因；如果本次是默认无公网开通，也会明确说明“服务已开启，但公网仍保持关闭”
5. 输出目标范围、槽位来源、执行的 API、`EnableWanNet` 的实际取值、当前状态、地域不支持时的候选地域和下一步使用方式；如果用户明确要求代为打开公网 / 外网访问，则会先停在“等你确认”的状态，并解释该动作的意义和风险

该 skill 的目标同样不是“给步骤”，而是 **尽量一轮把 REST / PostgREST 服务开通到可用，或关闭到最终状态**。

---

## 4. 运维面 PG 巡检

### 4.1 使用场景

适合查看单个实例的基础巡检结果，例如：

- CPU
- 内存
- 存储
- 连接数
- I/O
- 复制延迟

### 4.2 示例输入

- `PG巡检 ap-guangzhou postgres-abc12345`
- `健康检查 广州 postgres-abc12345 最近 1 小时`
- `监控巡检 上海 postgres-abc12345`

### 4.3 行为说明

该 skill 会：

1. 确认地域、实例 ID、可选时间窗
2. 归一化地域
3. 通过腾讯云监控接口发现支持的指标
4. 拉取固定指标集合
5. 返回基础巡检结果

该 skill **不会**：

- 扩查备份、账号、参数、安全组、SSL
- 进入排障流程
- 输出修复动作
- 自动做原因分析

### 4.4 输出格式

```text
一、巡检摘要：总体状态 / 关键发现 / 指标覆盖情况
二、巡检对象：地域 / 实例 / 时间范围 / 调用的监控 API
三、健康快照：CPU / 内存 / 存储 / 连接 / I/O / 复制延迟等核心指标的状态摘要
四、指标明细：指标名 / 数值 / 单位 / 数据状态 / 事实说明
五、风险与人工复核项：attention / abnormal / manual review needed 项
六、数据说明：unsupported / no-data 指标、时间窗说明
```

---

## 5. 运维面慢 SQL 查询

### 5.1 使用场景

适合查看单个实例在指定时间窗内的慢 SQL 基础信息，例如：

- Top N 慢 SQL
- SQL 文本或抽象 SQL
- 数据库名
- 用户名
- 客户端地址
- 执行耗时
- 执行次数
- 总耗时

### 5.2 示例输入

- `慢SQL查询 ap-guangzhou postgres-abc12345 最近 1 小时`
- `查看广州 postgres-abc12345 的慢查询`
- `慢SQL分析 上海 postgres-abc12345 昨天下午 3 点到 5 点`

### 5.3 行为说明

该 skill 会：

1. 确认地域、实例 ID、时间窗
2. 归一化地域
3. 调用 `DescribeSlowQueryList`
4. 在需要聚合排序时调用 `DescribeSlowQueryAnalysis`
5. 返回慢 SQL 基础信息

该 skill **不会**：

- 进行原因排序
- 做根因分析
- 关联错误日志进行推断
- 输出调优、升配、改参数建议

### 5.4 输出格式

```text
一、查询范围：地域 / 实例 / 时间窗口 / 排序方式
二、慢 SQL 基础列表：SQL / 数据库 / 用户 / 耗时 / 次数 / 总耗时 / 时间
```

---

## 6. 常见问题

### Q1：提示缺少凭证怎么办？

确认环境变量已设置：

```bash
echo $TENCENTCLOUD_SECRET_ID
echo $TENCENTCLOUD_SECRET_KEY
echo $TENCENTCLOUD_REGION
```

如果你用的是桌面客户端，再额外确认：**真正触发 skill 的那个进程** 是否能读到这些变量，而不只是某个临时终端能读到。

### Q2：支持中文地域吗？

支持常见中文地域，例如：

- `广州` → `ap-guangzhou`
- `上海` → `ap-shanghai`
- `成都` → `ap-chengdu`
- `北京` → `ap-beijing`

### Q3：为什么没有直接执行变更动作？

这是当前统一管控 skill 的设计要求：

- 先根据真实任务目标选对 lane 与 API
- 先返回只读事实
- 写类、费用类或高风险动作必须明确确认

### Q4：为什么 PG 巡检和慢 SQL 还单独保留？

这是当前的分层设计：

- `tencent-pg-management` 负责管控面任务的自然语言路由与受控执行
- `tencent-pg-inspection` 与 `tencent-pg-slowquery-diagnosis` 属于运维面场景，负责稳定返回巡检与慢 SQL 观察结果

这样可以减少误路由，提升结果稳定性。

---

## 7. 注意事项

1. 首次请求建议直接带上 `地域 + 实例 ID`
2. `tencent-pg-management` 的目标是 **一句话直达管理 API**，不是复杂排障工作流
3. 如需更复杂能力，应继续扩展当前独立 skill 的直连 OpenAPI / Monitor 动作，而不是引入额外前置服务
4. 如果你需要英文说明，可直接查看：[USAGE_GUIDE_EN.md](./USAGE_GUIDE_EN.md)
