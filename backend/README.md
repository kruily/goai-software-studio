# backend

Go 后端。基于 **go-zero**，起步为单体、可水平扩展为微服务。

## 核心约定

**模板只手写两类东西**，其余全部由 `goctl` 依据 `.api` / `.proto` 生成：

1. **`shared/`** — 跨模块复用的手写资产（本模板已提供骨架）：
   - `shared/pkg/` — 基础设施**抽象接口**。业务代码只依赖接口，不直接绑定具体 SDK：
     - `storage` 对象存储（MinIO/OSS/S3...）· `taskqueue` 异步队列（Asynq/Kafka...）
     - `aiprovider` 大模型/多模态供应商（OpenAI/Anthropic/Eino...）
     - `database` GORM 连接（PostgreSQL/MySQL/SQLite，按 Driver 分发） · `redisx` Redis · `jwtx` JWT
   - `shared/utils/` — `response` 统一响应 · `errorx` 业务错误 · `consts` 通用常量
2. **`.api` / `.proto`** 接口契约 + **`model/`** GORM 模型 — 由 AI 按需编写。

> 换厂商（存储/队列/AI）时只新增 `shared/pkg/{域}` 下的适配器文件，业务零改动 —— 用 `add-infra-adapter` 技能。

## 模块怎么长出来（不预置任何示例）

一个业务域就是一个模块，直接命名 `user`、`order`，**不套 `services/` 外壳**。每个模块自包含 `api/rpc/mq/cron` 子目录：

```
backend/
├── shared/   model/   scripts/
└── {module}/            # 如 user、order
    ├── api/             # REST → {module}.api.go（入口 main）
    │   ├── desc/ etc/ internal/{config,handler,logic,svc,types,middleware}
    ├── rpc/             # zRPC → {module}.rpc.go（微服务；单体无）
    ├── mq/              # MQ 消费者（手写，单体挂载到 api 进程）
    ├── cron/            # 定时任务（手写，单体挂载到 api 进程）
```

- **单体**：入口在 `api/{module}.api.go`。mq/cron 在 `ServiceContext` 中挂载，不拆出独立入口。无 rpc。
- **微服务**：每个 `api/`、`rpc/`、`mq/`、`cron/` 各自为独立部署单元（各有入口文件）。多 api 由 `deploy/nginx` 转发。
- 演进不重构：单体只拆进程，模块边界不变。

`internal/` 等 **全部由 goctl 生成，模板不预置**，避免与生成物冲突。

## 相关技能

| 目的 | 技能 |
|---|---|
| 加 API 模块（写 .api → goctl 生成 → 补 logic） | `gozero-add-api` |
| 加 GORM 模型（含 migrate 注册、结构化 JSON 自定义类型、跨库可移植） | `gorm-add-model` |
| 加异步任务（mq 消费者） | `add-worker-task` |
| 新增基础设施适配器（换存储/队列/AI） | `add-infra-adapter` |

## module 前缀

`go.mod` 当前为占位 `GOAI_MODULE`，启动项目后全仓库替换为真实前缀（见 bootstrap-project 技能）。
