# PostgreSQL Skill

PostgreSQL Skill is the unified skill entry for TencentDB for PostgreSQL. It is now organized around **five scenarios**: **management-plane entry**, **one-sentence mem0 deployment**, **one-sentence REST deployment**, **operations inspection**, and **operations slow SQL lookup**.

- `tencent-pg-management` handles the management-plane entry.
- `tencent-pg-mem0-deploy` and `tencent-pg-rest-deploy` handle one-sentence service enablement for expansion services.
- `tencent-pg-inspection` and `tencent-pg-slowquery-diagnosis` are operations-plane scenarios and currently stay read-only.

## 1. Use Cases and Access Links

| Scenario | Description |
|---|---|
| Management-plane Entry | Infer whether the request is about instance overview, instance changes, backup recovery, or access security, then route to the aligned OpenAPI actions |
| One-sentence mem0 Deployment | Auto-fill deployment slots, inspect current state first, then push mem0 toward a ready-to-use state |
| One-sentence REST Deployment | Auto-fill region and instance, inspect current state first; if the current instance region does not support REST, proactively return alternative region candidates before retrying REST / PostgREST |
| Operations Inspection | Review a fixed metric set from monitor APIs |
| Operations Slow SQL Lookup | Review basic slow SQL facts in a fixed time window |

### 1.2 Access Links

| Entry | Description | Link |
|---|---|---|
| GitHub Repository | View source code and documents | [Open GitHub Repository](https://github.com/TencentCloudCommunity/mcp-server/tree/feat/postgres-stdio-npx-support/src/postgres/skills) |
| ClawHub | View and distribute from marketplace | [Open ClawHub Page](https://clawhub.ai/tencent-adm/tencentdb-postgresql-skill) |
| SkillHub | View details from skill marketplace | [Open SkillHub Page](https://skillhub.cn/skills/tencentdb-postgresql-skill) |

## 2. Feature Overview

| Feature | Description | Output |
|---|---|---|
| Management-plane Entry | Covers instance status review, instance changes, backup recovery, and access-security governance, and routes to the smallest API set inside the skill | Current facts, recognized intent, extracted slots, inspected APIs, risks, safe next steps, and explicit waiting wording plus action meaning and main risk when confirmation is still needed |
| One-sentence mem0 Deployment | Covers target resolution, read-only preflight, `OpenMem0Service`, and polling until ready | Target scope, slot sources, current mem0 state, a ready-to-use address, and official acquisition links for missing slots |
| One-sentence REST Deployment | Covers target resolution, read-only preflight, alternative-region lookup when the current region is blocked, `OpenPostgRESTService` and polling until ready | Target scope, slot sources, current REST state, a ready-to-use access address, candidate regions when the current region is blocked, and explicit waiting wording plus action risk when user confirmation is still needed |
| Operations Inspection | Fixed monitor metrics for a target instance | Inspection summary, health snapshot, metric details, and manual review items |
| Operations Slow SQL Lookup | Read-only slow SQL lookup for a target instance | Query summary, key findings, detailed slow SQL list, and manual review items |

## 3. How to Use

Recommended input information:

```text
Region: for example ap-guangzhou
Instance ID: for example postgres-abc12345
Natural-language task: for example show instance status and evaluate whether a spec upgrade fits
Optional scope: for example last 1 hour / last 24 hours / specific database / target spec / public access / AgenticBaseId
```

## 4. Safety Boundaries

| Boundary | Description |
|---|---|
| Infer intent before calling APIs | `tencent-pg-management` classifies the management task first and then chooses the smallest aligned API set |
| Dedicated deployment skills inspect before opening | `tencent-pg-mem0-deploy` and `tencent-pg-rest-deploy` both run read-only checks before deciding whether to open the service |
| No blind writes under ambiguity or blockers | When the target is not unique, prerequisites are missing, or blockers remain, the skill does not execute write actions |
| Operations-plane scenarios currently stay read-only | `PG inspection` and `slow SQL lookup` do not execute management actions |
| Runtime-only secrets | Credentials must not be written into repository files or URLs |
