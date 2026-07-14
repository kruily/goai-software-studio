# 技术栈选型菜单

`studio` 访谈时的可选项与接入方式参考。**这里是候选菜单，不是强制**——
除后端框架固定 go-zero 外，其余由项目按需选，结果记入 `tech-selection.md`。

## 架构

| 选项 | 何时选 | 落地 |
|------|-------|------|
| 单体 | 规模小、边界未清晰、快速起步 | 一个模块的 api 进程中挂载 mq+cron |
| 微服务 | 多服务独立部署 | 各模块的 api/rpc/mq/cron 独立部署 |

> 单体可平滑演进到微服务（module → service）。

### 微服务/单体：模块化结构（不用 go.work）

单体和微服务共用同一套模块结构。区别仅在于：单体用一个进程跑模块的 api+mq+cron，
微服务让模块的 api/rpc/mq/cron 各自独立部署。

```
backend/
├── go.mod                             # 一个 module，所有模块共属一个 Go 项目
├── shared/                            # 跨模块手写资产（pkg、utils、goctl）
├── model/                             # GORM 模型（base、各域 model、migrate.go）
├── scripts/                           # 生成/初始化脚本
└── {module}/                          # 业务模块，如 user、order
    ├── api/                           # REST 服务
    │   ├── desc/                      # .api 定义（front/admin/）
    │   ├── etc/                       # 配置 yaml
    │   ├── internal/{config,handler,logic,svc,types,middleware}
    │   └── {module}.api.go            # 入口 main（goctl 生成）
    ├── rpc/                           # zRPC 服务（微服务时才需要）
    │   ├── pb/ etc/ internal/         # goctl 生成
    │   └── {module}.rpc.go            # 入口 main
    ├── mq/                            # MQ 消费者（手写；单体时挂载到 api 进程）
    │   ├── internal/{handler,logic,svc,config}
    │   └── {module}.mq.go             # 入口 main（微服务独立部署时）
    └── cron/                          # 定时任务（手写；单体时挂载到 api 进程）
        ├── internal/{handler,logic,svc,config}
        └── {module}.cron.go           # 入口 main（微服务独立部署时）
```

- 开发/小团队**不推荐 go.work**：单 module 零跨模块 ceremony、一次 `go mod tidy`、直接 import 共享包。go.work 适合多团队独立版本控制的大规模场景。
- go-zero 官方 `bookstore` 示例即采用单 module 多目录模式。
- **单体**：入口在 `{module}/api/{module}.api.go`。mq 与 cron 的 logic 通过 api 的 `ServiceContext` 挂载到 api 进程，不拆独立入口。`rpc/` 不存在。
- **微服务**：`{module}/api`、`{module}/rpc`、`{module}/mq`、`{module}/cron` 各自为独立部署单元（各有入口文件 `{module}.{type}.go`）。多 api 由 `deploy/nginx` 转发（**不使用 go-zero gateway**）。
- 演进：从单体到微服务只需把 mq/cron 取出加上独立入口，模块边界不变，代码零重构。
- 版本管理：各服务 Docker 镜像独立 tag；代码库统一管理，按功能分支开发。

## 后端框架

- **go-zero**（固定）— REST + zRPC + 代码生成。工作室的 `shared/goctl`、`shared/pkg` 均围绕它。

## 数据库（GORM 驱动）

| 选项 | 驱动 | 备注 |
|------|------|------|
| PostgreSQL | `gorm.io/driver/postgres` | 默认推荐，jsonb 支持好 |
| MySQL | `gorm.io/driver/mysql` | 广泛使用 |
| SQLite | `gorm.io/driver/sqlite` | 本地/轻量 |

连接经 `shared/pkg/database`，按 `Config.Driver` 分发。

## 消息队列

| 选项 | 何时选 | 落地 |
|------|-------|------|
| 无 | 无异步需求 | 只保留 `taskqueue` 抽象不接实现 |
| Asynq (Redis) | 常规异步任务 | `taskqueue` 的 asynq 适配器（工作室默认） |
| NATS | 低延迟 pub/sub、多消费者组、无中间件 | `taskqueue` 的 nats 适配器 |
| Kafka | 高吞吐、日志流、事件溯源 | `taskqueue` 的 kafka 适配器 |
| Pulsar | 多租户、多订阅、持久化 | 按 `taskqueue` 抽象实现适配器 |

> 各队列实现放 `shared/pkg/taskqueue/` 下，用 `add-infra-adapter` 技能新增。

## 对象存储

| 选项 | 何时选 | 实现 |
|------|-------|------|
| MinIO | 开发/自托管 | `storage` 的 minio 适配器 |
| OSS | 阿里云生产 | `storage` 的 oss 适配器 |
| COS | 腾讯云生产 | `storage` 的 cos 适配器 |
| S3 | AWS / 通用云 | `storage` 的 s3 适配器 |
| R2 | Cloudflare | `storage` 的 r2 适配器 |

各实现放 `shared/pkg/storage/` 下，用 `add-infra-adapter` 技能新增。

## 前端 / 客户端（顶层平铺）

| 选项 | 技术 | 目录 |
|------|------|------|
| Web SPA / Admin | Vite + React/Vue/Svelte | `frontend/`、`admin-web/` |
| 移动端 | Flutter | `mobile/` |
| 服务端渲染 | Templ + HTMX + Tailwind（在 backend 内） | 适合 AI 聊天/流式界面 |
| 纯 API | 无前端 | — |

## UI 设计 MCP（可多选）

| 选项 | 作用 | 配置 |
|------|------|------|
| Magic (21st.dev) | 自然语言→生成组件 | `@21st-dev/magic`，需 API key |
| Figma MCP | 读既有设计交实现 | 官方 `mcp.figma.com` 或 `figma-developer-mcp` |
| shadcn/ui MCP | 从注册表拉组件 | `npx shadcn@latest mcp` |
| **Ardot（腾讯）** | 腾讯自研 AI 设计工具，类 Figma；支持组件/变量/Dev Mode 导出 CSS/React/Vue | 通过 Ardot 桌面客户端 UI 配置 MCP（非 JSON）；详见 [docs.ardot.tencent.com](https://docs.ardot.tencent.com) |

> Ardot MCP 提供 18 个工具（查询/设计系统/编辑/导出/辅助），通过桌面客户端右上角 MCP 配置入口启用。不通过 `.mcp.json` 配置，由 bootstrap 时在 AGENTS.md 说明启用方式。

由 bootstrap 写入 `.mcp.json` 与四端，`sync-agents` 保持一致。

## 代码智能工具

见 `code-intelligence.md`。gopls（默认）+ Serena（推荐）+ CodeGraphContext（可选）。

## AI 能力（若需要）

- 若项目含大模型/多模态：在业务侧定义自己的 provider 抽象接口（工作室不预置 aiprovider），
  按需接入 Anthropic/OpenAI Go SDK 或 Eino。流式经 SSE/WebSocket，`context` 可取消。
