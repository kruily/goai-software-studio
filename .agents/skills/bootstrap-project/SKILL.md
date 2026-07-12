---
name: bootstrap-project
description: 开新项目时使用。访谈用户产出项目设定文档 PROJECT.md，据此确定技术架构（单体/微服务、数据库、队列、存储、前端、设计工具、代码智能工具、module 前缀），脚手架化后端骨架，写好四端 agent 配置与 MCP，产出技术选型文档 rules/tech-selection.md，并初始化 openspec。是工作室开发生命周期第②③步的引擎。
---

# bootstrap-project

你是工作室的项目启动向导。你把一个"想法"推进到"可以开始功能迭代的骨架项目"，
对应生命周期的第 ② 步（项目设定）与第 ③ 步（技术架构）。详见 `rules/lifecycle.md`。

## 使用时机

- 用户克隆模板后第一次开项目。
- 用户说"开始一个新项目""帮我 bootstrap""初始化项目"。
- 已有 PROJECT.md 但需要重新确认/调整技术架构时。

## 核心原则

- **访谈只问够定架构的最小信息**，不追求一次问完所有产品细节——功能细节留给后续 OpenSpec 迭代。
- **不替用户拍板技术栈**：给出推荐并说明理由，但选型由用户确认。工作室只固定 go-zero 作为后端框架。
- **每一步产出落到文件**，不停留在对话里：PROJECT.md、rules/tech-selection.md、四端配置、openspec/。
- 遵循 AGENTS.md 的行为约束：最小改动、构建即验收、指令不清先澄清。

## 执行流程

### 阶段 1：访谈 → 产出 PROJECT.md

依据 `references/project-md-template.md` 模板，在**仓库根**创建 `PROJECT.md`。围绕以下维度访谈用户，边问边填：

1. **一句话定位**：这是什么产品、为谁解决什么问题。
2. **目标用户**：主要用户群体与使用场景。
3. **核心功能清单**：粗粒度列出，不展开细节。
4. **非功能约束**（这些直接决定架构）：
   - 规模量级（并发/数据量/预期用户数）
   - 是否需要实时（SSE/WebSocket 流式）
   - 是否需要异步任务（队列/定时任务）
   - 是否含 AI 能力（大模型/多模态）

**访谈完成标志**：PROJECT.md 已写入仓库根，且用户已确认当前内容可以接受。
**在用户确认之前，不要进入阶段 2。** 如果用户不确定功能细节，不要替用户假设，就写已有信息。
**在用户主动说类似"可以了""继续""开始搭建"之前，不要进入阶段 3。**

把回答写入新建的 `PROJECT.md` 对应章节，更新"最后更新"时间。

### 阶段 2：推荐并确认技术架构

依据 PROJECT.md 的非功能约束，给出**带理由的推荐**，逐项与用户确认。参考 `rules/tech-stack-catalog.md`：

- **架构**：单体 / 微服务。规模小、边界未清晰 → 推荐单体（可后续演进）；已知多团队/强隔离 → 微服务。
- **module 前缀**：询问用户确认 module 前缀。推荐格式 `github.com/{organization}/{project}`，纯本地可直接用项目名。
- **数据库**：PostgreSQL / MySQL（GORM 驱动）。
- **消息队列**：需要异步 → Asynq(Redis) / Kafka；否则先只留 taskqueue 抽象不接实现。
- **对象存储**：MinIO(开发) / OSS / S3。
- **前端/客户端**（可多选，平铺为顶层目录）：Web(Vite) / Admin / 移动端(Flutter) / 服务端渲染 / 纯 API。
- **UI 设计 MCP**（可多选）：Magic(21st.dev) / Figma / shadcn。让用户选自己可用的。
- **代码智能工具**：gopls（默认必装）/ Serena（推荐）/ CodeGraphContext（可选加装）。

**技术选型确认标志**：所有选型已逐项确认，且用户已表示认可。
**在用户明确说"可以""开始""继续"之前，不要进入阶段 3。**
如果用户在访谈或选型中说"不确定""你定""先这样"等模糊表述，应主动追问一句：
"那我就按当前选型开始搭建了，可以吗？" 得到肯定答复后再进入阶段 3。

### 阶段 3：脚手架化（逐一执行以下步骤）

按确认的架构落地。**不要跳过阶段 1 和阶段 2 直接进入阶段 3**。
每步完成后向用户报告，等待用户确认后再继续下一步。

**Step 3.1：替换 module 前缀**
```bash
find . -type f \( -name "*.go" -o -name "*.tpl" -o -name "go.mod" \) -exec sed -i '' 's/GOAI_MODULE/{选定的module}/g' {} \;
```
验证：`cd backend && go build ./...` 通过。

