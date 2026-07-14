# Agent 指南

工作室是 **agent 驱动**的：8 个专职 agent 各司其职，沿开发生命周期协作。
定义真源在 `.agents/agents/{name}.md`，由 `sync-agents` 生成到四端（Claude/Codex/opencode/pi）。

## 核心模型：技能驱动对话，agent 承载执行

这是理解整个工作室的关键。多数 AI agent 平台里，专职 agent 是 **subagent**——
它们在隔离上下文里跑完返回，**结构上无法与用户持续对话**（Claude Code 的 subagent
连 `AskUserQuestion` 都用不了）。因此工作室采用如下模型：

- **交互对话由技能在主线程承载。** `studio`（需求向导）、`studio`（技术选型）、`spec-driven`（功能规划）这些技能
  在主线程运行，直接和用户来回访谈、确认、推进。它们就是 PM / tech-lead 等交互角色的"人格"。
- **agent 是被委派的执行单元。** backend-dev、frontend-dev、devops 等在隔离上下文执行
  非交互的工作，跑完返回结果。用户与它们不直接对话。
- **入口是 `studio` 技能**（`/studio`），它判断用户当前阶段并引导到对应技能。

所以"某个 agent 和用户对话"的实质，是"该 agent 对应的技能在主线程和用户对话"。
例如 PM 的交互体现在 `studio` / `spec-driven` 技能里，而 `project-manager`
agent 定义本身主要用于被委派做非交互分析（如整理需求文档）。

**对话式协作原则**（不做 auto-pilot）：提问 → 给 2-4 个选项及利弊 → 用户决定 →
起草 → 得到批准后才落地。

## 阵容与职责

| Agent | 角色 | 绑定技能 | 生命周期 |
|-------|------|---------|----------|
| **project-manager** | 产品经理：聊想法、维护 PROJECT.md、写 proposal、UAT 验收、发版决策 | studio, spec-driven | ①②③a④d⑤c⑥b |
| **tech-lead** | 技术负责人：技术选型、架构设计、API 契约、任务拆分、派发 | spec-driven, dispatch-dev | ②b③c③d④a⑤a |
| **backend-dev** | Go 后端：go-zero api/rpc/mq/cron、model、pkg 适配器 | gozero-add-api, gorm-add-model, add-worker-task, add-infra-adapter | ④a |
| **frontend-dev** | 前端/客户端：脚手架、实现、对接 API | scaffold-frontend | ④a |
| **ui-designer** | UI 设计：用设计 MCP 出组件页面 | （用设计 MCP） | ③b |
| **devops** | 运维：构建部署、CI、发布 | scaffold-deploy | ⑤a |
| **code-reviewer** | 代码审查：只读检查边界违规 | —（只读） | ④b |
| **test-engineer** | 测试：从③a 介入按 spec 持续生测试、E2E、负载 | generate-tests, e2e-runner, load-test | ③a④c⑤b |

## 协作流（agent 驱动的生命周期）

```
用户想法
  → ① 想法
  → ② project-manager: PROJECT.md → tech-lead: 技术选型 + 脚手架
  → ③ 功能规划
      ├── ③a project-manager: proposal.md ← Gate: PRD 签批
      ├── ③b ui-designer: 设计稿           ← Gate: 设计签批
      ├── ③c tech-lead: design.md + API    ← Gate: API 冻结
      └── ③d tech-lead: tasks.md           ← Gate: Sprint 承诺
  → ④ 实现+测试
      ├── ④a backend-dev + frontend-dev 并行
      ├── ④b code-reviewer: 每 PR 审查     ← Gate: CR 通过
      ├── ④c test-engineer: 持续测试
      └── ④d project-manager: UAT           ← Gate: UAT 通过
  → ⑤ 部署+发布
      ├── ⑤a devops: 构建部署
      ├── ⑤b test-engineer: 部署后测试
      └── ⑤c project-manager + devops: 发布决策 ← Gate: 发布门
  → ⑥ 复盘
      ├── ⑥a tech-lead: archive
      ├── ⑥b project-manager: 回写 PROJECT.md
      └── ⑥c 团队回顾
  → (循环：下一个功能，回到③)
```

## 职责分界（避免重叠）

- **PM 管需求、TL 管技术**：PM 描述"要什么"，TL 决定"怎么实现"并派发。
- **TL 不重复造派发轮子**：OpenSpec 的 tasks.md 已拆任务；TL 用 `dispatch-dev` 按其派发。
- **code-reviewer 只读**：报告问题不改代码。
- **test-engineer 只测不改**：③a 就介入按 spec 持续写测试，不等到代码写完再测。
- **ui-designer 在③b 设计，不在实现后才设计**：设计先于或并行于架构，而非实现后补。
- **devops 物料只在 deploy/**：不往应用代码塞部署逻辑。
- **每个阶段结束后有明确的 Gate**：不满足条件的不可进入下一阶段。

## 跨端机制

- 定义真源：`.agents/agents/{name}.md`（MD+frontmatter，tools 用小写）。
- 生成目标：`.claude/agents/`（MD，tools 首字母大写 + skills 预加载）、`.opencode/agents/`（无 tools 行）、`.pi/agents/`（MD + tools YAML 数组）、`.codex/agents/{name}.toml`（TOML，正文→developer_instructions）。
- 由 `sync-agents` 技能中的脚本自动生成。改了真源后运行它。

## per-agent 专属技能

- Claude Code：agent frontmatter 的 `skills:` 预加载；去掉 `Skill` 工具可禁用技能调用。
- opencode：`permission.skill` glob 白名单。
- 推荐的外部技能见 `recommended-skills.md`。

## 调用方式

- Claude Code：自动委派（按 description）或 `@agent-{name}`。
- opencode：自动或 `@{name}`，或 Task 工具。
- Codex：显式请求子 agent（Ultra 级可自动）。
- pi：subagent 扩展（`.pi/agents/`）。
