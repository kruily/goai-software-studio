# AGENTS.md

> 本文件是本工作室**唯一的规范真源**。Codex、opencode、pi 原生读取本文件；
> Claude Code 通过 `CLAUDE.md` 的 `@AGENTS.md` 导入本文件。
> 修改规范只改这里，再用 `sync-agents` 技能校验四端一致。

## 这是什么

一个 **Go 为主的全栈 AI 开发工作室模板**。克隆后按固定生命周期把一个想法推进到上线代码，
所有项目共享同一套规范、约束与工具链，并可被 Claude Code / Codex / opencode / pi 四类 agent 一致理解。

## 开发生命周期（核心工作流）

详见 `rules/lifecycle.md`。六步：

1. **想法** — 一句话或一段话描述要做什么。
2. **项目设定** — `bootstrap-project` 访谈产出 `PROJECT.md` 并确定技术架构，脚手架化、写四端配置、初始化 `openspec/`。
3. **功能规划** — 按功能迭代。PM 产出 `proposal.md`（需求细化），UI Designer 产出设计稿，Tech Lead 产出 `design.md`（架构+API 契约）和 `tasks.md`（任务拆分）。每步有评审门（PRD 签批 → 设计签批 → API 冻结 → Sprint 承诺）。
4. **实现 + 测试** — `dispatch-dev` 派发 sub-agents 并行开发前端+后端。`code-reviewer` 每 PR 审查。`test-engineer` 按 spec 持续产出测试。PM 做 UAT 验收。
5. **部署 + 发布** — DevOps 构建部署，Test Engineer 部署后测试，PM + DevOps 做发布决策。
6. **复盘 + Archive** — OpenSpec archive 合并 specs，回写 `PROJECT.md`，回顾。

