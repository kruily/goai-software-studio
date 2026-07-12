# .api 与 .proto 接口定义规范

go-zero 的接口契约文件。`.api` 用于 REST，`.proto` 用于 zRPC。由 AI 按本规范编写，
再用 `goctl` 生成代码。**不改生成物，只改契约文件。**

## .api 规范（REST）

### 文件组织

```
backend/
├── api/
│   ├── desc/
│   │   ├── import.api          # 聚合入口，import 所有 .api 文件
│   │   ├── front/              # 前台接口
│   │   │   ├── user.api
│   │   │   └── order.api
│   │   └── admin/              # 后台接口
│   │       └── user.api
```

### 语法

```go
syntax = "v1"

// 类型定义
type (
    {Module}{Action}Req {
        {Field} {Type} `json:"..." validate:"..."`  // validate tag 需配校验插件
    }

    {Module}{Action}Resp {
        Id   int64  `json:"id"`
        Name string `json:"name"`
    }

    {Module}Item {
        // 列表项结构
    }
)

// 服务定义
@server (
    prefix: /api/v1/{front|admin}/{domain}  // 以小写域名结尾，例如 /api/v1/front/user
    group:  {domain}                         // 域名，goctl 据此落到 handler/{域}/、logic/{域}/
)
service app-api {
    @doc "接口功能简述"
    @handler {HandlerName}
    post /{action} ({RequestType}) returns ({ResponseType})  // action 为 camelCase 动词，如 getProfile
}
```

### 命名约定

| 元素 | 格式 | 示例 |
|------|------|------|
| 请求类型 | `{Module}{Action}Req` | `UserProfileReq` |
| 响应类型 | `{Module}{Action}Resp` | `UserProfileResp` |
| 列表项 | `{Module}Item` | `UserItem` |
| handler 名 | `{Module}{Action}` | `UserGetProfile` |
| 路由 | `/api/v1/{front\|admin}/{domain}/{action}` | `/api/v1/front/user/getProfile` |
| 状态字段 | 显式字符串，不用魔法数字 | `status: "pending" / "active" / "disabled"` |

### @server 注解

- `prefix`：路由前缀，以小写域名结尾，例如 `prefix: /api/v1/front/user`。
- `group`：**必须加**。域内聚的关键——goctl 按此落到 `handler/{group}/`、`logic/{group}/`。
- **鉴权**：**不使用 go-zero 内建 jwt 中间件**。studio 定义独立鉴权中间件，通过 `@server` 的 `middleware:` 注解挂载到需要认证的路由组。鉴权中间件放 `{module}/api/internal/middleware/`。

### 请求方法

**全部接口使用 POST**，不使用 RESTful 语义。路径为 camelCase 动词（动作名），完整路径 = `prefix + / + action`：

| 完整路由 | 语义 |
|----------|------|
| `/api/v1/front/user/getProfile` | 查用户资料 |
| `/api/v1/front/user/createUser` | 创建用户 |
| `/api/v1/front/user/updateUser` | 更新用户 |
| `/api/v1/front/order/listOrders` | 查询订单列表 |

> 使用 POST 简化客户端调用与 CORS 配置。动作名用 camelCase 动词（getX/createX/updateX/deleteX/listX）。

### 校验

字段用 `validate` tag（需启用校验插件）：

```go
type CreateUserReq {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"gte=0,lte=150"`
}
```

生成时加校验插件：`goctl api plugin -p goctl-validate="validate --translator" ...`

### import.api

聚合入口，goctl 据此一次性生成所有域：

```go
syntax = "v1"
import "front/user.api"
import "front/order.api"
import "admin/user.api"
```

### 生成

```bash
# 单体
goctl api format --dir api/desc
backend/scripts/gen-api.sh                              # 默认 user 模块
backend/scripts/gen-api.sh order/api/desc/import.api     # 其他模块
```

## .proto 规范（zRPC，微服务）

### 文件组织

```
backend/
├── api/
│   └── desc/
│       └── front/
├── rpc/
│   ├── proto/
│   │   ├── user/
│   │   │   └── user.proto
│   │   └── order/
│   │       └── order.proto
│   └── user/
│       ├── pb/ etc/ internal/
```

### 语法

```protobuf
syntax = "proto3";
package user;
option go_package = "./;user";

service UserService {
  rpc GetProfile(GetProfileReq) returns(GetProfileResp);
}

message GetProfileReq {
  int64 id = 1;
}

message GetProfileResp {
  int64  id    = 1;
  string name  = 2;
  string status = 3;
}
```

### 生成

```bash
backend/scripts/gen-rpc.sh rpc/proto/user/user.proto rpc/user
```

入口在 `rpc/user/` 下（goctl 生成）。

## 通用约定

- **统一响应体不入契约文件**：`.api` / `.proto` 只描述业务 data；`{code, msg, data}` 在 `shared/utils/response` 处理。
- 不改生成物：types.go、routes、pb 文件由 goctl 生成，改契约文件重新生成。
- 用改造版模板：`--home ./backend/shared/goctl`（保证走 response.Success/Error、透传错误码）。
