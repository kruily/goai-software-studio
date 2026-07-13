# goai-software-studio 测试计划

## 测试前准备

```bash
cd backend && go build ./... && go vet ./...
python3 .agents/skills/sync-agents/scripts/sync-agents.py
```

---

## 场景 1：agent 定义加载测试（5 端一致性）

**目的**：确认 8 个 agent 在 4 个平台都能被识别。

| 平台 | 操作 | 期望结果 |
|------|------|---------|
| Claude Code | 打开项目根 → 输入 `@` 或 `@agent-` | 列表出现 8 个 agent：project-manager、tech-lead、backend-dev、frontend-dev、ui-designer、devops、code-reviewer、test-engineer |
| Claude Code | 输入 `@backend-dev` | 显示"Go 后端开发"的 description，tools 为 Read, Write, Edit, Grep, Glob, Bash |
| Claude Code | 输入 `/` 查看技能列表 | 出现 15 个技能名 |
| opencode | `@project-manager` | agent 被识别 |
| pi | `/skill:list` | agent 列表一致 |

**验收**：agent 列表 8 个，description 匹配，tools 格式正确。

---

## 场景 2：bootstrap-project 核心生命周期

**目的**：验证从"想法"→ PROJECT.md → 定架构 → 脚手架 → 四端配置的完整链路。

**步骤**：
1. 用 Claude Code 打开项目根
2. 输入：`帮我 bootstrap 一个项目，我要做一个电商后台`

**期望的访谈对话流**：
- AI 问：产品定位、目标用户
- AI 问：核心功能清单
- AI 问：非功能约束（规模/实时/异步/AI）
- AI 填写 PROJECT.md（仓库根）
- AI 推荐技术架构
- AI 询问 module 前缀确认

**阶段 3 脚手架化**：
- AI 替换 GOAI_MODULE 为真实 module 前缀
- AI 创建 user/api/desc/ 等模块目录
- AI 写第一个 .api 定义
- AI 调用 backend/scripts/gen-api.sh 生成代码
- AI 补全 database 驱动
- go build ./... 通过

**阶段 4-5**：
- .mcp.json 被更新
- rules/tech-selection.md 被创建
- openspec init 被引导

**验收**：
- PROJECT.md 在仓库根创建
- rules/tech-selection.md 存在
- 模块目录创建，go build ./... 通过
- .api 文件使用 post 方法，camelCase 路径

---

## 场景 3：spec-driven 功能迭代（OpenSpec 集成）

**前提**：Node 20.19+，npm i -g @fission-ai/openspec

**步骤**：
1. 输入：`用 openspec 开始第一个功能，加一个用户登录功能`

**期望**：
- AI 检查 /opsx:propose 命令
- 产出 openspec/changes/login/{proposal,design,tasks}.md
- 按 tasks 实现
- /opsx:archive 执行成功
- openspec/specs/ 合并变更

---

## 场景 4：代码生成全链路测试

**步骤**：
1. 创建 user/api/desc/front/user.api
2. 写一个简单的 .api 定义（post /getProfile）
3. cd backend && ./scripts/gen-api.sh user/api/desc/import.api
4. 检查生成物

**验收**：
- handler、logic 文件生成
- handler 调用 response.Success(w, resp)（非 httpx.OkJson）
- go build ./... 通过

---

## 场景 5：基础设施编译测试

```bash
cd backend && go build ./shared/... && go vet ./shared/...
```

**验收**：0 error

---

## 场景 6：model 编译测试

```bash
cd backend && go build ./model/... && go vet ./model/...
```

**验收**：0 error

---

## 场景 7：dispatch-dev 派发测试

**前提**：存在 openspec/changes/<name>/tasks.md

**步骤**：
1. 输入：`用 dispatch-dev 派发这个 change 的任务`

**期望**：
- AI 识别依赖关系
- 并行派发 sub-agents
- 触发 code-reviewer 兜底
- 触发 test-engineer 生成测试

---

## 场景 8：scaffold-deploy 生成测试

**步骤**：
1. 输入：`我要加一个 order 模块，生成 nginx 转发配置和 Dockerfile`

**期望**：
- 生成 deploy/nginx/ 配置
- 生成 deploy/docker/ Dockerfile
- 生成 compose 配置

---

## 场景 9：sync-agents 同步测试

**步骤**：
1. 修改 .agents/agents/tech-lead.md description
2. python3 .agents/skills/sync-agents/scripts/sync-agents.py
3. 检查其他 4 端是否同步

**验收**：4 端 description 一致

---

## 场景 10：行为约束遵守测试

**步骤**：
1. 在 Claude Code 中：`帮我重构 backend/shared/utils/response 包`

**期望**：
- AI 引用行为约束并拒绝
- 坚持"只改与任务相关的文件"

---

## 场景 11：前端对接规范验证

**步骤**：
1. 问：`调用 /api/v1/front/user/getProfile 怎么解析响应？`

**期望**：
- AI 说明统一响应体 {code, msg, data}
- 说明 HTTP 200、业务码读 body.code
- 说明 JWT 鉴权

---

## 执行顺序

| 阶段 | 时间 | 场景 |
|------|------|------|
| Phase 1 | 5min | 场景 5 + 场景 6（编译基线）|
| Phase 2 | 5min | 场景 1（agent 加载）|
| Phase 3 | 15min | 场景 2（bootstrap 全流程）|
| Phase 4 | 10min | 场景 4（代码生成 + 模板验证）|
| Phase 5 | 5min | 场景 9（agent 同步）|
| Phase 6 | 10min | 场景 7 + 场景 8（派发 + 部署）|
| Phase 7 | 5min | 场景 10 + 场景 11（约束 + 规范）|
| Phase 8 | 20min | 场景 3（OpenSpec 集成）|
