# PostgreSQL 技能包使用说明

> English Version: [README_EN.md](./README_EN.md)

如果你想在 CodeBuddy / WorkBuddy 等 AI 客户端中直接使用腾讯云数据库 PostgreSQL 技能，先看这份 README 就够了。

当前可直接安装和使用的能力已经收敛为 **5 个技能包**：

- **管控面统一入口**：`tencent-pg-management`
- **扩展服务场景**：`tencent-pg-mem0-deploy`、`tencent-pg-rest-deploy`
- **运维观察场景**：`tencent-pg-inspection`、`tencent-pg-slowquery-diagnosis`

推荐一次性安装完整 bundle：`tencentdb-postgresql-skill-v1.0.3.zip`。

## 你可以用这些技能做什么

| 技能包 | 适合让它帮你做什么 | 可以直接这样说 | 使用边界 |
|---|---|---|---|
| `tencent-pg-management` | 查看实例状态、恢复时间窗、账号权限、SSL 配置，或先评估是否适合做变更 | `查看实例状态并评估是否适合升级规格`、`帮我看恢复时间窗`、`检查账号权限和 SSL 配置` | 会先返回只读事实；写类、费用类或高风险动作仍需你明确确认 |
| `tencent-pg-mem0-deploy` | 开通或关闭 mem0 服务 | `帮我给广州 postgres-abc12345 开通 mem0`、`帮我关闭广州 postgres-abc12345 的 mem0` | 只有目标唯一且前置检查通过时才会直接执行；密钥只从运行时环境读取 |
| `tencent-pg-rest-deploy` | 开通、关闭、调用或排查 REST / PostgREST 服务 | `帮我给广州 postgres-abc12345 开通 REST 服务`、`帮我调用这个 REST 路径`、`帮我排查 REST 502` | 会先做只读取证；安全组规则本身不支持直接改，只有更换绑定安全组时才会在确认后继续 |
| `tencent-pg-inspection` | 做实例健康巡检，查看固定指标集合 | `PG巡检`、`健康检查`、`资源巡检` | 当前保持只读，返回固定指标集合与保守风险提示 |
| `tencent-pg-slowquery-diagnosis` | 查看慢 SQL 基础事实 | `慢SQL查询`、`查看慢查询`、`慢SQL分析` | 当前保持只读，只返回慢 SQL 基础事实，不做根因和调优建议 |

## 如何选择安装方式

### 1. 只安装一个技能包

适合你的目标很明确，只想用一类能力，例如只安装：

- `tencent-pg-management-v1.0.3.zip`
- `tencent-pg-mem0-deploy-v1.0.3.zip`
- `tencent-pg-rest-deploy-v1.0.3.zip`
- `tencent-pg-inspection-v1.0.3.zip`
- `tencent-pg-slowquery-diagnosis-v1.0.3.zip`

### 2. 一次安装完整 bundle

适合你希望一次性获得完整 PostgreSQL 技能集合：

- `tencentdb-postgresql-skill-v1.0.3.zip`

如果你安装的是 bundle，通常可以先从根级 `SKILL.md` 进入，再按具体任务进入对应技能目录。

## 你实际需要关注的只有这 5 个技能包

当前对外发布和安装时，应以这 5 个技能包为准：

- `tencent-pg-management`
- `tencent-pg-mem0-deploy`
- `tencent-pg-rest-deploy`
- `tencent-pg-inspection`
- `tencent-pg-slowquery-diagnosis`

## 运行前准备

### 标准运行时变量

所有当前 skill 都优先依赖下面这组标准环境变量：

- `TENCENTCLOUD_SECRET_ID`
- `TENCENTCLOUD_SECRET_KEY`
- `TENCENTCLOUD_REGION`
- 可选：`TENCENTCLOUD_SESSION_TOKEN`

补充说明：

- **密钥只放运行时环境**，不要写进仓库文件、URL、查询参数或聊天记录
- **地域优先使用标准地域码**，例如 `ap-guangzhou`，也可以使用中文如“广州”
- 如果宿主进程使用自定义变量名，需要先映射到标准 `TENCENTCLOUD_*` 变量

### 如何把变量放到正确的运行环境

下面只保留一个原则：**让真正触发 skill 的那个进程能读到同一组 `TENCENTCLOUD_*` 变量**。

- **终端里直接跑**：在启动命令前 `export`
- **macOS 图形客户端**：优先用 `launchctl setenv`
- **Windows 客户端**：优先配置用户级环境变量后重启客户端
- **容器 / 自建宿主**：用运行时环境注入，而不是把密钥写进仓库配置

更完整的按客户端示例，请直接看：[USAGE_GUIDE.md](./USAGE_GUIDE.md)。

## 使用方式

### 1. 你可以直接怎么说

这套技能的目标就是：**尽量让你用一句自然语言就把任务说清楚**。

推荐示例：

- `查看实例状态并评估是否适合升级规格`
- `帮我看恢复时间窗`
- `帮我给广州 postgres-abc12345 开通 mem0`
- `帮我关闭广州 postgres-abc12345 的 mem0`
- `帮我给广州 postgres-abc12345 开通 REST 服务`
- `帮我关闭广州 postgres-abc12345 的 REST 服务`
- `帮我调用这个 REST 路径`
- `帮我排查 REST 502`
- `PG巡检`
- `慢SQL查询`

如果你能直接给出 `地域 + 实例 ID`，通常会减少补充追问。

## 本地开发与打包

### 开发入口

通常修改下面这些内容：

- 对应场景目录下的 `SKILL.md`
- 对应场景目录下的 `references/`
- 公共规则：`references/common/`
- 打包脚本：`scripts/package-skill.mjs`、`scripts/package-all.mjs`

### 推荐命令

```bash
cd src/postgres/skills
npm run verify
npm run package
npm run release
```

说明：

- `npm run verify`：校验 skill 结构
- `npm run package`：生成所有单独 zip 与 bundle zip
- `npm run release`：清理、校验并打包

## 相关文档

- [English README](./README_EN.md)
- [完整中文使用指南](./USAGE_GUIDE.md)
- [English Usage Guide](./USAGE_GUIDE_EN.md)
- `CURRENT_SKILL_REFERENCE.md`：当前 skill 现状快照，适合后续对话快速接手
- `SKILL_SECTION.md`：适合放到 bundle 根级 `SKILL.md` 的摘要结构
- `references/common/error_handling.md`：公共缺参 / 缺凭证 / 地域错误回复模板
- `references/common/region_normalization.md`：公共地域归一化规则
