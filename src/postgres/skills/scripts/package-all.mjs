import { cpSync, existsSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { execFileSync } from 'node:child_process';
import { COMMON_REFERENCES_DIR, DIST_DIR, SKILLS, SKILLS_ROOT, packageSkill, resolveVersion } from './package-skill.mjs';

const BUNDLE_SLUG = 'tencentdb-postgresql-skill';
const BUNDLE_DISPLAY_NAME = 'TencentDB PostgreSQL Skill';
const BUNDLE_DESCRIPTION = 'Bundle entry for TencentDB PostgreSQL skills. These skills call Tencent Cloud PostgreSQL OpenAPI directly and can use the full 48 aligned actions documented by each child skill under references/. Explicit confirmation is required before any write, fee-impacting, or high-risk action.';
const BUNDLE_SKILL_GUIDE = {
  'tencent-pg-inspection': {
    displayName: 'PG 巡检',
    scenarios: '日常健康检查、备份检查、资源与配置巡检',
    exampleRequests: '`PG巡检`、`健康检查`、`备份检查`、`资源水位检查`',
    focus: '实例健康、备份状态、网络 / SSL / 只读上下文、参数姿态、风险总结',
    outputs: '健康状态、风险项、重点关注项、建议动作',
  },
  'tencent-pg-slowquery-diagnosis': {
    displayName: '慢 SQL 诊断',
    scenarios: '慢 SQL 分析、延迟波动、性能退化定位',
    exampleRequests: '`慢SQL分析`、`SQL性能诊断`、`查询为什么变慢`',
    focus: '慢查询证据、错误日志、实例上下文、可能原因排序、安全优化建议',
    outputs: '慢 SQL 排名、原因分析、风险提示、优化建议',
  },
  'tencent-pg-ops-troubleshooter': {
    displayName: '运维排障',
    scenarios: '连接异常、备份失败、SSL 问题、实例异常排查',
    exampleRequests: '`PG排障`、`实例异常排查`、`SSL问题`、`备份失败`',
    focus: '故障分类、分模块证据采集、Runbook 风格处置建议',
    outputs: '问题线索、阻塞点、下一步动作、需确认的高风险操作',
  },
};

function ensureDistDir() {
  rmSync(DIST_DIR, { recursive: true, force: true });
  mkdirSync(DIST_DIR, { recursive: true });
}

function bundleFilter(src) {
  const name = path.basename(src);
  if (name.startsWith('.')) {
    return false;
  }
  if (name === 'dist') {
    return false;
  }
  return true;
}

function parseFrontmatter(markdown) {
  const match = markdown.match(/^---\n([\s\S]*?)\n---\n?/);
  if (!match) {
    return {};
  }

  const metadata = {};
  for (const line of match[1].split('\n')) {
    const separatorIndex = line.indexOf(':');
    if (separatorIndex === -1) {
      continue;
    }

    const key = line.slice(0, separatorIndex).trim();
    let value = line.slice(separatorIndex + 1).trim();
    if (!key) {
      continue;
    }

    if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
      value = value.slice(1, -1);
    }

    metadata[key] = value;
  }

  return metadata;
}

function readSkillDescriptors() {
  return SKILLS.map((skillName) => {
    const skillFile = path.join(SKILLS_ROOT, skillName, 'SKILL.md');
    const content = readFileSync(skillFile, 'utf8');
    const meta = parseFrontmatter(content);

    return {
      skillName,
      name: meta.name || skillName,
      description: meta.description || '',
      descriptionZh: meta.description_zh || meta.description_en || meta.name || skillName,
      entryPath: `references/${skillName}/SKILL.md`,
      directoryPath: `references/${skillName}/`,
    };
  });
}

function extractReadmeSection(markdown, startHeading, endHeading) {
  const startIndex = markdown.indexOf(startHeading);
  if (startIndex === -1) {
    return '';
  }

  const endIndex = markdown.indexOf(endHeading, startIndex + startHeading.length);
  const section = endIndex === -1 ? markdown.slice(startIndex) : markdown.slice(startIndex, endIndex);
  return section.trim();
}

