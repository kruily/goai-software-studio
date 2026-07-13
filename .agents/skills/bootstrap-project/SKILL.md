---
name: bootstrap-project
description: PROJECT.md 已确认后使用。只做工程管道铺设：初始化 openspec、写四端 agent 配置与 MCP、回写 PROJECT.md。不做技术选型、不推荐架构、不创建代码。技术选型由 tech-lead 在功能规划阶段完成。
---

# bootstrap-project

你是**工程管道铺设**工具。你的职责：**在 PROJECT.md 已确认的基础上**，把项目的基础工程
容器准备好。**你不做技术选型**（那是 tech-lead 在 spec-driven 中做的）、**不推荐架构**、
**不创建代码**。

## 前置条件

- `PROJECT.md` 已存在于仓库根，且用户已确认内容。
- studio 已完成需求访谈。

## 核心原则

- 只做管道铺设（openspec init、四端配置）,不碰任何与技术选型/代码相关的事。
- tech-lead 在 spec-driven 的 ③c 阶段做技术选型，你不替用户推荐或确认。

## 执行流程

### 阶段 1：管道铺设

1. 回写 `PROJECT.md` 更新时间戳。
2. 初始化 openspec（尝试 `openspec init`，不成则手建目录骨架）。
3. 四端 agent 配置同步（`.mcp.json`、`opencode.json`、`.codex/config.toml`、`.pi/settings.json`）。
4. 检查 `.claude/skills` symlink 是否有效。
5. 告诉用户：工程管道铺设完成。下一步是 "功能规划阶段，tech-lead 会协助你做技术选型和 architecture design"。

**bootstrap-project 不做的事（由 tech-lead 在 spec-driven 的 ③c 完成）：**
- ❌ 不推荐或确认技术选型（数据库、架构、队列、前端等）。
- ❌ 不创建任何后端/前端目录或模块。
- ❌ 不替换 GOAI_MODULE。
- ❌ 不写 .api 文件。
- ❌ 不调用 gen-api.sh。

## 禁止

- ❌ 不做技术选型（那是 tech-lead 的事）。
- ❌ 不创建任何代码文件或模块目录。
- ❌ 不替换 GOAI_MODULE。
- ❌ 不写 .api 文件。
