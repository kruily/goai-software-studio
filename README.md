# goai-software-studio

一个 **Go 为主的全栈 AI 开发工作室模板**（harness 层）。它提供开发流程、规范约束和 agent 编排能力，
**不直接包含后端代码**。后端骨架来自独立模板仓库 `github.com/kruily/go-ai-backend-template`，
由 tech-lead 在第一个功能迭代中派发 backend-dev 去克隆。

可被 **Claude Code / Codex / opencode / pi** 四类 AI agent 一致理解与驱动。

## 特点

- **纯 harness 层**：不包含后端 Go 代码。工作室只有规范、agent 定义和技能。后端来自独立模板。
- **单一真源**：`AGENTS.md` 是唯一规范源，桥接到四端；`sync-agents` 保持一致。
- **agent 驱动**：8 个专职 agent 沿生命周期协作。交互对话由技能在主线程承载，agent 只做被委派的执行。
- **技能标准化**：14 个技能采用 Agent Skills 开放格式，四端通用。
- **spec 驱动**：OpenSpec 贯穿功能规划到复盘；`PROJECT.md` 管全局产品蓝图。
- **后端固定 go-zero**、技术栈由项目选型：单体/微服务同构，可平滑演进。

## 开发生命周期

```
① 想法（用户一句话）
   → ② 项目设定（studio 技能）
       聊需求 → PROJECT.md → 管道铺设(openspec init+四端配置) → 引导功能规划
   → ③ 功能规划（spec-driven，每功能）
       需求(proposal) → 设计(UI稿) → 技术架构(含选型) → 任务拆分(tasks)
   → ④ 实现+测试（dispatch-dev）
       后端先clone模板+首次脚手架 → 前后端并行 → code-reviewer → 持续测试 → UAT
   → ⑤ 部署+发布（devops）
       构建部署 → E2E/负载测试 → 发布决策
   → ⑥ 复盘+Archive
       archive specs + 回写 PROJECT.md + 回顾
```

每阶段结束后有明确的门（PRD 签批 → 设计签批 → API 冻结 → Sprint 承诺 → CR 通过 → UAT 通过 → 发布门）。
详见 `rules/lifecycle.md`。

## 核心模型：技能驱动对话，agent 承载执行

- **交互对话由技能在主线程承载。** `studio`、`spec-driven` 在主线程和用户对话。
- **agent 是执行单元。** backend-dev、frontend-dev 等由 tech-lead 派发，在隔离上下文执行。
- **入口是 `studio` 技能**：判断当前阶段并引导到对应技能。
- **对话式协作**：不做 auto-pilot。提问 → 给选项 → 用户决定 → 批准后才落地。

## 后端模板

工作室不直接提供后端代码。后端骨架由 `go-ai-backend-template` 独立模板提供：

```bash
# tech-lead 在第一个功能迭代中派发 backend-dev 执行
git clone git@github.com:kruily/go-ai-backend-template.git backend
cd backend
find . -type f \( -name "*.go" -o -name "*.tpl" -o -name "go.mod" \) -exec sed -i '' 's/GOAI_MODULE/github.com\/myorg\/myapp/g' {} \;
go build ./...
```

`studio` 技能只做需求访谈和管道铺设，不执行脚手架。脚手架由 tech-lead -> backend-dev 派发。
后端模板包含：shared/pkg（storage/taskqueue/database/jwtx/redisx 接口+实现）、shared/utils（response/errorx/contexts/consts）、shared/goctl（60个改造版代码生成模板）、model/base（BaseRepo泛型CRUD）、scripts（gen-api/gen-rpc/gen-model）。

## 快速开始

### 1. 克隆
```bash
git clone <模板仓库 URL> my-project
cd my-project
```

### 2. 安装工具链
```bash
./backend/scripts/setup-agents.sh
```

### 3. 用 AI agent 启动

用 Claude Code / Codex / opencode / pi 打开项目，然后告诉它：

> "我想开始一个项目"

流程会自动推进：
- `studio` 技能聊需求 → 产出 `PROJECT.md`
- 用户确认后 → `studio` 铺设管道（openspec init + 四端配置）
- 进入功能规划阶段 → tech-lead 在 ③c 做技术选型
- tech-lead 派发 backend-dev 执行脚手架（clone go-ai-backend-template + 替换 GOAI_MODULE + 首次代码生成）

## 8 个专职 agent

所有 agent 都是 subagent（被技能委派的执行单元），不与用户直接对话。

| Agent | 职责 | 绑定技能 |
|-------|------|---------|
| **project-manager** | 产品需求：聊想法、维护 PROJECT.md、写 proposal | spec-driven |
| **tech-lead** | 技术负责人：架构决策、design、拆 tasks、派发（含首个脚手架任务）| spec-driven, dispatch-dev |
| **backend-dev** | Go 后端：首次 clone 模板 + api/rpc/mq/cron、model、pkg | gozero-add-api, gorm-add-model, add-worker-task, add-infra-adapter |
| **frontend-dev** | 前端/客户端：脚手架、实现、对接 API | scaffold-frontend |
| **ui-designer** | UI 设计：用设计 MCP 出组件与页面 | （用设计 MCP） |
| **devops** | 运维部署：docker、compose、nginx、CI | scaffold-deploy |
| **code-reviewer** | 边界审查：按规范检查架构违规 | —（只读） |
| **test-engineer** | 测试工程师：从 propose 后就开始持续测试 | generate-tests, e2e-runner, load-test |

