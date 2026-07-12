# 后端开发规范

后端框架固定为 **go-zero**；具体技术栈（数据库/队列/存储）由项目选型（见 `tech-selection.md`）。
本文规定 go-zero 服务应如何开发。速览见 `AGENTS.md`，目录说明见 `backend/README.md`。

## 1. 手写 vs 生成

**手写的东西：**

1. **`backend/shared/`** — 跨模块复用的手写资产：
   - `shared/pkg/` 基础设施**抽象接口**（storage、taskqueue、database、redisx、jwtx）。
   - `shared/utils/` 工具（response、errorx、contexts、consts）。
   - `shared/goctl/` 改造版代码生成模板。
2. **`.api` / `.proto` / `model`** — 接口契约与数据模型。
3. **业务实现代码**：
   - `api/rpc` 的 **logic**（goctl 生成空壳，业务逻辑手写填入）。
   - `mq` 消费者（handler + logic，goctl 不生成，全手写）。
   - `cron` 定时任务（handler + logic，goctl 不生成，全手写）。

**goctl 生成、不手改的东西：** `internal/` 下的 handler（api）、svc、types、config、routes、`{module}.api.go`/`{module}.rpc.go` 入口、pb 文件。用脚本 `backend/scripts/gen-*.sh` 生成。

> logic 是"生成空壳 + 手写填肉"：goctl 生成方法签名与 `todo` 占位，业务在其中手写。mq/cron 无 goctl 生成，整套手写（结构见 §7）。

## 2. 模块即业务域

一个业务域就是一个模块，直接命名 `user`、`order`，**不套 `services/` 外壳**。每个模块是**自包含的 go-zero 服务**，拥有 `api/rpc/mq/cron` 子目录。

### 结构

```
backend/
├── shared/              # 跨模块手写资产（pkg、utils、goctl）
├── model/               # GORM 模型
├── scripts/             # 生成/初始化脚本
└── {module}/            # 业务模块（如 user、order）
    ├── api/             # REST 服务（goctl 生成 + 手写 logic）
    │   ├── desc/        # .api 定义
    │   ├── etc/         # 配置 yaml
    │   ├── internal/
    │   │   ├── config/ handler/ logic/ svc/ types/ middleware/
    │   └── {module}.api.go  # ★入口 main（单体：在 ServiceContext 中挂载 mq/cron）
    ├── rpc/             # zRPC 服务（微服务时才有）
    │   ├── pb/ etc/ internal/
    │   └── {module}.rpc.go   # 入口 main
    ├── mq/              # MQ 消费者（手写）
    ├── cron/            # 定时任务（手写）
```

### 入口约定

- **微服务**：每个 `api/`、`rpc/`、`mq/`、`cron/` 目录都含自己的入口 `{module}.{type}.go`（goctl 生成的 `{module}.api.go`、`{module}.rpc.go`），各自独立部署。
- **单体**：入口统一写在 `api/{module}.api.go`（goctl 生成）。mq 与 cron 通过 `ServiceContext` 或 `ServiceGroup` 挂载到 api 进程，**不拆出独立入口**。
  - mq 消费者在 `mq/` 目录下实现 handler，在 `svc.ServiceContext` 中启动消费循环。
  - cron 定时任务在 `cron/` 目录下实现，在 `svc.ServiceContext` 中注册定时调度。
  - 这样从单体演进到微服务时，只需把 `mq/` 或 `cron/` 抽出加上独立入口 `{module}.mq.go`，代码零改动。

## 3. 分层边界

- **handler**：只解析/校验请求、调 logic。**不写业务逻辑**。
- **logic**：业务编排。依赖从 `svc.ServiceContext` 取。
- **model**：数据访问唯一入口。handler/logic 不直接拼 SQL。
- **service**（可选）：跨 logic 复用的领域服务。

## 4. 基础设施走抽象

业务代码只依赖 `shared/pkg` 的接口：
- `storage.Storage`、`taskqueue.TaskQueue`、`database`（GORM 连接）、`jwtx`、`redisx`。
- **不直接 import** MinIO/OSS/Asynq/Kafka/具体 DB SDK。
- 换厂商 = 在对应 pkg 包新增 `{driver}.go` 实现 + 按 `Config.Driver` 分发（用 `add-infra-adapter` 技能），业务零改动。

## 5. 统一响应体与错误码

- **响应体不入 .api**：`{code, msg, data}` 只存在于 `shared/utils/response`；`.api` 只描述业务 data。
- **错误**：业务错误返回 `errorx.CodeError`（集中码表 `errorx.New` / `errorx.NewByCode`）；HTTP 状态统一 200，客户端读 `body.code` 判断。
- handler 由改造版模板生成：成功走 `response.Success(w, resp)`，失败走 `response.Error(w, err)`（用 `errors.As` 透传错误码）。

## 6. 数据模型（GORM，多数据库）

- 模型目录在 `model/{domain}/{name}.go`（顶层，不在 internal 内）。
- **三层架构**：
  1. **`type()` 块**内定义 Struct + Interface（嵌入 `base.BaseRepo[T]` + 领域方法）+ 私有 struct（嵌入 `baseRepo[T]`）。
  2. **`TableName()`** 在 type 块外。
  3. **构造函数 + receiver 方法**散列在文件下半部分，每个方法带一行中文注释。
