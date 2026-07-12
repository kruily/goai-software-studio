---
name: dispatch-dev
description: 当一个功能的 tasks.md 就绪、需要把多个可并行任务派发给 sub-agents 并行开发时使用。读取 openspec/changes/<name>/tasks.md，拆分任务、派发实现、汇总结果，并用 code-reviewer 兜底检查。是生命周期第⑤步的派发编排技能。
---

# dispatch-dev

你把一个功能的任务清单派发给 sub-agents 并行实现，对应生命周期第 ⑤ 步。

## 使用时机

- 某个 OpenSpec change 的 `tasks.md` 已产出，任务较多或可并行。
- 用户说"派发开发""并行实现这些任务""开始做这个 change"。

## 前置

- 存在 `openspec/changes/<name>/tasks.md`（由 spec-driven / OpenSpec propose 产出）。

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

### 4. 汇总与兜底（含测试触发）

- 收集各 sub-agent 结果，汇总改动文件。
- 运行 `code-reviewer` 子 agent 按工作室规范检查边界违规。
- 在 `backend/` 下 `go build ./...` 验证整体构建。
- **触发 test-engineer**：汇总后通知 test-engineer agent 调用 `generate-tests` 技能，为本次新增/改动的接口生成测试。测试通过后才进入回写。
- 若 test-engineer 发现 bug → 返回给对应 sub-agent 修复 → 重新 code-review → 重新测试。

### 5. 回写

- 更新 `tasks.md` 勾选完成项。
- 确认 test-engineer 的测试结果已通过。
- 交给 spec-driven 走 archive（回写 specs 与 PROJECT.md）。

## 完成后

报告：派发了几个任务、并行/串行安排、构建与 review 结果、还剩哪些未完成。

## 禁止

- 不派发有未满足依赖的任务（会拿到错误上下文）。
- 不让多个 sub-agent 并行改同一文件。
- 不跳过 code-reviewer 与构建验证就宣告完成。
