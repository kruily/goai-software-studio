---
name: add-infra-adapter
description: 当需要为某个 pkg 抽象接口实现新的基础设施适配器时使用（如新增对象存储厂商 OSS/S3、新增队列实现 Kafka、切换数据库驱动）。在对应 pkg 包内新增实现文件、实现接口全部方法、按 Config.Driver 分发，业务代码零改动。
---

# add-infra-adapter

你为 `shared/pkg` 下的某个抽象接口实现新适配器。这是"换厂商业务零改动"的落地点。

## 使用时机

- 换/加对象存储厂商（MinIO→OSS/S3/COS）。
- 换/加队列实现（Asynq→Kafka/Pulsar）。
- 切换/新增数据库驱动。
- 用户说"接入 xxx 存储""换成 Kafka""支持 OSS"。

## 核心规范

- **接口不动，只加实现**：接口定义在 `pkg/{域}/{域}.go`，实现放同包 `{driver}.go`（如 `storage/oss.go`）。
- **Config.Driver 分发**：每个 pkg 的 `Config` 有 `Driver` 字段；构造函数按它选择实现。
- **业务零改动**：业务只依赖接口，新增适配器不触碰任何 logic。
- **敏感配置**：密钥来自安全配置，禁止入库、入日志、下发客户端。

## 执行流程

### 1. 定位接口

- 读 `shared/pkg/{域}/{域}.go`，明确要实现的接口方法全集与 Config 字段。

### 2. 写适配器

- 新建 `shared/pkg/{域}/{driver}.go`，定义 struct 实现接口**全部方法**。
- 构造函数 `New{Driver}{域}(cfg Config) ({接口}, error)`。
- 复用包内公共 helper（如 `storage.DefaultPresignExpiry`）。

### 3. 接入分发

- 在包内提供/更新 `New(cfg Config)` 工厂：按 `cfg.Driver` 返回对应实现。
- 引入的第三方 SDK 加入 `go.mod`（需批准新依赖）。

### 4. 装配

- 在 `ServiceContext`（或 bootstrap 的装配处）按配置构造该实现注入。
- 更新 `rules/tech-selection.md` 记录新增的选型。

### 5. 验证

- `go build ./...`；`go vet ./...`。

## 完成后

报告：实现了哪个接口的哪个 driver、新增依赖、Config 变化、构建结果。

## 禁止

- 不修改接口签名去迁就某个 SDK（若接口不够用，先和用户确认再改接口，且要顾及所有实现）。
- 不在业务代码里直接 new 具体实现（只经工厂/ServiceContext）。
- 不硬编码密钥。