## 14 个工具型技能（四端通用）

| 阶段 | 技能 | 作用 |
|------|------|------|
| 入口 | `studio` | 总入口向导，判断当前阶段并引导 |
| 项目设定 | `studio` | 需求访谈 → PROJECT.md → 管道铺设 → 引导 spec-driven |
| 全流程 | `sync-agents` | AGENTS.md 真源同步到四端 + MCP + agent 定义 |
| 全流程 | `author-skill` | 编写项目自己的业务技能 |
| 功能迭代 | `spec-driven` | 引导 OpenSpec propose → apply → archive |
| 派发 | `dispatch-dev` | 读 tasks.md 派发 sub-agents 并行开发 |
| 后端 | `gozero-add-api` | 加 go-zero API 模块 |
| 后端 | `gorm-add-model` | 加 GORM 模型 + migrate 注册 |
| 后端 | `add-worker-task` | 加异步任务（mq 消费者） |
| 后端 | `add-infra-adapter` | 实现 pkg 抽象的新适配器 |
| 前端 | `scaffold-frontend` | 按选型初始化前端项目 |
| 部署 | `scaffold-deploy` | 按选型生成 nginx/Dockerfile/compose/CI |
| 测试 | `generate-tests` | 从 spec 场景生成 Go 单元/集成测试 |
| 测试 | `e2e-runner` | 用 Playwright 执行浏览器 E2E |
| 测试 | `load-test` | 用 k6 执行 API 负载/压力测试 |

## 目录

```
├── AGENTS.md / CLAUDE.md         # 规范真源 + Claude 桥接
├── .mcp.json                     # 共享 MCP（gopls + serena）
├── opencode.json / .codex / .pi  # 四端 agent 配置
├── .agents/
│   ├── skills/                   # 14 个工具型技能真源（四端通用）
│   └── agents/                   # 8 个专职 agent 定义真源
├── .claude/                      # Claude 版 agents + skills symlink
├── PROJECT.md                    # 项目设定文档（bootstrap 生成）
├── openspec/                     # 功能级 spec（OpenSpec）
├── rules/                        # 12 份工作室固定规范文档
├── backend/                      # 仅 README（指向 go-ai-backend-template）
│   └── README.md                 # 说明：clone 模板到此目录
├── deploy/                       # 运维部署约束说明
└── frontend/ / admin-web/ (bootstrap 生成)
```

## 文档索引

| 文档 | 内容 |
|------|------|
| `rules/lifecycle.md` | 六步开发生命周期 + 7 道评审门 |
| `rules/agents-guide.md` | 核心模型说明 + 8 个 agent 职责协作 |
| `rules/backend-spec.md` | go-zero 后端开发规范 |
| `rules/frontend-spec.md` | 前端开发规范 |
| `rules/frontend-backend-spec.md` | 前后端对接规范 |
| `rules/api-proto-conventions.md` | .api/.proto 定义规范（全部 POST）|
| `rules/tech-stack-catalog.md` | 技术栈选型菜单 |
| `rules/code-intelligence.md` | gopls / Serena / CGC 代码智能工具 |
| `rules/recommended-skills.md` | 按 agent 推荐的外部社区技能 |
| `rules/testing-conventions.md` | Go 测试规范 |
| `rules/template-boundary.md` | 模板三层边界模型 |
| `rules/service-split-patterns.md` | 微服务拆分决策 |

## 约定

- **后端框架固定 go-zero**；数据库/队列/存储/前端由项目选型。
- **后端代码来自独立模板** `github.com/kruily/go-ai-backend-template`。
- **studio 只做需求访谈和管道铺设，不做技术选型**。选型由 tech-lead 在 ③c 完成。
- **全部接口使用 POST**，路径 camelCase 动词。
- **统一响应体不入 .api**。
- **基础设施走抽象**：业务只依赖 `shared/pkg` 接口。

## 后端脚本（clone 模板后存在）

```bash
backend/scripts/gen-api.sh                     # 生成 REST
backend/scripts/gen-rpc.sh <proto> <out>       # 生成 zRPC
backend/scripts/gen-model.sh <pg|mysql> ...    # 从数据库生成 model
```

## 四端支持

| 平台 | agent 机制 | skills 机制 |
|------|-----------|-------------|
| **Claude Code** | `.claude/agents/`（subagent） | `.claude/skills/` symlink |
| **opencode** | `.opencode/agents/`（subagent） | 原生读取 `.agents/skills/` |
| **Codex** | `.codex/agents/`（TOML） | 原生读取 `.agents/skills/` |
| **pi** | `.pi/agents/`（MD） | 原生读取 `.agents/skills/` |