**PROJECT.md vs openspec/**：`PROJECT.md` 是全局产品蓝图（粗、稳）；`openspec/` 是功能级增量（细、频繁）。archive 时增量回流，`PROJECT.md` 保持当前全貌。每个功能经过 7 道 Gate（PRD 签批 → 设计签批 → API 冻结 → Sprint 承诺 → CR 通过 → UAT 通过 → 发布门）后才关闭。

## 仓库结构

```
├── AGENTS.md / CLAUDE.md            # 规范真源 + Claude 桥接
├── .mcp.json                        # 共享 MCP（gopls + serena）
├── opencode.json / .codex / .pi     # 四端 agent 配置
├── PROJECT.md                       # 项目设定文档（bootstrap 生成，开始项目后存在）
├── openspec/                        # 功能级 spec（proposal/design/tasks + specs 真源）
├── rules/                            # 工作室固定规范文档（生命周期/Gate/agent/后端/前端/测试）
├── .agents/skills/                  # 工具型技能真源（Agent Skills 标准格式，四端通用）
├── .claude/skills -> ../.agents/skills  # symlink，让 Claude Code 读到同一批技能
├── backend/                         # Go 后端（clone go-ai-backend-template 后存在），见 backend/README.md
├── frontend/ admin-web/ mobile/ ... # 前端/客户端，平铺，由脚手架生成
├── deploy/                          # 运维部署（docker/compose/k8s/nginx/ci/env）
└── backend/scripts/             # gen-*.sh（来自 go-ai-backend-template）
```

## 后端开发规范

后端**框架固定为 go-zero**（`shared/goctl` 代码生成模板与 `shared/pkg` 抽象均围绕它）。
但**具体技术栈由每个项目自行选型**——数据库、消息队列、对象存储、前端形态等，
由 `bootstrap-project` 访谈确定并写入**技术选型文档 `rules/tech-selection.md`**（供整个工作室共享参考）。
本工作室不预设这些选型；此处只规定 go-zero 服务应如何开发。完整规范见 `rules/backend-spec.md` 与 `backend/README.md`（clone `go-ai-backend-template` 后存在）。

- **只手写两类东西**：`backend/shared/`（pkg 抽象 + utils + 改造版 goctl 模板）与 `.api`/`.proto`/`model`。其余全部由 `goctl` 生成，**不手写 goctl 生成物**（main/svc/handler/routes/config 等）。
- **基础设施走抽象**：业务代码只依赖 `shared/pkg` 的接口（`storage.Storage`、`taskqueue.TaskQueue` 等），不直接依赖具体厂商 SDK。换厂商用 `add-infra-adapter` 技能，业务零改动。
- **模块即业务域**：直接命名 `user`、`order`，不套 `services/` 外壳。单体与微服务同构建模。
  - **单体**：一个 go-zero 服务；`.api` 用 `group:` 注解落到 `handler/{域}`、`logic/{域}`；`mq`/`cron` 手写目录，用 go-zero `ServiceGroup` 与 REST 合进一个进程。
  - **微服务**：同名模块水平扩展，每模块可含 `api`/`rpc`/`mq`/`cron`，各自独立部署；多 api 服务由 `deploy/nginx` 转发（不用 go-zero gateway）。
- **统一响应体不入 .api**：`{code,msg,data}` 只存在于 `shared/utils/response`；`.api` 只描述业务 data。
- **错误码**：业务错误返回 `errorx.CodeError`（集中码表），HTTP 状态统一 200，客户端读 `body.code` 判断。
- **代码生成**：用改造版模板 `--home ./backend/shared/goctl`；封装在 `gozero-add-api` / `gorm-add-model` 技能中。

## 技能（Skills）

**核心模型：技能驱动对话，agent 承载执行。** 交互式对话（需求访谈、确认、推进）由技能在主线程承载；
agent 是被委派的执行单元，在隔离上下文跑完返回，不与用户直接对话。入口是 `studio` 技能。
详见 `rules/agents-guide.md`。

工具型技能采用 **Agent Skills 开放标准**（`SKILL.md` 格式），真源放 `.agents/skills/{name}/SKILL.md`，四端通用：

- **Codex / opencode / pi** 原生读取 `.agents/skills/`（零配置）。
- **Claude Code** 通过 `.claude/skills -> ../.agents/skills` 的 symlink 读取同一批。
- frontmatter 只用规范严格子集（`name` 需匹配目录名、`description`），保证四端兼容。

| 阶段 | 技能 | 用途 |
|------|------|------|
| 入口 | `studio` | 总入口向导，判断当前阶段并引导到对应技能 |
| 项目设定/架构 | `bootstrap-project` | 访谈→PROJECT.md→定架构→脚手架→写四端配置→产出 tech-selection.md→init openspec |
| 全程 | `sync-agents` | AGENTS.md 真源同步/校验到四端配置 |
| 全程 | `author-skill` | 编写项目自己的业务技能 |
| 功能迭代 | `spec-driven` | 引导 OpenSpec propose→apply→archive |
| 派发 | `dispatch-dev` | 读 tasks.md 派发 sub-agents 并行开发 |
| 后端 | `gozero-add-api` | 加 go-zero API 模块 |
| 后端 | `gorm-add-model` | 加 GORM 模型 + migrate 注册 |
| 后端 | `add-worker-task` | 加异步任务（mq 消费者） |
| 后端 | `add-infra-adapter` | 实现 pkg 抽象的新适配器 |
| 前端 | `scaffold-frontend` | 按选型初始化前端项目 |
| 部署 | `scaffold-deploy` | 按选型生成 nginx/Dockerfile/compose/CI 配置 |
| 测试 | `generate-tests` | 从 OpenSpec spec 场景生成 Go 单元/集成测试 |
| 测试 | `e2e-runner` | 用 Playwright 执行浏览器 E2E 测试 |
| 测试 | `load-test` | 用 k6 执行 API 负载/压力测试 |

> 无技能触发器的场景（如 Codex/pi 的隐式调用），可直接读对应 `.agents/skills/{name}/SKILL.md` 按步骤执行。

## 常用命令

在 `backend/` 下（clone `go-ai-backend-template` 后）：

```bash
go build ./...                 # 构建（提交前必须通过）
go vet ./...                   # 静态检查
gofmt -w <files>               # 格式化
# API 代码生成（改 .api 后）：见 gozero-add-api 技能封装的 goctl 命令
```

## 行为约束（成本与质量控制）

> 以下为默认硬约束。项目可在自身 `AGENTS.md` 追加，或用 `AGENTS.override.md` 临时强制覆盖（Codex 逃生舱）。

- **严格按需求实现**，不推断额外需求，不过度设计，优先最小直接改动。
- **只改与任务相关的文件**，不重构无关代码，不擅自引入新依赖（需批准）。
- **默认不跑测试、不新增测试**，除非用户明确要求"运行测试"。
- **构建即验收**：改动后运行构建命令；构建通过即停，不追加额外"改进"；构建失败只修与失败相关的错误。
- **不自动反复重试**，不做无谓的自我反思循环。
- 指令不清时**先澄清再动手**；受阻时说明卡点并给出最小下一步。
- 输出精简：列出改动文件、简述变更；不堆砌冗长日志与解释。
- **导航优先用代码智能工具**（gopls MCP 的定义/引用/符号搜索）而非全仓 grep；避免不必要的整仓扫描。

## 代码智能工具

- **gopls MCP**（默认）— Go 语义导航：定义、引用、包 API、诊断。
- **Serena**（推荐）— 符号级检索与编辑，跨语言（含前端）。
- **CodeGraphContext**（可选加装）— 调用/依赖图，用于"谁调用了 X""改动影响范围"。

配置见 `.mcp.json` 与 `rules/code-intelligence.md`；由 `bootstrap-project` 按需写入。

## 安全

- `docker-compose` / `.env` 中的凭证均为开发默认值，**不提交生产密钥**。
- 密钥、AccessSecret 等禁止入库、入日志、下发客户端。
- AI/存储/队列供应商一律走抽象接口。
