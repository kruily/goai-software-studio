---
name: gozero-add-api
description: 当需要为后端新增一个 REST API 模块或给已有模块加接口时使用。编写 .api 定义（用 group 注解分域）、更新 import.api、用改造版 goctl 模板生成 handler/logic/types，然后在 logic 中实现业务。遵循工作室 go-zero 规范：统一响应体不入 .api、业务逻辑不入 handler。
---

# gozero-add-api

你为后端新增 API 模块或接口。遵循工作室 go-zero 规范（见 `rules/backend-spec.md`、`backend/README.md`）。

## 使用时机

- 新增一个业务域的 REST 接口，或给已有域加接口。
- 用户说"加个接口""新增 API""做 user 模块的 xxx"。

## 核心规范

- **统一响应体不入 .api**：`.api` 只描述业务 data；`{code,msg,data}` 由 `shared/utils/response` 处理。
- **命名**：`{Module}{Action}Req` / `{Module}{Action}Resp` / `{Module}Item`；路由 `/api/v1/{front|admin}/{domain}/{action}`（action 为 camelCase 动词，如 getProfile）；全部使用 POST，状态用显式字符串。
- **域内聚**：`@server(group: {域})` 让 goctl 生成到 `handler/{域}`、`logic/{域}`。
- **handler 不写业务**：只解析请求、调 logic。业务在 logic。
- **用改造版模板**：`--home ./backend/shared/goctl`（保证走 response.Success/Error、透传错误码）。

## 执行流程

### 0. 前置：目录确认

**不要手动创建模块目录**。`mkdir` 创建的目录结构与 goctl 生成的不一致。
如果模块目录还不存在，先运行 gen-api.sh 让它通过 goctl 生成初始结构：

```bash
# 首次生成会报 .api 文件不存在的错误，goctl 会提示需要哪些目录。
# 此时只写 .api 文件和 import.api，然后运行：
backend/scripts/gen-api.sh {module}/api/desc/import.api
# goctl 会自动创建 handler/logic/svc/types/config/routes 等目录和文件。
```

如果 goctl 提示 desc 目录不存在，就只创建 `{module}/api/desc/` 这一个目录，其余让 goctl 生成。

### 1. 写 .api

- 在 `{module}/api/desc/front/{域}.api`（或 `admin/`）定义 type 与 service。
- **必须先写 .api 文件**，再用 gen-api.sh 生成代码。不手动创建 handler/logic/svc/types 目录。
- 用 `@server(group: {域}, prefix: /api/v1/front/{域})`。全部使用 POST，路径为 camelCase 动词，如 `post /getProfile`。
- **不使用 go-zero 内建 jwt 中间件**，在 `@server` 中用 `middleware:` 注解挂载 studio 自定义鉴权中间件。

### 2. 注册到 import.api

- 在 `{module}/api/desc/import.api` 中 `import "front/{域}.api"`。

### 3. 生成代码（只写 .api，不手建目录）

用生成脚本（封装了改造版模板与校验插件），从仓库根运行。**脚本会自动创建 handler/logic/svc 等目录**：

```bash
backend/scripts/gen-api.sh               # 默认 user 模块
backend/scripts/gen-api.sh order/api/desc/import.api  # 其他模块
```

> 脚本内部用 `--home ./shared/goctl` 保证走 response.Success/Error、透传错误码。不要手敲原生 goctl 命令绕过脚本。

### 4. 实现 logic

- 在生成的 `{module}/api/internal/logic/` 下填业务。
- 依赖从 `svcCtx` 取（DB、pkg 抽象、领域 service）。
- 用户身份用 `contexts.GetUserID(ctx)`；错误返回 `errorx.New(code, msg)` 或 `errorx.NewByCode(code)`。
- 需要新模型 → 调 `gorm-add-model`；需要异步 → 调 `add-worker-task`。

### 5. 装配

- 若引入新依赖（model、service），在 `ServiceContext` 中装配。
- `go build ./...` 验证。

## 完成后

报告：新增/修改的 .api、生成的文件、实现的 logic、构建结果。

## 禁止

- 不在 .api 里定义 code/msg/data 通用返回体。
- 不在 handler 写业务逻辑。
- 不手改 goctl 生成的 types.go / routes。
- 不用原生 goctl 模板（必须 `--home ./shared/goctl`）。