function renderBundleSkill(version, descriptors) {
  const sections = descriptors
    .map((descriptor) => {
      const guide = BUNDLE_SKILL_GUIDE[descriptor.skillName] || {};
      const displayName = guide.displayName || descriptor.descriptionZh || descriptor.skillName;
      const scenarios = guide.scenarios || descriptor.descriptionZh || '请进入该目录查看完整 skill 说明。';
      const exampleRequests = guide.exampleRequests || '请进入对应目录查看示例请求。';
      const outputs = guide.outputs || '结构化结论与建议';
      return `| \`${descriptor.skillName}\` | ${displayName} | ${scenarios} | ${exampleRequests} | ${outputs} | \`${descriptor.entryPath}\` |`;
    })
    .join('\n');
  const skillsReadme = readFileSync(path.join(SKILLS_ROOT, 'README.md'), 'utf8');
  const alignedActionsSection = extractReadmeSection(skillsReadme, '## 当前开放的 48 个对齐 Action', '## 目录结构');
  const appendixSection = alignedActionsSection
    ? alignedActionsSection.replace('## 当前开放的 48 个对齐 Action', '## 附录：当前开放的 48 个对齐 Action')
    : '';

  const lines = [
    '---',
    `name: "${BUNDLE_DISPLAY_NAME}"`,
    `description: "${BUNDLE_DESCRIPTION}"`,
    `description_zh: "${BUNDLE_DISPLAY_NAME}"`,
    `description_en: "${BUNDLE_DISPLAY_NAME}"`,
    `version: ${version}`,
    '---',
    '',
    `# ${BUNDLE_DISPLAY_NAME}`,
    '',
    '## 简介',
    '',
    `${BUNDLE_DISPLAY_NAME} 是面向腾讯云数据库 PostgreSQL 的任务型技能包入口，统一覆盖 **PG 巡检**、**慢 SQL 诊断** 和 **运维排障** 三类高频场景。它直接调用腾讯云 PostgreSQL OpenAPI，不依赖已部署的 MCP Server，并且严格限制在当前仓库已经完成参数对齐的 **48 个 Action** 范围内执行。`,
    '',
    '这个总入口适合用来快速理解整个技能包能做什么、如何使用，以及应该进入哪个子 skill；如果你已经知道自己的问题类型，可以直接进入 `references/` 下对应子目录继续查看更细的执行说明。',
    '',
    '## 适用场景',
    '',
    '- 需要做 PostgreSQL **日常巡检、发布后回归或健康检查**。',
    '- 需要定位 **慢 SQL、性能抖动、查询变慢** 等问题。',
    '- 需要处理 **连接失败、备份异常、账号权限、网络 / SSL、只读链路** 等运维问题。',
    '- 希望由 AI 基于证据先做分析，再在必要时进入对齐后的 OpenAPI 处置动作。',
    '',
    '## 包含的技能',
    '',
    '| Skill 目录 | 能力名称 | 适合处理的问题 | 典型请求 | 输出结果 | 入口 |',
    '|---|---|---|---|---|---|',
    sections,
    '',
    '## 使用前准备',
    '',
    '推荐用户优先准备以下运行时环境变量：',
    '',
    '- `TENCENTCLOUD_SECRET_ID`',
    '- `TENCENTCLOUD_SECRET_KEY`',
    '- `TENCENTCLOUD_REGION`',
    '- 可选：临时凭证场景补充 `TENCENTCLOUD_SESSION_TOKEN`',
    '',
    '补充说明：',
    '',
    '- 兼容读取 `MCP_REQUEST_SECRET_ID`、`MCP_REQUEST_SECRET_KEY`、`MCP_REQUEST_SESSION_TOKEN`、`MCP_SECRET_ID`、`MCP_SECRET_KEY`。',
    '- 地域既可以直接写 `ap-guangzhou`，也可以先写 `广州`、`上海`、`成都`、`北京` 等常见别名，执行前应先归一化为标准地域码。',
    '- 建议先通过腾讯云 [API 密钥管理控制台](https://console.cloud.tencent.com/cam/capi) 检查可用凭证，通过 [地域与域名说明](https://cloud.tencent.com/document/product/1596/77930) 确认目标实例所在地域。',
    '- 优先使用官方 SDK；如果 SDK 不可用，不应把它当成首次使用的硬阻塞，而应退回到本地生成 TC3 签名 HTTPS 请求。',
    '- 密钥只应存在于运行时环境或安全上下文中，不要写入代码、仓库文件、URL 或查询参数。',
    '',
    '## 使用方法',
    '',
    '1. **确认目标范围**：至少准备好地域和实例 ID；如果还不知道实例 ID，可先进入控制台 `https://console.cloud.tencent.com/postgres` 查询。',
    '2. **选择匹配的子 skill**：巡检类问题进入 `tencent-pg-inspection`，性能类问题进入 `tencent-pg-slowquery-diagnosis`，故障排查类问题进入 `tencent-pg-ops-troubleshooter`。',
    '3. **发起请求**：建议在请求里明确地域、实例 ID、问题类型，以及可选的时间范围。',
    '4. **查看结果**：先阅读结论和证据摘要，再决定是否需要执行后续动作。',
    '',
    '推荐输入格式示例：',
    '',
    '```text',
    'ap-guangzhou postgres-abc12345',
    '',
    '请做一次 PG 巡检',
    '请分析最近 1 小时慢 SQL',
    '请排查这个实例为什么 SSL 连接异常',
    '```',
    '',
    '## 结果说明',
    '',
    '执行完成后，通常会得到以下几类输出：',
    '',
    '- **结论概览**：当前实例或问题的整体状态判断。',
    '- **证据摘要**：来自 OpenAPI 的关键证据，例如实例状态、备份、慢 SQL、错误日志、网络或 SSL 信息。',
    '- **风险与原因**：风险项、疑似原因排序或阻塞点说明。',
    '- **建议动作**：安全下一步、建议复核项，以及需要明确确认后才能执行的处置动作。',
    '',
    '## 安全与边界',
    '',
    '- 默认应以 **只读证据采集** 为先，不直接执行会影响实例状态的动作。',
    '- 任何 **写类、费用类或高风险动作** 都必须先说明影响面并获得明确确认。',
    '- 只允许使用当前仓库中已完成参数对齐的 **48 个 PostgreSQL OpenAPI Action**；未对齐的动作不应调用。',
    '- 不应混用不同地域、不同实例或未经确认的目标范围。',
    '',
    '## 使用限制与确认机制',
    '',
    '为确保分析结果准确且操作安全，使用该 skill 时默认遵循以下约束：',
    '',
    '- **先确认目标范围**：建议在请求中提供地域与实例 ID；如果信息不足，skill 会先提示补充必要信息，再继续分析。',
    '- **自动匹配对应能力**：系统会根据你的问题类型自动选择巡检、慢 SQL 诊断或运维排障能力，并进入对应说明继续执行。',
    '- **默认先做只读分析**：会优先采集实例状态、备份、慢 SQL、错误日志、网络或 SSL 等证据，再给出结论与建议。',
    '- **高风险动作必须确认**：涉及写类、费用类或可能影响实例状态的动作时，会先说明影响面、目标实例与预期结果，获得明确确认后才继续。',
    '- **仅在支持范围内调用**：只会使用当前已完成参数对齐的 48 个 PostgreSQL OpenAPI Action；如果凭证缺失、地域不合法或能力未开放，会先提示修正方式。',
    '',
    '## 目录与扩展说明',
    '',
    '- 根目录入口：`SKILL.md`',
    '- 根目录元数据：`_meta.json`',
    '- 公共规则目录：`references/common/`',
    '- 子技能目录：`references/<skill-name>/`',
    '- 子技能入口：`references/<skill-name>/SKILL.md`',
    '- 公共规则文档：`@references/common/region_normalization.md`、`@references/common/error_handling.md`',
    '',
  ];

  if (appendixSection) {
    lines.push(appendixSection, '');
  }

  return `${lines.join('\n')}\n`;
}

