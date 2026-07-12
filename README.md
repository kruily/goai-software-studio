# goai-software-studio

一个 **Go 为主的全栈 AI 开发工作室模板**。克隆它开新项目：沿固定的开发生命周期，把一个想法
从 PRD 推进到上线代码；所有项目共享同一套规范、约束与工具链，并可被 **Claude Code / Codex /
opencode / pi** 四类 AI agent 一致理解与驱动。

## 特点

- **规范单一真源**：`AGENTS.md` 是唯一规范源，桥接到四端；`sync-agents` 保持一致。
- **agent 驱动**：8 个专职 agent（产品/技术负责人 + 后端/前端/UI/运维/审查/测试）沿生命周期协作。
- **技能标准化**：`.agents/skills/` 用 Agent Skills 开放格式，四端通用。
- **spec 驱动**：OpenSpec 管功能级迭代；`PROJECT.md` 管全局产品蓝图；测试从 spec 场景自动生成。
- **后端固定 go-zero**、技术栈由项目选型：单体/微服务同构，可平滑演进。
- **测试完整**：test-engineer agent 从 spec 生成 Go 单元/集成测试、Playwright E2E、k6 负载测试。

## 开发生命周期

```
① 想法  → ② project-manager: PROJECT.md
   → ③ tech-lead: 定架构+脚手架
   → ④ OpenSpec 功能迭代(propose/apply/archive)
   → ⑤ 派发实现: backend-dev/frontend-dev/ui-designer/devops
     → code-reviewer 边界审查 → test-engineer 生成测试+执行+报 bug
```

详见 `rules/lifecycle.md` 与 `rules/agents-guide.md`。

## 快速开始

模板本身就是一个 Git 仓库，**不需要运行初始化脚本**。

### 1. 克隆
```bash
git clone <模板仓库 URL> my-project
cd my-project
```

### 2. 安装工具链
```bash
./backend/scripts/setup-agents.sh
```
安装 gopls（代码智能）、goctl（代码生成）、OpenSpec（spec 驱动）。

### 3. 用 AI agent 启动
用 Claude Code / Codex / opencode / pi 打开项目，然后告诉它：

> "帮我 bootstrap 这个项目"

`project-manager` agent 会自动启动 `bootstrap-project` 技能，依次完成：
- 和你聊想法 → 产出 `PROJECT.md`（项目设定文档）
- 推荐技术架构（单体/微服务、数据库、队列、前端等）
- 替换 module 前缀 `GOAI_MODULE` 为你的项目名（如 `github.com/kruily/myapp`）
- 写四端 agent 配置与 MCP
- 初始化 openspec，准备开始功能迭代

## 8 个专职 agent

| Agent | 职责 | 绑定技能 |
|-------|------|---------|
| **project-manager** | 产品需求：聊想法、维护 PROJECT.md、写 OpenSpec proposal | bootstrap-project, spec-driven |
| **tech-lead** | 技术负责人：架构决策、design、拆 tasks、派发 | spec-driven, dispatch-dev |
| **backend-dev** | Go 后端：go-zero api/rpc/mq/cron、model、pkg | gozero-add-api, gorm-add-model, add-worker-task, add-infra-adapter |
| **frontend-dev** | 前端/客户端：脚手架、实现、对接 API | scaffold-frontend |
| **ui-designer** | UI 设计：用设计 MCP 出组件与页面 | （用设计 MCP） |
| **devops** | 运维部署：docker、compose、nginx、CI | — |
| **code-reviewer** | 边界审查：按规范检查架构违规 | —（只读） |
| **test-engineer** | 测试工程师：从 spec 生成测试、执行、报 bug | generate-tests, e2e-runner, load-test |

## 13 个工具型技能（四端通用）

