---
name: spec-driven
description: 用于按功能迭代开发时，引导 OpenSpec 的规范驱动流程：propose 产出 proposal/design/tasks，apply 实现，archive 回写 specs 真源与 PROJECT.md。是生命周期第④步的引导技能，把一个功能需求推进成可派发的任务清单并落地。
---

# spec-driven

你引导 OpenSpec 的规范驱动开发流程，对应生命周期第 ④ 步（功能迭代）。
每个功能是一次 OpenSpec change。详见 `rules/lifecycle.md`。

## 使用时机

- 架构已就绪（bootstrap 完成），要开发某个具体功能。
- 用户说"做 X 功能""加个 Y""迭代 Z"。

## 前置

- **OpenSpec 已安装并初始化**。检查 `openspec/AGENTS.md` 与 `/opsx:*` 命令是否存在；若无：
  ```bash
  npm install -g @fission-ai/openspec@latest   # 需 Node 20.19+
  openspec init                                 # 生成结构 + 各 agent 的 /opsx 命令
  ```
  - 若无法使用 `/opsx:*` 命令，按本技能的步骤手动创建 `openspec/changes/<name>/` 目录与 proposal/design/tasks 文件。
- 已有 `PROJECT.md`（bootstrap 阶段产出）。

## 角色分工（agent 驱动）

- **project-manager**：产出 `proposal.md`（为什么做、要什么、验收标准）。
- **tech-lead**：补 `design.md`（技术方案，遵循后端规范）、拆 `tasks.md`。
- **backend-dev / frontend-dev**：按 tasks 实现。
- 本技能是这些 agent 共用的流程引导。

## OpenSpec 循环

```
/opsx:propose <功能>   → openspec/changes/<name>/{proposal,design,tasks}.md
/opsx:apply            → 按 tasks 实现
/opsx:archive          → 合并 delta 到 openspec/specs，更新 PROJECT.md
```

## 执行流程

### 1. Propose（产出提案）

- 与用户澄清这个功能的**边界与验收标准**（够写 tasks 即可，不过度设计）。
- 运行/引导 `/opsx:propose <name>`，产出：
  - `proposal.md`：为什么做、改什么。
  - `design.md`：技术方案——落到工作室约定（哪个模块、走哪些 `.api`/`model`、是否需要 mq/cron、用哪些 pkg 抽象）。
  - `tasks.md`：可执行清单，**标注可并行项**，供 dispatch-dev 派发。
- design 必须遵循后端规范（不手写 goctl 生成物、基础设施走抽象、统一响应体不入 .api）。

### 2. Apply（实现）

- 小功能可直接 apply；多任务/可并行的功能 → 交给 `dispatch-dev` 技能派发 sub-agents。
- 后端实现调用 `gozero-add-api` / `gorm-add-model` / `add-worker-task`。
- 前端实现按选型进行。

### 3. Archive（回写）

- 功能完成、构建通过后 `/opsx:archive`：把 delta 合并进 `openspec/specs/`（系统当前行为真源）。
- **更新 PROJECT.md**：勾选对应功能、追加"变更记录"一行，更新时间。

## 完成后

报告：本次 change 的名称、产出的 tasks、实现状态、specs/PROJECT.md 是否已回写。

## 禁止

- 不跳过 propose 直接写代码（规范驱动的核心是先对齐再实现）。
- 不在 design 里违反后端规范。
- 完成后不忘回写 specs 与 PROJECT.md（否则真源漂移）。
