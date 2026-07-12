---
name: scaffold-frontend
description: 当需要初始化前端或客户端项目时使用。按 bootstrap 阶段选定的形态（Web SPA/Admin 后台/移动端/服务端渲染）用官方脚手架在仓库顶层平铺生成项目，并配置对接后端 JSON API（统一响应体 code/msg/data）。前端各自使用脚手架自带结构。
---

# scaffold-frontend

你按选型初始化前端/客户端项目。前端不强加内部结构——用官方脚手架生成、沿用其约定。

## 使用时机

- bootstrap 阶段确定了前端/客户端形态，需要落地。
- 用户说"建前端""加个管理后台""初始化 App"。

## 布局约定

- **顶层平铺**：前端/客户端是仓库顶层平级目录，**不套 clients/ 外壳**。常见：`frontend/`、`admin-web/`、`mobile/`、`desktop/`。
- **结构靠脚手架**：由官方脚手架生成，沿用其目录约定，模板不干预内部结构。

## 执行流程

### 1. 确认形态与技术

- 从 `rules/tech-selection.md` / `PROJECT.md` 读已选形态；未定则询问。
- 常见选择：
  - Web SPA / Admin：Vite + (React/Vue/Svelte)。
  - 移动端：Flutter。
  - 服务端渲染：Templ + HTMX + Tailwind（在 backend 内，见 tech-stack-catalog）。

### 2. 用官方脚手架生成

- 在仓库顶层运行对应脚手架（示例）：
  - `npm create vite@latest admin-web -- --template vue-ts`
  - `flutter create mobile`
- 生成后不重排其目录结构。

### 3. 对接后端

- 配置 API base URL 与鉴权（JWT，放 Authorization 头）。
- **按统一响应体解析**：后端返回 `{code, msg, data}`，`code==0` 为成功，非 0 按 `errorx` 码处理；封装一个响应拦截器统一处理。
- 若选了 UI 设计 MCP（Magic/Figma/shadcn），说明如何用它生成组件到该项目。

### 4. 记录

- 在 `rules/tech-selection.md` 补充前端选型；在 `PROJECT.md` 技术架构表更新。

## 完成后

报告：生成了哪个目录、用的脚手架、API 对接配置、下一步。

## 禁止

- 不把前端塞进 clients/ 包裹目录。
- 不重排脚手架生成的内部结构去套后端的分层。
- 不在前端硬编码密钥或后端服务端密钥。
