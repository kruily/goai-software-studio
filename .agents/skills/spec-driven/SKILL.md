---
name: spec-driven
description: 贯穿功能规划(③)、实现+测试(④)、部署+发布(⑤)、复盘+archive(⑥)的 OpenSpec 生命周期引导技能。proposal 产出需求+验收标准，design 产出架构+API+数据模型+技术选型，tasks 产出任务拆分。不跳过任何 Gate。
---

# spec-driven

你是工作室的**生命周期引导技能**。每个功能是一个 OpenSpec change，贯穿从③需求到⑥复盘。
详见 `rules/lifecycle.md`。

## 使用时机

- PROJECT.md 已确认，bootstrap 已铺设管道。
- 用户说"开始功能规划""写 PRD""做第一个功能"。

## 前置

- OpenSpec 已初始化（bootstrap 阶段完成）。
- 已有 PROJECT.md。

## 角色分工

| 阶段 | Agent | 产出 | 备注 |
|------|-------|------|------|
| ③a 需求细化 | project-manager | proposal.md（边界+验收标准 GIVEN/WHEN/THEN）| PM 主导 |
| ③b UI 设计 | ui-designer | 设计稿（先设计系统→原子组件→页面）| 与③c并行 |
| ③c 技术架构 | tech-lead | design.md（技术选型+API 契约+数据模型+模块划分）| TL 选型决策 |
| ③d 任务拆分 | tech-lead | tasks.md（标注依赖与可并行）| TL 产出 |
| ④a 实现 | backend-dev / frontend-dev | 代码 | dispatch-dev 派发 |
| ④b 审查 | code-reviewer | 每 PR 审查 | — |
| ④c 测试 | test-engineer | 按 spec 产出测试 | 从③a就介入 |
| ④d UAT | project-manager | 验收确认 | — |
| ⑤ 部署 | devops | deploy/ 物料 | — |
| ⑥ 复盘 | tech-lead + pm | archive + 回顾 | — |

## 技术选型（③c，由 tech-lead 完成）

tech-lead 在 ③c 阶段做技术选型。逐项推荐与确认：

- 架构：单体 / 微服务
- 数据库：PostgreSQL / MySQL
- 消息队列：Asynq / Kafka / 无
- 对象存储：MinIO / OSS / S3
- 前端形态：admin-web / studio-web / 移动端
- UI 设计 MCP：Magic / Figma / shadcn / Ardot
- module 前缀

选型确认后记录到 rules/tech-selection.md。**bootstrap-project 不做选型，tech-lead 做。**

## OpenSpec 循环

```
/opsx:propose <name>   → ③a→③b→③c→③d（每步有 Gate）
/opsx:apply            → ④（实现+测试+UAT）
/opsx:archive          → ⑥（合并 specs + 回写 PROJECT.md + 回顾）
```

## 执行流程

### 阶段 ③a：需求细化（PM → proposal.md）

- project-manager 与用户澄清功能边界、验收标准。
- 产出 openspec/changes/<name>/proposal.md。

**Gate: PRD 签批**——用户确认。确认后进入③b+③c（并行）。

### 阶段 ③b：UI 设计（ui-designer，与③c并行）

- 先建立设计系统（调色板、图标库、Logo）。
- 按调色板配色 → 原子组件 → 复合组件 → 页面。
- 每个界面覆盖 loading/empty/error/edge case 状态。
- 产出设计稿供 frontend-dev 实现。

ui-designer 使用配置的设计 MCP（Ardot/Figma/shadcn）在主线程作图。
如果设计 MCP 需要子 agent 无法访问（Ardot 通过主线程 MCP 通话），ui-designer 在主线程承担。

**Gate: 设计签批**——PM + Tech Lead + UI Designer 确认。设计冻结后传给 frontend-dev。

### 阶段 ③c：技术架构（tech-lead → design.md）

- tech-lead 产出 openspec/changes/<name>/design.md：
  - 技术选型记录
  - 架构决策（模块划分、数据流）
  - API 契约（遵循 gozero-add-api 规范：全 POST、group/prefix、camelCase）
  - 数据模型（不嵌套 BaseModel，每个模型显式定义字段）
- 数据模型中不使用 datatypes.JSON 或 map[string]interface{}（用 shared/model/custom_type）。

**Gate: API 冻结**——契约确定后④a可开始。

### 阶段 ③d：任务拆分（tech-lead → tasks.md）

- tech-lead 将设计拆解为 tasks.md：
  - 标注依赖关系、可并行项。
  - 首个任务：backend-dev clone go-ai-backend-template + 首次脚手架。
  - 其余任务按 design.md 的模块逐一分配。

**Gate: Sprint 承诺**——PM + Tech Lead 确认。

### 阶段 ④：实现 + 测试（dispatch-dev）

- tech-lead 用 dispatch-dev 按 tasks.md 派发 sub-agents。
- backend-dev 的先执行 clone 模板 + 首次脚手架。
- 后端任务调用 gozero-add-api / gorm-add-model / add-worker-task。
- test-engineer 从③a 就介入，按 proposal 的 GIVEN/WHEN/THEN 持续产出测试。
- code-reviewer 审查每个 PR。
- PM 做 UAT。

**Gate: UAT 通过**后进入⑤。

### 阶段 ⑤：部署 + 发布

- devops 构建部署。
- test-engineer 部署后测试。
- PM + DevOps 发布决策。

**Gate: 发布门**后进入⑥。

### 阶段 ⑥：复盘 + Archive

- /opsx:archive：合并 specs。
- PM 回写 PROJECT.md。
- 简短回顾。

**循环：下一个功能回到③。**

## 禁止

- 不跳过 proposal 直接实现。
- tech-lead 不做选型，bootstrap-project 也不做——tech-lead 在③c做。
- design.md 不能包含 datatypes.JSON / map[string]interface{} / BaseModel 嵌套。
- .api 格式必须遵循 gozero-add-api 规范（全 POST、group/prefix、camelCase）。
- ③a未签批不进③b/③c；③c未冻结不进④。
