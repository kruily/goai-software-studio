---
name: dispatch-dev
description: 当功能规划完成（tasks.md 就绪）后，把多个可并行任务派发给 sub-agents 并行开发。生命周期第④步（实现+测试）的编排技能。
---

# dispatch-dev

你把一个功能的 tasks.md 派发给 sub-agents 并行实现，对应生命周期第 ④ 步（实现 + 测试）。

## 使用时机

- 功能规划已完成（③a-③d 所有门已通过），`tasks.md` 已就绪。
- PM 说"可以开始开发了"或用户说"实现这个功能"。
- **不要**在③c（API 冻结）之前派发实现——没有契约后端和前端都会跑偏。

## 执行流程

### 1. 读任务

- 读 `openspec/changes/<name>/tasks.md` 与对应 `design.md`。
- 识别任务的**依赖关系**与**可并行标记**。数据模型/接口契约类任务通常是其他任务的前置。

### 2. 规划派发

- 把任务分层：先做前置（如 `gorm-add-model` 建模型、`gozero-add-api` 定接口），再并行做依赖它们的实现。
- 每个 sub-agent 任务要**自包含**：说明改哪个模块、遵循哪些规范（引用 AGENTS.md 后端约定）、验收标准（构建通过）。

### 3. 派发实现

- 无依赖冲突的任务**并行派发** sub-agents。
- 会改同一批文件的任务**串行**，或用隔离工作区避免冲突。
- 后端任务让 sub-agent 调用对应技能（add-api/add-model/add-worker-task）。

### 4. 汇总、审查与测试

- 收集各 sub-agent 结果，汇总改动文件。
- 运行 `code-reviewer` 子 agent 按工作室规范检查边界违规（④b）。
- 在 `backend/` 下 `go build ./...` 验证整体构建。
- test-engineer 在③a 后就已按 spec 场景持续产出测试。汇总时确认所有测试已通过。
- 若测试未通过 → 返回给对应 sub-agent 修复 → 重新 code-review → 重新测试。

### 5. 回写 + UAT

- 更新 `tasks.md` 勾选完成项。
- 确认 code-reviewer 已通过、test-engineer 测试已通过。
- 通知 PM 进行 UAT 验收（④d）。PM 确认验收标准满足后，交给 spec-driven 走 archive（⑥）。

## 完成后

报告：派发了几个任务、并行/串行安排、构建与 review 结果、还剩哪些未完成。

## 禁止

- 不派发有未满足依赖的任务（会拿到错误上下文）。
- 不让多个 sub-agent 并行改同一文件。
- 不跳过 code-reviewer 与构建验证就宣告完成。