function renderBundleMeta(version) {
  return `${JSON.stringify(
    {
      slug: BUNDLE_SLUG,
      version,
      publishedAt: Date.now(),
    },
    null,
    2,
  )}\n`;
}

function stageExpandedSkills(stageDir) {
  const referencesDir = path.join(stageDir, 'references');
  mkdirSync(referencesDir, { recursive: true });

  if (existsSync(COMMON_REFERENCES_DIR)) {
    cpSync(COMMON_REFERENCES_DIR, path.join(referencesDir, 'common'), {
      recursive: true,
    });
  }

  for (const skillName of SKILLS) {
    const sourceDir = path.join(SKILLS_ROOT, skillName);
    const targetDir = path.join(referencesDir, skillName);
    cpSync(sourceDir, targetDir, {
      recursive: true,
      filter: bundleFilter,
    });
  }
}

function createBundle(version) {
  const stageDir = mkdtempSync(path.join(os.tmpdir(), `${BUNDLE_SLUG}-`));
  const bundlePath = path.join(DIST_DIR, `${BUNDLE_SLUG}-v${version}.zip`);
  const descriptors = readSkillDescriptors();

  try {
    writeFileSync(path.join(stageDir, 'SKILL.md'), renderBundleSkill(version, descriptors), 'utf8');
    writeFileSync(path.join(stageDir, '_meta.json'), renderBundleMeta(version), 'utf8');

    stageExpandedSkills(stageDir);

    rmSync(bundlePath, { force: true });
    execFileSync('zip', ['-rq', bundlePath, 'SKILL.md', '_meta.json', 'references'], {
      cwd: stageDir,
      stdio: 'inherit',
    });

    return bundlePath;
  } finally {
    rmSync(stageDir, { recursive: true, force: true });
  }
}

function createSinglePackages(version) {
  return SKILLS.map((skillName) => packageSkill(skillName, version, DIST_DIR));
}

function main() {
  const versionArg = process.argv.slice(2).find((token) => !token.startsWith('--'));
  const version = resolveVersion(versionArg);

  ensureDistDir();

  const packagedFiles = createSinglePackages(version);
  const bundlePath = createBundle(version);

  const output = [
    ...packagedFiles.map((filePath) => `- ${path.basename(filePath)}`),
    `- ${path.basename(bundlePath)}`,
  ].join('\n');

  console.log(`generated skill release assets (v${version}):\n${output}`);
}

main();
