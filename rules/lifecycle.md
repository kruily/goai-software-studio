# 开发生命周期

本工作室把一个想法推进到上线代码，走固定的五步。每一步有对应的产物、技能与工具。

```
① 想法
   ↓
② 项目设定 —— bootstrap-project 访谈 → PROJECT.md（活文档）
   ↓
③ 技术架构 —— 依据 PROJECT.md 定架构 → 脚手架化 + 写四端配置 + 初始化 openspec/
   ↓
④ 功能迭代 —— OpenSpec：/opsx:propose → apply → archive（回写 specs 与 PROJECT.md）
   ↓
⑤ UI 设计 + 派发 —— 设计 MCP 出 UI → dispatch-dev 派发 sub-agents 实现 → code-reviewer 兜底
```

---

## ① 想法

一句话或一段话说明要做什么产品、解决什么问题、大致规模。信息不必完整——第 ② 步会补齐。

## ② 项目设定 → `PROJECT.md`

运行 `bootstrap-project`，它先访谈**够定架构的最小信息**（产品定位、目标用户、核心功能、是否需要实时/异步/AI、规模量级），产出 `PROJECT.md`。

- `PROJECT.md` 是**活文档**：全局产品蓝图，粗粒度、稳定、少变。
- 后续每次功能迭代（第 ④ 步）结束时回写更新，始终反映"当前全貌"。

## ③ 技术架构

`bootstrap-project` 依据 `PROJECT.md` 推荐并与你确认：

- 架构：单体 / 微服务
- 数据库、消息队列、对象存储选型
- 前端/客户端形态（Web / Admin / 移动端 / 服务端渲染）
- UI 设计 MCP（Magic / Figma / shadcn，多选）
- 代码智能工具（gopls 默认 / Serena / CGC）
- module 前缀，`bootstrap-project` 会询问确认

确认后：脚手架化目录、用改造版 goctl 模板生成后端骨架、写好 `.mcp.json` 与四端配置、`openspec init`，并产出**技术选型文档 `rules/tech-selection.md`**（记录所有选型决定，供整个工作室共享参考）。

## ④ 功能迭代（OpenSpec）

每个功能是一次 OpenSpec change：

```
/opsx:propose <功能>   # 产出 openspec/changes/<name>/{proposal,design,tasks}.md
/opsx:apply            # 按 tasks 实现
/opsx:archive          # 合并 delta 到 openspec/specs（真源），并更新 PROJECT.md 对应章节
```

- `openspec/specs/` 是系统当前行为的真源；`changes/` 是提案增量。
- 后端实现调用 `gozero-add-api` / `gorm-add-model` / `add-worker-task` 等技能。
- 用 `spec-driven` 技能引导整个循环。

## ⑤ UI 设计 + 派发开发

- **UI 设计**：用第 ③ 步选配的设计 MCP，从 `PROJECT.md`/`tasks.md` 出组件与页面。
- **派发**：`dispatch-dev` 读 `tasks.md`，把可并行的任务派发给 sub-agents 实现。
- **兜底**：`code-reviewer` 子 agent 按工作室规范检查边界违规。

---

## 产物对照

| 产物 | 粒度 | 何时更新 |
|------|------|----------|
| `PROJECT.md` | 全局产品蓝图（由 bootstrap 生成） | 第 ②③ 步创建，每次迭代回写 |
| `rules/tech-selection.md` | 技术选型决定 | 第 ③ 步创建，选型变更时更新 |
| `openspec/specs/` | 系统当前行为真源 | 每次 archive |
| `openspec/changes/<name>/` | 单功能提案 | 每次 propose |
| 代码 | 实现 | apply / 派发开发 |
