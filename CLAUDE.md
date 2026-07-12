@AGENTS.md

# CLAUDE.md — Claude Code 专属补充

上方 `@AGENTS.md` 导入了工作室的规范真源。以下仅为 Claude Code 专属约定，
不重复真源内容。修改通用规范请改 `AGENTS.md`。

## 技能（Skills）

技能真源在 `.agents/skills/`（13 个），通过 `.claude/skills -> ../.agents/skills` symlink 供 Claude Code 读取。
完整技能清单见 `AGENTS.md` 的「技能」章节。Claude Code 按 `description` 自动匹配触发，
也可用 `/<skill-name>` 手动调用。

## 子 agent（专职 agent）

8 个专职 agent 定义在 `.claude/agents/`，由 `sync-agents.py` 从 `.agents/agents/` 真源自动生成：

| agent | 职责 | 绑定技能 |
|-------|------|---------|
| `project-manager` | 产品需求负责人 | bootstrap-project, spec-driven |
| `tech-lead` | 技术负责人 | spec-driven, dispatch-dev |
| `backend-dev` | Go 后端开发 | gozero-add-api, gorm-add-model, add-worker-task, add-infra-adapter |
| `frontend-dev` | 前端/客户端开发 | scaffold-frontend |
| `ui-designer` | UI 设计 | （用设计 MCP） |
| `devops` | 运维部署 | — |
| `code-reviewer` | 边界审查（只读） | — |
| `test-engineer` | 测试工程师 | generate-tests, e2e-runner, load-test |

## 规划与确认

- 非平凡的实现任务优先进入 plan mode，先对齐方案再动手。
- 涉及派发多个 sub-agents 的功能开发，用 `dispatch-dev` 技能编排。
