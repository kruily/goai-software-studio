# Agent 指南

工作室是 **agent 驱动**的：7 个专职 agent 各司其职，沿开发生命周期协作。
定义真源在 `.agents/agents/{name}.md`，由 `sync-agents` 生成到四端（Claude/Codex/opencode/pi）。

## 阵容与职责

| Agent | 角色 | 绑定技能 | 生命周期 |
|-------|------|---------|----------|
| **project-manager** | 产品需求负责人：聊想法、维护 PROJECT.md、写 OpenSpec proposal | bootstrap-project, spec-driven | ①②③④ |
| **tech-lead** | 技术负责人：架构决策、写 design、拆 tasks、派发 | spec-driven, dispatch-dev | ④⑤ |
| **backend-dev** | Go 后端：go-zero api/rpc/mq/cron、model、pkg 适配器 | gozero-add-api, gorm-add-model, add-worker-task, add-infra-adapter | ⑤ |
| **frontend-dev** | 前端/客户端：脚手架、实现、对接 API | scaffold-frontend | ⑤ |
| **ui-designer** | UI 设计：用设计 MCP 出组件页面 | （用 MCP） | ⑤ |
| **devops** | 运维：deploy/ 下 docker/compose/nginx/CI | — | ⑤ |
| **code-reviewer** | 代码审查：只读检查边界违规 | —（只读） | ⑤ |
| **test-engineer** | 测试：读 spec 生测试（单测/集测/E2E/load），报 bug、验修复 | generate-tests, e2e-runner, load-test | ⑤ |

## 协作流（agent 驱动的生命周期）

```
用户想法
  → project-manager：澄清 → PROJECT.md → OpenSpec proposal
  → tech-lead：架构决策 → design.md → 拆 tasks.md → 派发
      ├→ backend-dev：后端实现
      ├→ ui-designer → frontend-dev：设计 → 前端实现
      └→ devops：部署物料
      └→ test-engineer：生成测试 → 执行 → 报 bug → 验修复
  → code-reviewer：边界审查兜底
  → test-engineer：生成测试、执行、报 bug、验修复
  → project-manager / spec-driven：archive 回写 specs + PROJECT.md
```

## 职责分界（避免重叠）

- **PM 管需求、TL 管技术**：PM 描述"要什么"，TL 决定"怎么实现"并派发。
- **TL 不重复造派发轮子**：OpenSpec 的 tasks.md 已拆任务；TL 用 `dispatch-dev` 按其派发。
- **code-reviewer 只读**：报告问题不改代码。
- **test-engineer 只测不改**：生成测试、执行、报 bug、验证修复，但不动生产逻辑。
- **devops 物料只在 deploy/**：不往应用代码塞部署逻辑。

## 跨端机制

- 定义真源：`.agents/agents/{name}.md`（MD+frontmatter，tools 用小写）。
- 生成目标：`.claude/agents/`（MD，tools 首字母大写 + skills 预加载）、`.opencode/agents/`（MD + permission.skill）、`.pi/agents/`（MD）、`.codex/agents/{name}.toml`（TOML，正文→developer_instructions）。
- 由 `sync-agents` 技能生成/校验。改了真源后运行它。

## per-agent 专属技能

- Claude Code：agent frontmatter 的 `skills:` 预加载；去掉 `Skill` 工具可禁用技能调用。
- opencode：`permission.skill` glob 白名单。
- 推荐的外部技能见 `recommended-skills.md`。

## 调用方式

- Claude Code：自动委派（按 description）或 `@agent-{name}`。
- opencode：自动或 `@{name}`，或 Task 工具。
- Codex：显式请求子 agent（Ultra 级可自动）。
- pi：subagent 扩展（`.pi/agents/`）。