- 示例：
  ```go
  package {domain}

  import (
      "time"
      "GOAI_MODULE/model/base"
      "gorm.io/gorm"
  )

  type (
      // {Name} 表示某业务实体。
      {Name} struct {
          Id        int64          `gorm:"primaryKey;autoIncrement"`
          Name      string         `gorm:"size:128;not null"`
          Status    string         `gorm:"size:32;not null;default:pending;index"`
          CreatedBy int64          `gorm:"not null;index"`
          CreatedAt time.Time      `gorm:"autoCreateTime"`
          UpdatedAt time.Time      `gorm:"autoUpdateTime"`
          DeletedAt gorm.DeletedAt `gorm:"index"`
      }

      // {Name}Model 定义数据查询接口。
      {Name}Model interface {
          base.BaseRepo[{Name}]
          FindByName(session *gorm.DB, name string) (*{Name}, error)
      }

      {name}Model struct {
          base.BaseRepo[{Name}]
      }
  )

  func ({Name}) TableName() string { return "{table}" }

  func New{Name}Model() {Name}Model {
      return &{name}Model{BaseRepo: base.NewBaseRepo[{Name}]()}
  }

  func (m *{name}Model) FindByName(session *gorm.DB, name string) (*{Name}, error) {
      var v {Name}
      err := session.Where("name = ?", name).First(&v).Error
      return &v, err
  }
  ```
- 连接经 `shared/pkg/database`（按 `Config.Driver` 支持 postgres/mysql/sqlite）。
- 新模型注册到 `model/migrate.go`（`AutoMigrate`）。软删用 `gorm.DeletedAt`。

## 7. 异步任务与定时任务

mq 与 cron 都是**全手写**（goctl 不生成），沿用 go-zero 分层：`internal/{handler,logic,svc,config}` + 入口。

### mq（消息队列消费者）

```
{module}/mq/
├── {module}.mq.go              # 入口 main：ServiceGroup + 队列消费者
└── internal/
    ├── config/config.go
    ├── svc/service_context.go
    ├── handler/
    │   ├── routes.go           # 注册 topic → handler
    │   └── {event}.go          # handler 桥接 logic
    └── logic/{event}.go        # 业务逻辑
```

- 走 `taskqueue.TaskQueue` 抽象。队列 payload **只带 TaskID**，详情回表读。
- 投递链：api logic → task 表落库 → `Enqueue(topic, taskID, opts)` → 消费者 `Consume` → 回表执行。
- handler 只桥接：`func {Event}Handler(svcCtx) HandlerFunc { return func(ctx, event) error { return logic.New{Event}Logic(ctx, svcCtx).{Event}(event) } }`。
- routes.go 集中注册 topic→handler 映射。

### cron（定时任务）

```
{module}/cron/
├── {module}.cron.go           # 入口 main：cobra 命令
└── internal/
    ├── config/config.go
    ├── svc/servicecontext.go
    ├── handler/
    │   ├── routes.go           # 注册各 cron 为 cobra subcommand
    │   └── {job}_handle.go     # 返回 *cobra.Command，RunE 调 logic
    └── logic/{job}_logic.go    # 业务逻辑
```

- 每个定时任务是一个 cobra subcommand，由外部调度器（crontab / k8s CronJob）按 `{binary} {job}` 触发。
- handler 返回 `*cobra.Command`，`RunE` 内调 `logic.New{Job}Logic(cmd.Context(), svc).{Job}()`。

### 单体 vs 微服务

- **单体**：mq/cron 的 logic 通过 `api` 的 `ServiceContext` 挂载到 api 进程，不拆独立入口。
- **微服务**：`{module}/mq/{module}.mq.go`、`{module}/cron/{module}.cron.go` 各自独立部署。

## 8. 代码生成脚本

在 `backend/` 下：

```bash
backend/scripts/gen-api.sh [user/api/desc/import.api]    # 生成 REST（改造版模板）
backend/scripts/gen-rpc.sh <proto> <out>                  # 生成 zRPC（微服务）
backend/scripts/gen-model.sh <pg|mysql> <dsn> <table> <out>   # 生成 model
```

均用 `--home ./shared/goctl`。封装在 `gozero-add-api` / `gorm-add-model` 技能。

## 9. 命名与约定

- `.api` 类型：`{Module}{Action}Req` / `{Module}{Action}Resp` / `{Module}Item`。
- 路由：`/api/v1/{front|admin}/{domain}/{action}`（action 为 camelCase 动词，如 `getProfile`、`createUser`）。
- **全部接口使用 POST**，不使用 RESTful 语义，简化客户端调用与 CORS。
- 状态：显式字符串（`consts.StatusPending` 等），不用魔法数字。
- 用户身份：`contexts.GetUserID(ctx)`，不散落 context key。

## 10. 安全

- 密钥、AccessSecret 禁止硬编码、入库、入日志、下发客户端。
- `docker-compose` / `.env` 中的凭证均为开发默认值，不提交生产密钥。