| 阶段 | 技能 | 作用 |
|------|------|------|
| 项目设定 | `bootstrap-project` | 访谈 → PROJECT.md → 定架构 → 脚手架 → 写配置 → openspec init |
| 全流程 | `sync-agents` | AGENTS.md 真源同步到四端 + MCP + agent 定义 |
| 全流程 | `author-skill` | 编写项目自己的业务技能 |
| 功能迭代 | `spec-driven` | 引导 OpenSpec propose → apply → archive |
| 派发 | `dispatch-dev` | 读 tasks.md 派发 sub-agents 并行开发 |
| 后端 | `gozero-add-api` | 加 go-zero API 模块 |
| 后端 | `gorm-add-model` | 加 GORM 模型 + migrate 注册 |
| 后端 | `add-worker-task` | 加异步任务（mq 消费者） |
| 后端 | `add-infra-adapter` | 实现 pkg 抽象的新适配器 |
| 前端 | `scaffold-frontend` | 按选型初始化前端项目 |
| 测试 | `generate-tests` | 从 spec 场景生成 Go 单元/集成测试 |
| 测试 | `e2e-runner` | 用 Playwright 执行浏览器 E2E |
| 测试 | `load-test` | 用 k6 执行 API 负载/压力测试 |

## 目录

```
├── AGENTS.md / CLAUDE.md         # 规范真源 + Claude 桥接
├── .mcp.json                     # 共享 MCP（gopls + serena）
├── opencode.json / .codex / .pi  # 四端 agent 配置
├── .agents/
│   ├── skills/                   # 13 个工具型技能真源（四端通用）
│   └── agents/                   # 8 个专职 agent 定义真源
├── .claude/                      # Claude 版 agents + skills symlink
├── PROJECT.md                    # 项目设定文档（bootstrap 生成，开始后存在）
├── openspec/                     # 功能级 spec（OpenSpec）
├── rules/                        # 12 份工作室固定规范文档
├── backend/                      # Go 后端（go-zero），见 backend/README.md
│   ├── shared/                   # pkg 抽象 + utils + goctl 模板（可编译）
│   ├── model/                    # GORM 模型（BaseRepo 泛型）
│   └── scripts/                  # gen-*.sh / setup-agents.sh / sync-agents.py
├── deploy/                       # 运维部署约束说明 + env 示例
└── frontend/ admin-web/ mobile/  # 前端/客户端（bootstrap 生成，顶层平铺）
```

## 文档索引

| 文档 | 内容 |
|------|------|
| `rules/lifecycle.md` | 五步开发生命周期 |
| `rules/agents-guide.md` | 8 个 agent 的职责与协作 |
| `rules/backend-spec.md` | go-zero 后端开发规范 |
| `rules/frontend-spec.md` | 前端开发规范 |
| `rules/frontend-backend-spec.md` | 前后端对接规范（统一响应体/鉴权/流式） |
| `rules/api-proto-conventions.md` | .api/.proto 定义规范 |
| `rules/tech-stack-catalog.md` | 技术栈选型菜单 |
| `rules/code-intelligence.md` | gopls / Serena / CGC 代码智能工具 |
| `rules/recommended-skills.md` | 按 agent 推荐的外部社区技能 |
| `rules/testing-conventions.md` | Go 测试规范（单测/集测/E2E/load） |
| `rules/template-boundary.md` | 模板三层边界模型与设计哲学 |
| `rules/service-split-patterns.md` | 微服务拆分决策（5 症状 + 拆分步骤） |

## 约定

- **后端框架固定 go-zero**；数据库/队列/存储/前端由项目选型，记入 `rules/tech-selection.md`。
- **只手写两类**：`backend/shared/`（pkg 抽象 + utils + goctl）与 `.api`/`.proto`/`model`；其余由 goctl 生成。
- **全部接口使用 POST**，路径格式 `/{action}`（camelCase 动词），路由前缀含域名。
- **统一响应体不入 .api**：`{code, msg, data}` 由 `shared/utils/response` 处理。
- **基础设施走抽象**：业务只依赖 `shared/pkg` 接口，不直连具体 SDK。
- **行为约束**（最小改动、构建即验收、默认不跑测试等）见 `AGENTS.md`。

## 后端脚本

所有脚本从仓库根执行：

```bash
backend/scripts/gen-api.sh                  # 生成 REST（默认 user/api/desc/import.api）
backend/scripts/gen-api.sh order/api/...    # 其他模块
backend/scripts/gen-rpc.sh <proto> <out>    # 生成 zRPC
backend/scripts/gen-model.sh <pg|mysql>     # 从数据库生成 model
backend/scripts/setup-agents.sh             # 安装工具链
python3 backend/scripts/sync-agents.py      # 从真源同步 agent 到四端
```
