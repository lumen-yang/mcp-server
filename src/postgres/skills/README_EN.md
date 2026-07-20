# PostgreSQL Skill Package Guide

> 中文版：[README.md](./README.md)

If you want to use TencentDB for PostgreSQL skills directly in CodeBuddy, WorkBuddy, or similar AI clients, this README is the fastest way to get started.

The current public skill surface has been consolidated into **five installable skills**:

- **Unified management entry**: `tencent-pg-management`
- **Extension-service skills**: `tencent-pg-mem0-deploy`, `tencent-pg-rest-deploy`
- **Operations and observation skills**: `tencent-pg-inspection`, `tencent-pg-slowquery-diagnosis`

The recommended one-shot install is the full bundle: `tencentdb-postgresql-skill-v1.0.3.zip`.

## What these skills can help you do

| Skill package | Best for | Example requests | Boundary |
|---|---|---|---|
| `tencent-pg-management` | Check instance status, recovery window, account privileges, SSL settings, or evaluate whether a change is appropriate | `show instance status and evaluate whether a spec upgrade fits`, `check the recovery window`, `review account privileges and SSL` | Returns read-only facts first; write, fee-impacting, or high-risk actions still require your explicit confirmation |
| `tencent-pg-mem0-deploy` | Open or close mem0 service | `open mem0 for ap-guangzhou postgres-abc12345`, `close mem0 for ap-guangzhou postgres-abc12345` | Executes directly only when the target is unambiguous and preflight passes; secrets stay in runtime only |
| `tencent-pg-rest-deploy` | Open, close, call, or troubleshoot REST / PostgREST service | `open REST service for ap-guangzhou postgres-abc12345`, `call this REST path for me`, `troubleshoot REST 502` | Starts from read-only evidence first; security-group rules themselves are not edited directly, and replacing the bound security-group set still needs explicit confirmation |
| `tencent-pg-inspection` | Run a PostgreSQL health inspection with a fixed metric set | `PG inspection`, `health check`, `resource inspection` | Read-only only; returns fixed metrics and conservative risk prompts |
| `tencent-pg-slowquery-diagnosis` | Look up slow-query facts | `slow SQL lookup`, `show slow queries`, `slow SQL analysis` | Read-only only; returns basic facts without root-cause analysis or tuning advice |

## How to choose an installation mode

### 1. Install one skill only

Use this when your goal is narrow and you only want one capability, for example:

- `tencent-pg-management-v1.0.3.zip`
- `tencent-pg-mem0-deploy-v1.0.3.zip`
- `tencent-pg-rest-deploy-v1.0.3.zip`
- `tencent-pg-inspection-v1.0.3.zip`
- `tencent-pg-slowquery-diagnosis-v1.0.3.zip`

### 2. Install the full bundle once

Use this when you want the complete PostgreSQL skill set in one import:

- `tencentdb-postgresql-skill-v1.0.3.zip`

If you install the bundle, the usual starting point is the root `SKILL.md`, then the matching skill directory for your task.

## The only five skill packages most users need

For public installation and release, these are the real entry points:

- `tencent-pg-management`
- `tencent-pg-mem0-deploy`
- `tencent-pg-rest-deploy`
- `tencent-pg-inspection`
- `tencent-pg-slowquery-diagnosis`

## Runtime prerequisites

### Standard runtime variables

All current skills depend first on this standard runtime variable set:

- `TENCENTCLOUD_SECRET_ID`
- `TENCENTCLOUD_SECRET_KEY`
- `TENCENTCLOUD_REGION`
- optional: `TENCENTCLOUD_SESSION_TOKEN`

Additional notes:

- **Secrets must stay in runtime only** and should never be written into repository files, URLs, query strings, or chat history
- **Prefer a standard region code** such as `ap-guangzhou`, although common Chinese city names can also be normalized
- If the host process uses custom variable names, map them into the standard `TENCENTCLOUD_*` names first

### Put the variables into the correct runtime

Keep one rule in mind: **the actual process that triggers the skill must be able to read the same `TENCENTCLOUD_*` variables**.

- **Running directly from Terminal**: use `export` before launch
- **macOS GUI clients**: prefer `launchctl setenv`
- **Windows desktop clients**: set user-level environment variables, then restart the client
- **Containers or self-hosted runtimes**: inject variables at runtime instead of storing secrets in repository config files

For step-by-step examples by client type, see [USAGE_GUIDE.md](./USAGE_GUIDE.md) and [USAGE_GUIDE_EN.md](./USAGE_GUIDE_EN.md).

## Usage

### 1. What you can say directly

The goal of this skill set is simple: **you should be able to describe the task in one natural-language sentence whenever possible**.

Recommended prompts:

- `show instance status and evaluate whether a spec upgrade fits`
- `check the recovery window`
- `open mem0 for ap-guangzhou postgres-abc12345`
- `close mem0 for ap-guangzhou postgres-abc12345`
- `open REST service for ap-guangzhou postgres-abc12345`
- `close REST service for ap-guangzhou postgres-abc12345`
- `call this REST path for me`
- `troubleshoot REST 502`
- `PG inspection`
- `slow SQL lookup`

If you can provide `region + instance ID` directly, the skill usually needs fewer follow-up turns.

## Local development and packaging

### Main edit points

Most changes happen in:

- `SKILL.md` under the target scenario directory
- `references/` under the target scenario directory
- shared rules in `references/common/`
- packaging scripts: `scripts/package-skill.mjs` and `scripts/package-all.mjs`

### Recommended commands

```bash
cd src/postgres/skills
npm run verify
npm run package
npm run release
```

What they do:

- `npm run verify`: validate skill structure
- `npm run package`: generate all individual zip packages and the bundle zip
- `npm run release`: clean, verify, and package

## Related documents

- [中文版 README](./README.md)
- [中文使用指南](./USAGE_GUIDE.md)
- [English Usage Guide](./USAGE_GUIDE_EN.md)
- `CURRENT_SKILL_REFERENCE.md`: a snapshot of the current skill state for future handoff and reuse
- `SKILL_SECTION.md`: the compact summary structure used for bundle-style root entries
- `references/common/error_handling.md`: shared missing-parameter, missing-credential, and region-error reply templates
- `references/common/region_normalization.md`: shared region normalization rules
