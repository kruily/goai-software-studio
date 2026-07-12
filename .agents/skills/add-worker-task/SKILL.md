---
name: add-worker-task
description: 当需要新增异步任务时使用。通过 pkg/taskqueue 抽象定义任务 topic、在 api logic 中 Enqueue、在 mq 消费者中注册 handler 回表执行。单体走 ServiceGroup 内的 mq 目录，微服务走模块的 mq 服务。遵循「队列只带 TaskID、详情回表」规范。
---

# add-worker-task

你为后端新增异步任务，走 `shared/pkg/taskqueue` 抽象。遵循工作室队列规范。

## 使用时机

- 某操作耗时/可异步（如生成、批处理、通知），需投递到队列后台执行。
- 用户说"加个异步任务""后台跑""放队列"。

## 核心规范

- **只依赖抽象**：用 `taskqueue.TaskQueue` 接口，不直接依赖 Asynq/Kafka SDK。
- **队列只带 TaskID**：`Payload{TaskID}`；任务详情落 task 表，消费时回表读，避免队列成为业务数据源。
- **投递链**：api logic → task 表落库 → `TaskQueue.Enqueue(topic, taskID, opts)` → 消费者 `Consume` → 回表执行。
- **合进程**：单体的消费者在模块 `mq/` 目录下实现 handler，由 api 进程的 `ServiceContext` 挂载启动。微服务在模块 `mq/` 独立部署（含独立入口）。

## 执行流程

### 1. 定义任务与 topic

- 约定 topic 命名，如 `task.{域}.{action}`（前缀订阅友好）。
- 如需 task 表记录任务状态（pending/running/success/failed，用 `consts` 常量），先 `gorm-add-model` 建 task 模型。

### 2. 投递（生产侧）

- 在 api logic 中：落库任务记录 → `svcCtx.TaskQueue.Enqueue(ctx, "task.{域}.{action}", taskID, opts)`。

### 3. 注册消费者（消费侧）

mq 全手写，结构：`{module}/mq/internal/{handler,logic,svc,config}`。

- **handler**（`internal/handler/{event}.go`）：桥接到 logic：
  ```go
  func {Event}Handler(svcCtx *svc.ServiceContext) queue.HandlerFunc {
      return func(ctx context.Context, event *queue.Event) error {
          return logic.New{Event}Logic(ctx, svcCtx).{Event}(event)
      }
  }
  ```
- **routes.go**（`internal/handler/routes.go`）：集中注册 topic→handler 映射。
- **logic**（`internal/logic/{event}.go`）：按 `msg.TaskID` 回表读详情 → 执行 → 更新任务状态。
- **单体**：logic 通过 api 的 `ServiceContext` 挂载到 api 进程。
- **微服务**：`{module}/mq/{module}.mq.go` 独立入口部署（ServiceGroup + 队列消费者）。

### 4. 验证

- `go build ./...`；确认订阅 Pattern 与投递 topic 匹配。

## 完成后

报告：新增的 topic、投递点、消费者、task 模型（若有）、构建结果。

## 禁止

- 不在队列 payload 里塞业务详情（只带 TaskID）。
- 不直接依赖具体队列 SDK（走 taskqueue 抽象；换实现用 add-infra-adapter）。
- 消费失败不做无限自动重试（按 EnqueueOptions.MaxRetry 控制）。
