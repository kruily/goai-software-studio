# 前端开发规范

前端/客户端由 `frontend-dev` 负责，`scaffold-frontend` 技能初始化。前端不强加内部结构——
用官方脚手架生成、沿用其约定。设计由 `ui-designer` 用设计 MCP 产出。

## 1. 布局：顶层平铺

前端/客户端是仓库**顶层平级目录**，不套 `clients/` 外壳：

- `frontend/` — 主 Web 应用
- `admin-web/` — 管理后台
- `mobile/` — 移动端（如 Flutter）
- `desktop/` — 桌面端

建哪些由 bootstrap 选型决定；结构由官方脚手架生成。

## 2. 技术选型

见 `tech-stack-catalog.md`。常见：Vite + React/Vue/Svelte（SPA/Admin）、Flutter（移动端）、
Templ+HTMX+Tailwind（服务端渲染，在 backend 内，适合 AI 流式界面）。

## 3. 对接后端

- **统一响应体**：后端返回 `{code, msg, data}`。`code == 0` 成功；非 0 按 `errorx` 码处理。
  用响应拦截器统一解包与错误提示，不在每个调用处重复判断。
- **鉴权**：JWT 放 `Authorization` 头。
- **流式**：AI/实时场景用 SSE 或 WebSocket 接后端流。
- **类型安全**：可用 `openapi-to-typescript` 类技能从后端契约生成 TS 客户端。

## 4. 与设计协作

- `ui-designer` 用设计 MCP（Magic/Figma/shadcn）产出组件与页面。
- `frontend-dev` 落地设计，优先复用所选 UI 库组件，保持一致性。
- 设计要覆盖统一响应体驱动的状态：加载 / 成功 / 错误（按 code）。

## 5. 安全

- 不硬编码密钥；**绝不**在前端放服务端密钥。
- 敏感操作走后端；前端只持临时凭证（如预签名 URL、短期 token）。

## 6. 禁止

- 不把前端塞进 `clients/` 包裹目录。
- 不重排脚手架生成的内部结构去套后端分层。
- 不在前端绕过统一响应体各写各的解析。
