---
name: spec-driven
description: 贯穿功能规划(③)、实现+测试(④)、部署+发布(⑤)、复盘+archive(⑥)的 OpenSpec 生命周期引导技能。提案产出 proposal/design/tasks，apply 推进实现，archive 回写 specs 与 PROJECT.md。
---

# spec-driven

你是工作室的**生命周期引导技能**。OpenSpec 不是某一步——它贯穿整个功能开发：从③需求细化到⑥复盘。
每个功能是一个 OpenSpec change。详见 `rules/lifecycle.md`。

## 使用时机

- 项目 setup 完成（bootstrap 结束），要开始开发某个具体功能。
- 用户说"做 X 功能""加个 Y""迭代 Z"。

## 前置

- **OpenSpec 已安装并初始化。** 检查 `openspec/AGENTS.md` 与 `/opsx:*` 命令。
- 已有 `PROJECT.md` 和 `rules/tech-selection.md`（bootstrap 产出）。
- 用户已确认进入功能规划阶段（③）。

## 角色分工（agent 驱动，按阶段）

| 阶段 | Agent | 产出 |
|------|-------|------|
| ③a 需求细化 | project-manager | proposal.md（为什么做、要什么、验收标准 GIVEN/WHEN/THEN）|
| ③b UI 设计 | ui-designer | 线框图/高保真/状态矩阵（与③c并行） |
| ③c 技术架构 | tech-lead | design.md（技术方案、API 契约、数据模型）|
| ③d 任务拆分 | tech-lead | tasks.md（标注依赖与可并行项）|
| ④a 实现 | backend-dev / frontend-dev | 代码（并行开发）|
| ④b 审查 | code-reviewer | 每 PR 审查通过 |
| ④c 测试 | test-engineer | 按 spec 持续产出测试 |
| ④d UAT | project-manager | 验收确认 |
| ⑤ 部署 | devops | deploy/ 物料 |
| ⑥ 复盘 | tech-lead + pm | archive + 回顾 |

## OpenSpec 循环

```
/opsx:propose <name>   → 分阶段产出（③a→③b→③c→③d），每步有门
/opsx:apply            → 按 tasks 实现 + 测试 + UAT（④）
/opsx:archive          → 合并 delta 到 openspec/specs + 回写 PROJECT.md + 回顾（⑥）
```

## 执行流程

### 阶段 ③a：需求细化（Propose → proposal.md）

- project-manager 与用户澄清功能的**边界与验收标准**。
- 产出 `openspec/changes/<name>/proposal.md`：
  - 功能定位与边界
  - 验收标准（GIVEN/WHEN/THEN 场景）
  - 成功指标

**Gate: PRD 签批**——PM 确认 proposal.md 内容可接受。用户确认后进入③b。

### 阶段 ③b：UI/UX 设计（与③c 并行）

- ui-designer 用配置的设计 MCP 产出线框图 → 高保真 → 交互原型 → 状态矩阵。
- 不阻塞③c（技术架构可以与设计并行推进）。

**Gate: 设计签批**——PM + Tech Lead + UI Designer 三方确认。
设计冻结后传给 frontend-dev 做④a。

### 阶段 ③c：技术架构（Propose → design.md）

- tech-lead 产出 `openspec/changes/<name>/design.md`：
  - 技术方案（模块划分、组件交互、数据流）
  - API 契约（`.api` / `.proto` 文件）
  - 数据模型（GORM 定义）
  - 架构决策记录

**Gate: API 冻结**——契约确定后前后端可并行开发。

### 阶段 ③d：任务拆分（Propose → tasks.md）

- tech-lead 将设计拆解为 tasks.md：
  - 标注依赖关系
  - 标注可并行项
  - 每个任务指向具体的 agent（backend-dev / frontend-dev / devops / test-engineer）

**Gate: Sprint 承诺**——PM + Tech Lead 确认 scope 可实现。users 确认后可进入④。

### 阶段 ④：实现 + 测试（Apply）

- 调用 `dispatch-dev` 按 tasks.md 派发实现。
- test-engineer 在 spec propose 后就按 GIVEN/WHEN/THEN 持续产出测试。
- code-reviewer 审查每个 PR。
- PM 做 UAT 验收。

**Gate: UAT 通过**后进入⑤。

### 阶段 ⑤：部署 + 发布

- devops 构建部署。
- test-engineer 部署后测试。
- PM + DevOps 做发布决策。

**Gate: 发布门**后进入⑥。

### 阶段 ⑥：复盘 + Archive

- `/opsx:archive`：合并 delta 到 `openspec/specs/`。
- PM 回写 `PROJECT.md`（勾选功能、追加变更记录）。
- 团队做简短回顾：哪些做得好、哪些可改进。

**循环：** 下一个功能回到③。

## Gate 快速对照

| 门 | 位置 | 条件 | 签字 |
|----|------|------|------|
| PRD 签批 | ③a→③b | proposal.md 用户认可 | PM |
| 设计签批 | ③b→③c | 所有 UI 状态覆盖可实施 | PM + Tech Lead + UI Designer |
| API 冻结 | ③c→③d | 契约确定 | Tech Lead |
| Sprint 承诺 | ③d→④ | scope 可实现 | PM + Tech Lead |
| CR 通过 | ④b | 每 PR 审查 | Code Reviewer |
| UAT 通过 | ④d→⑤ | 验收标准满足 | PM |
| 发布门 | ⑤→⑥ | 测试通过 + 监控就绪 + 预案确认 | PM + DevOps |

## 完成后

报告：本次 change 的名称、pass 了哪些门、当前处于哪个阶段、还剩哪些门。

## 禁止

- 不跳过 proposal 直接实现（必须先对齐）。
- 不在③a 未确认时进入③c（需求不清晰就不要做架构）。
- 不在③c 未确认时进入④（API 没冻结就开工是浪费）。
- 不跳过 archive 就直接开始下一个 change。
