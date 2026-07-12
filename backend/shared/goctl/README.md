# shared/goctl — 改造版 goctl 代码生成模板

这是工作室定制的 **goctl 模板集**，让 `goctl` 生成的代码直接对齐工作室约定，
而不是原生默认输出。覆盖 `api`（REST）、`rpc`（zrpc）、`model`（GORM/SQL）、
`gateway`、`kube`、`docker`、`mongo`。

## 为什么要定制模板

- **统一响应体不入 .api**：改造版 `api/handler.tpl` 让每个 handler 成功走
  `response.Success(w, resp)`、失败走 `response.Error(w, err)`，`{code,msg,data}`
  只存在于 `shared/utils/response`，`.api` 只描述业务 data。
- **错误码可透传**：handler 失败时把 `err` 原样交给 `response.Error`，后者用
  `errors.As` 提取 `*errorx.CodeError` 的真实业务码（原生模板会丢码）。
- **参数错误码固定**：`httpx.Parse` 失败统一返回 `errorx.CodeInvalidParam`。

## 使用方式

生成时通过 `--home` 指向本目录（相对 backend 根）：

```bash
# 生成 api（单体：group 注解落到 handler/{域}、logic/{域}）
goctl api go -api api/desc/import.api --dir . --home ./shared/goctl --style go_zero

# 生成 model（GORM 风格，含自定义 customized.tpl）
goctl model pg datasource ... --home ./shared/goctl

# 生成 rpc（微服务）
goctl rpc protoc xxx.proto --home ./shared/goctl ...
```

> 具体命令由 `gozero-add-api` / `gorm-add-model` 技能封装，通常不手敲。

## 占位 module 路径

`api/handler.tpl` 中引用 `GOAI_MODULE/shared/utils/{response,errorx}`。
启动项目后用 `sed` 全仓库替换 `GOAI_MODULE` 为真实 module 前缀（见 bootstrap-project 技能），模板引用随之生效。

## 改动记录（相对原生 goctl 模板）

- `api/handler.tpl` — 走 `response.Success/Error`，透传错误码，参数错误用 `CodeInvalidParam`。
- 其余定制模板（`model/customized.tpl`、`api/sse_handler.tpl` 等）按需在此追加改动并记录。