**Step 3.2：写 .api 定义（不手建目录）**

先写 .api 接口定义和 import.api，**不要手动创建 handler/logic/svc 等目录**。
这些目录由 goctl 在 Step 3.4 中自动生成。

只创建 `backend/{service_name}/api/desc/` 这一个目录用于存放 .api 文件。
**单体服务名由用户确认。** 例如用户选择 app、server、admin 等，不要自己定名。
**整个后端只有一个服务目录，不要按域拆模块。**

```
backend/
└── {service_name}/
    └── api/desc/
        ├── import.api
        ├── front/
        │   └── {domain}.api**
```

**Step 3.3：写 .api 定义（续）**
在 `{service_name}/api/desc/front/{domain}.api` 按规范定义首条接口。
```go
syntax = "v1"

type GetProfileReq {
    Id int64 `json:"id"`
}

type GetProfileResp {
    Id   int64  `json:"id"`
    Name string `json:"name"`
}

@server (
    prefix: /api/v1/front/{domain}
    group:  {domain}
)
service app-api {
    @doc "获取资料"
    @handler GetProfile
    post /getProfile (GetProfileReq) returns (GetProfileResp)
}
```
更新 `import.api`：`import "front/{domain}.api"`。

**Step 3.4：调用 gen-api.sh 生成代码**
```bash
backend/scripts/gen-api.sh {service_name}/api/desc/import.api
```
验证：`cd backend && go build ./...` 通过。
此时 `{service_name}/api/internal/` 下已生成 handler/logic/svc/types/config/routes。

**Step 3.5：实现 logic 占位**
在生成的 `{service_name}/api/internal/logic/` 下，找到生成的方法签名，
在 TODO 位置填入最小实现（至少能返回 success 响应），确保 HTTP 端到端可通。

**Step 3.6：按数据库选型补全驱动**
按选型补 database driver。例如选 PostgreSQL：
- 在 `go.mod` 加 `gorm.io/driver/postgres`
- 在 `shared/pkg/database/database.go` 的 `openPostgres` 中填入 `gorm.Open(postgres.Open(dsn), &gorm.Config{})`
- `go mod tidy && go build ./...` 验证

**Step 3.7：建 model 基础模型**
按 gorm-add-model 规范在 `model/{domain}/` 建首个模型（可空但目录和 register 必须存在）。

**Step 3.8：前端/客户端**
若选了前端/客户端形态，调用 `scaffold-frontend` 技能。

**Step 3.9：deploy/ 基础**
调用 `scaffold-deploy` 技能生成 Dockerfile、docker-compose、nginx（若微服务）配置。

### 阶段 4：写四端配置与 MCP

- 更新 `.mcp.json`：启用 gopls；按选择启用 Serena；**写入所选的 UI 设计 MCP**（Magic/Figma/shadcn 的配置块）。
- 同步到 `opencode.json`、`.codex/config.toml`、`.pi/settings.json`（调用 `sync-agents` 技能保证一致）。

### 阶段 5：产出 tech-selection.md + 初始化 openspec

1. 写 `rules/tech-selection.md`：完整记录本次所有选型决定与理由，供整个工作室共享参考。
2. 同步把架构表回写进 `PROJECT.md` 的"技术架构"章节。
3. **安装并初始化 OpenSpec**（进入第 ④ 步功能迭代的前提）：
   - 先检查 Node 环境：`node --version`（需 ≥ 20.19）与 `npm --version`。
   - 若 Node 未安装或版本不足，提示用户安装 Node 20.19+，**不要继续尝试**。
   - Node 就绪后执行：
   ```bash
   npm install -g @fission-ai/openspec@latest
   openspec init                                 # 生成结构 + 各 agent 的 /opsx 命令
   ```
   - `openspec init` 会为已装的 agent 生成 `/opsx:*` 命令与 `openspec/AGENTS.md`（勿手改）。
   - 若后续 `/opsx:*` 命令不可用，回退方案是：手动读 `openspec/changes/<name>/*.md` 并按 `spec-driven` 技能的分步流程执行。

## 完成后

向用户汇报：产出了哪些文件、确定了什么架构、下一步用 `spec-driven` 开始第一个功能。
不要自动开始写业务功能——那是下一步用户驱动的事。

## 禁止

- 不在未确认选型时擅自脚手架化。
- 不手写 goctl 会生成的文件。
- 不把技术栈写死进 AGENTS.md（选型属于 PROJECT.md / tech-selection.md）。
