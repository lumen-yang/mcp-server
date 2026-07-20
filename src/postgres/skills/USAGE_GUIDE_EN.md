# PostgreSQL Skill Package — Usage Guide

> 中文版：[USAGE_GUIDE.md](./USAGE_GUIDE.md)

This guide explains how to install and use the current PostgreSQL skill packages. The package surface has now been consolidated into **five scenarios**:

- **Management-plane skill**: `tencent-pg-management`
- **Extension-service skills**: `tencent-pg-mem0-deploy`, `tencent-pg-rest-deploy`
- **Operations skills**: `tencent-pg-inspection`, `tencent-pg-slowquery-diagnosis`

In short:

- `tencent-pg-management` routes a natural-language request into one of four OpenAPI lanes: instance overview, instance change, backup and recovery, or access and security
- `tencent-pg-mem0-deploy` turns a one-sentence mem0 open / close request into an executable flow and returns the final service status
- `tencent-pg-rest-deploy` turns a one-sentence REST / PostgREST open, read-only call, or 502 troubleshooting request into an executable flow and returns the final status or troubleshooting conclusion
- The two operations skills focus on observation tasks and stay read-only for now

---

## 1. Before installation

### 1.1 Get your Tencent Cloud API credentials

1. Open [Tencent Cloud API Key Management](https://console.cloud.tencent.com/cam/capi)
2. If you do not have a usable key yet, click `Create Key`
3. Prefer a **least-privilege CAM sub-account** instead of reusing a long-lived high-privilege root account key
4. Save `SecretId` and `SecretKey` immediately after creation
   - `SecretId` can be viewed again later
   - `SecretKey` is usually shown in full only once; if you lose it, create a new one and rotate the old key
5. Open the [PostgreSQL console](https://console.cloud.tencent.com/postgres) and confirm the target region, preferably in standard form such as `ap-guangzhou`

> Never paste `SecretKey`, `SessionToken`, or any other sensitive value into chat, source code, screenshots, or tickets.

### 1.2 Set environment variables

Use these standard variable names:

- `TENCENTCLOUD_SECRET_ID`
- `TENCENTCLOUD_SECRET_KEY`
- `TENCENTCLOUD_REGION`
- add `TENCENTCLOUD_SESSION_TOKEN` only when you use temporary credentials

#### First decide how you are actually running the client

| How you use it | Typical clients / scenarios | Recommended method |
|---|---|---|
| You launch the host process, CLI, or debug command directly from Terminal | `node`, `npm run ...`, local scripts | Method A |
| You use a macOS GUI client launched from Finder / Dock / app icon | `WorkBuddy`, `CodeBuddy`, `Claude Desktop`, `Cherry Studio`, `Chatbox`, and similar apps | Method B |
| You use an IDE-style client | `Cursor`, `VS Code`, `Windsurf`, `Trae`, and similar tools | If launched from Terminal, use Method A; if launched from the GUI, use Method B; if you only run the host in the integrated terminal, Method A is also enough |
| You use a Windows desktop client | Windows desktop AI clients, IDEs, Electron apps | Method D |

#### Method A: CLI / command-line launch

Use this when you **start the host process, CLI, or local debug command directly from the command line**.

1. In the same terminal where you are about to launch the process, run:

```bash
export TENCENTCLOUD_SECRET_ID="your SecretId"
export TENCENTCLOUD_SECRET_KEY="your SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
# add this only when you use temporary credentials:
export TENCENTCLOUD_SESSION_TOKEN="your SessionToken"
```

2. Verify right away:

```bash
echo $TENCENTCLOUD_SECRET_ID
echo $TENCENTCLOUD_REGION
```

3. If the values are visible, continue with your launch command

> This fits **command-line launch only**. If you use a desktop client, a one-time `export` in some random terminal is usually not the main solution.

#### Method B: macOS GUI clients (`WorkBuddy`, `CodeBuddy`, `Claude Desktop`, `Cherry Studio`, `Chatbox`, and similar apps)

Use this when you work through a macOS desktop client instead of launching everything from CLI.

1. Open Terminal and run:

```bash
launchctl setenv TENCENTCLOUD_SECRET_ID "your SecretId"
launchctl setenv TENCENTCLOUD_SECRET_KEY "your SecretKey"
launchctl setenv TENCENTCLOUD_REGION "ap-guangzhou"
# add this only when you use temporary credentials:
launchctl setenv TENCENTCLOUD_SESSION_TOKEN "your SessionToken"
```

2. Verify that the variables are now present in the current login session:

```bash
launchctl getenv TENCENTCLOUD_SECRET_ID
launchctl getenv TENCENTCLOUD_REGION
```

3. **Fully quit** the client
4. Reopen the client and continue your task

> A GUI app launched from Dock or Finder often does **not** inherit the variables you exported in one terminal window. On macOS, `launchctl setenv` is usually the safer choice for desktop clients.

#### Method C: IDE / editor clients (`Cursor`, `VS Code`, `Windsurf`, `Trae`, and similar tools)

For these tools, first decide **whether the skill is triggered by the IDE main process, or whether you start the host yourself in the integrated terminal**.

##### Scenario 1: You start the command yourself in the integrated terminal

This is effectively still **Method A**. Set the variables in that terminal and run your command there.

```bash
export TENCENTCLOUD_SECRET_ID="your SecretId"
export TENCENTCLOUD_SECRET_KEY="your SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
```

##### Scenario 2: You want the IDE main process and its extensions to see the variables too

If you launch the IDE from Terminal, you can do:

```bash
export TENCENTCLOUD_SECRET_ID="your SecretId"
export TENCENTCLOUD_SECRET_KEY="your SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"

# choose the one you actually use
cursor .
# code .
# windsurf .
# trae .
```

If you normally launch the IDE from Dock / Finder / an app icon, go back to **Method B** and inject the variables into the macOS login session first.

> Simple rule: **give the variables to the process that actually launches the work**. If the integrated terminal starts the host, give them to that shell. If the IDE main process or plugin needs them, give them to the IDE process.

#### Method D: Windows clients

Use this when you work on Windows with a desktop client, IDE, or custom host.

##### Option 1: Set user-level environment variables (recommended)

In PowerShell, run:

```powershell
setx TENCENTCLOUD_SECRET_ID "your SecretId"
setx TENCENTCLOUD_SECRET_KEY "your SecretKey"
setx TENCENTCLOUD_REGION "ap-guangzhou"
# temporary session tokens are better injected per session instead of stored long-term with setx
```

Then:

1. Close and reopen your terminal or client
2. Verify with:

```powershell
echo $env:TENCENTCLOUD_SECRET_ID
echo $env:TENCENTCLOUD_REGION
```

##### Option 2: Inject only into the current session

```powershell
$env:TENCENTCLOUD_SECRET_ID = "your SecretId"
$env:TENCENTCLOUD_SECRET_KEY = "your SecretKey"
$env:TENCENTCLOUD_REGION = "ap-guangzhou"
$env:TENCENTCLOUD_SESSION_TOKEN = "your SessionToken"
```

This is more suitable for temporary debugging or one session only.

> `setx` does not update the already-open window automatically. Restart the terminal or client after changing it.

#### Other host setups

For Docker, CI, self-hosted runtimes, or hosts that already use custom variable names, this guide no longer keeps separate step-by-step methods.

Use one consistent rule instead: **before the real host process starts, inject the standard `TENCENTCLOUD_*` variables into the exact process environment that triggers the skill**. If your environment already uses different names internally, map them into the standard variables in your startup chain first.

#### Common mistakes

- **Mistake 1**: “I exported the variables in one terminal, so the desktop client must see them too.”  
  Not always. Many GUI apps launched from an icon do not inherit that terminal environment.
- **Mistake 2**: “I can just write `SecretKey` into a client config file or repo config.”  
  Do not do that. Prefer system environment variables, platform secrets, or runtime injection.
- **Mistake 3**: “I do not need to restart the client after `setx` or `launchctl setenv`.”  
  Many clients only reload environment variables on restart.
- **Mistake 4**: “A temporary `SessionToken` is effectively permanent.”  
  It expires and must be refreshed.

---

## 2. Install the skills

### 2.1 Download

Download the following zip files from the release:

- `tencent-pg-management-v1.0.3.zip`
- `tencent-pg-mem0-deploy-v1.0.3.zip`
- `tencent-pg-rest-deploy-v1.0.3.zip`
- `tencent-pg-inspection-v1.0.3.zip`
- `tencent-pg-slowquery-diagnosis-v1.0.3.zip`
- or `tencentdb-postgresql-skill-v1.0.3.zip`

### 2.2 Import into CodeBuddy / WorkBuddy

1. Open skill management
2. Import the target `.zip`
3. Enable the skill

### 2.3 Verify the installation

After installation, you can test with prompts such as:

- `show instance status and evaluate whether a spec upgrade fits` → should trigger `tencent-pg-management`
- `open mem0 for ap-guangzhou postgres-abc12345` → should trigger `tencent-pg-mem0-deploy`
- `close mem0 for ap-guangzhou postgres-abc12345` → should trigger `tencent-pg-mem0-deploy`
- `open REST service for ap-guangzhou postgres-abc12345` → should trigger `tencent-pg-rest-deploy`
- `close REST service for ap-guangzhou postgres-abc12345` → should trigger `tencent-pg-rest-deploy`
- `call this REST path for me` → should trigger `tencent-pg-rest-deploy`
- `troubleshoot REST 502` → should trigger `tencent-pg-rest-deploy`
- `PG inspection` → should trigger `tencent-pg-inspection`
- `slow SQL lookup` → should trigger `tencent-pg-slowquery-diagnosis`

---

## 3. Unified management entry

### 3.1 When to use it

Use this skill when you want to launch management tasks directly in natural language, such as:

- checking instance status, spec, version, task state, and read-only group state
- evaluating or preparing restart, scale up/down, isolate, de-isolate, create instance, or create read-only instance
- checking backup overview, recovery window, or evaluating backup creation, backup download, or clone restore
- checking account privileges, database owner, public access, security groups, and SSL settings

### 3.2 Example prompts

- `show instance status for ap-guangzhou postgres-abc12345`
- `help me evaluate whether ap-guangzhou postgres-abc12345 is suitable for a spec upgrade`
- `check the recovery window for ap-shanghai postgres-abc12345`
- `review account privileges and SSL for ap-beijing postgres-abc12345`
- `before resetting a password on postgres-abc12345, show me the current privileges first`

### 3.3 Behavior

This skill will:

1. extract the main goal and key slots such as region, instance ID, account, database, recovery window, target spec, and more
2. normalize the region
3. route the request into one of the four main lanes: instance overview, instance change, backup and recovery, or access and security
4. build the minimal API plan for the chosen lane and start with read-only evidence collection
5. return the current facts, blockers, and feasibility first, before deciding whether a confirmation step is needed
6. mark write, fee-impacting, or high-risk actions as `pending confirmation`; if a follow-up confirmation is still required, the reply will clearly say it is waiting for confirmation and explain the meaning and main risk of the next action

This skill will **not**:

- force the user to manually choose one internal sub-scenario first
- call APIs blindly when required scope information is missing
- execute restart, scaling, database creation, password reset, or public-access changes without explicit confirmation
- force operations inspection or slow-query observation requests into the management execution lane

### 3.4 Output format

```text
1. Scope: region / instance / optional object
2. Intent: instance overview / instance change / backup and recovery / access and security
3. Extracted slots: key fields that affect routing and API choice
4. APIs used: OpenAPI calls used for evidence collection or execution
5. Current facts / result: structured conclusions directly related to the task
6. Risks and confirmations: pending actions, fee impact, security impact; if a follow-up confirmation is still needed, the reply should clearly say it is waiting for confirmation, what will happen after confirmation, and the meaning and main risk of that action
7. Suggested next step
```

### 3.5 Related specialized skill: one-sentence mem0 open / close

This is suitable for very direct requests such as:

- `open mem0 for ap-guangzhou postgres-abc12345`
- `deploy mem0 to the current instance with AgenticBaseId ab-xxxxx`
- `close mem0 for ap-guangzhou postgres-abc12345`

This skill will:

1. identify whether the request is an `open` or `close` action and fill the needed slots, including `AgenticBaseId`, `LLMModel`, and so on
2. prefer defaults and secrets from runtime, instead of asking the user to paste them into chat
3. when critical parameters are missing, return direct console / product-site links, the shortest click path, runtime configuration examples, and an exact “what to say next” prompt; for non-secret fields such as `Region` or instance ID it may also offer “I can look this up for you”, while secret values such as `EmbeddingApiKey` still only get official entry links and runtime-injection guidance
4. run `DescribeDBInstanceAttribute` / `DescribeMem0Service` as prechecks
5. execute `OpenMem0Service` or `CloseMem0Service` according to the target direction
6. poll service status until it reaches a usable mem0 address, a confirmed closed state, or a clear blocking reason

The goal of this skill is not to hand out a plan, but to **push mem0 into a usable or final closed state in as few turns as possible**.

#### Parameter lookup entry points

- `AgenticBaseId`: open the PostgreSQL console at `https://console.cloud.tencent.com/postgres`
  - click path: `AI 应用` → choose region → `AgenticBase`
  - if a Base already exists, open the detail page and copy `AgenticBaseId`
  - otherwise click `新建`
- `EmbeddingApiKey`: open the Hunyuan console at `https://console.cloud.tencent.com/hunyuan`
  - click path: `立即接入管理` → `API Key 管理` → `创建 API KEY`
  - put the created key into runtime only, never into chat
- `LLMModel`: if you have no special requirement, use the default `auto`
  - only if you want to pin a specific model should you open the Hunyuan console and check the currently available model list

### 3.6 Related specialized skill: one-sentence REST open / close

This is suitable for very direct requests such as:

- `open REST service for ap-guangzhou postgres-abc12345`
- `deploy PostgREST to the current instance`
- `close REST service for ap-guangzhou postgres-abc12345`

This skill will:

1. identify whether the request is an `open` or `close` action and fill region and instance ID, preferring runtime defaults first
2. run `DescribeDBInstanceAttribute` / `DescribePostgRESTService` as read-only checks first; if REST is unsupported in the region, it returns candidate regions you can try next
3. for open requests, execute `OpenPostgRESTService(DBInstanceId, EnableWanNet=false)` by default — that means **enable REST service first while keeping public / external access closed by default**; for close requests, execute `ClosePostgRESTService`
4. poll service status until REST is confirmed usable, the service is confirmed closed, or a clear blocker is returned; if the open action used the default no-public mode, explicitly say that the service is enabled while public access stays closed
5. output the target scope, slot sources, executed APIs, the effective `EnableWanNet` value, current status, candidate regions when unsupported, and suggested next usage; if the user explicitly asks the skill to open public / external exposure, stop in a waiting-for-confirmation state first and explain the meaning and risk of that action

The goal of this skill is likewise not to hand out a manual, but to **push REST / PostgREST into a usable or final closed state in as few turns as possible**.

---

## 4. Operations PG inspection

### 4.1 When to use it

Use this when you want a basic inspection result for one instance, such as:

- CPU
- memory
- storage
- connections
- I/O
- replication lag

### 4.2 Example prompts

- `PG inspection ap-guangzhou postgres-abc12345`
- `health check for ap-guangzhou postgres-abc12345 in the last 1 hour`
- `resource inspection for ap-shanghai postgres-abc12345`

### 4.3 Behavior

This skill will:

1. confirm region, instance ID, and optional time window
2. normalize the region
3. discover supported metrics through Tencent Cloud Monitor APIs
4. fetch a fixed metric set
5. return the basic inspection result

This skill will **not**:

- expand into backup, account, parameter, security-group, or SSL checks
- enter a troubleshooting workflow
- output repair actions
- generate automatic root-cause analysis

### 4.4 Output format

```text
1. Inspection summary: overall state / key findings / metric coverage
2. Inspection target: region / instance / time range / monitor APIs used
3. Health snapshot: summarized CPU / memory / storage / connection / I/O / replication-lag state
4. Metric details: metric name / value / unit / data state / factual note
5. Risks and manual review items: attention / abnormal / manual review needed
6. Data notes: unsupported / no-data metrics and time-window notes
```

---

## 5. Operations slow-query lookup

### 5.1 When to use it

Use this when you want basic slow-query information for one instance in a given time window, such as:

- Top N slow queries
- SQL text or abstracted SQL
- database name
- username
- client address
- execution time
- execution count
- total time

### 5.2 Example prompts

- `slow SQL lookup ap-guangzhou postgres-abc12345 in the last 1 hour`
- `show slow queries for ap-guangzhou postgres-abc12345`
- `slow SQL analysis for ap-shanghai postgres-abc12345 from 3 PM to 5 PM yesterday`

### 5.3 Behavior

This skill will:

1. confirm region, instance ID, and time window
2. normalize the region
3. call `DescribeSlowQueryList`
4. call `DescribeSlowQueryAnalysis` when aggregation or ranking is needed
5. return basic slow-query facts

This skill will **not**:

- rank root causes
- perform true root-cause analysis
- infer from correlated error logs
- output tuning, scaling, or parameter-change advice

### 5.4 Output format

```text
1. Query scope: region / instance / time window / sort mode
2. Slow-query list: SQL / database / user / latency / count / total latency / timestamp
```

---

## 6. FAQ

### Q1: What if the reply says credentials are missing?

Confirm that the environment variables are present:

```bash
echo $TENCENTCLOUD_SECRET_ID
echo $TENCENTCLOUD_SECRET_KEY
echo $TENCENTCLOUD_REGION
```

If you use a desktop client, confirm one more thing: **the process that actually triggers the skill** can read them, not just a temporary terminal window.

### Q2: Are Chinese region names supported?

Yes. Common Chinese names can be normalized, for example:

- `广州` → `ap-guangzhou`
- `上海` → `ap-shanghai`
- `成都` → `ap-chengdu`
- `北京` → `ap-beijing`

### Q3: Why did it not directly execute a change action?

That is by design for the current management skill:

- pick the right lane and API based on the real task
- return read-only facts first
- require explicit confirmation for write, fee-impacting, or high-risk actions

### Q4: Why are PG inspection and slow-query lookup still separate skills?

That is the current layering model:

- `tencent-pg-management` handles natural-language routing and controlled execution for management tasks
- `tencent-pg-inspection` and `tencent-pg-slowquery-diagnosis` belong to the operations layer and return stable inspection / slow-query observations

This reduces misrouting and improves result stability.

---

## 7. Notes

1. On the first request, it is best to provide `region + instance ID`
2. The goal of `tencent-pg-management` is **one-sentence access to management APIs**, not a full troubleshooting workflow
3. If you need more advanced capabilities, extend the current skill set with direct OpenAPI / Monitor actions instead of adding extra prerequisite services
4. If you need the Chinese version, open [USAGE_GUIDE.md](./USAGE_GUIDE.md)
