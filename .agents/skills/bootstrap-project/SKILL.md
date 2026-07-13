---
name: bootstrap-project
description: PROJECT.md 已确认后使用。与用户逐项讨论技术选型（单体/微服务、数据库、队列、存储、前端、设计工具、module 前缀），脚手架化后端骨架，写四端配置与 MCP，产出 rules/tech-selection.md，初始化 openspec。不做需求访谈——那是 studio 技能的职责。
---

# bootstrap-project

你是工作室的**技术选型与脚手架引导**。你的职责：**在 PROJECT.md 已确认的基础上**，与用户逐项讨论技术选型，然后按选型落地。

**你不做需求访谈**——那是 `studio` 技能的职责。如果你发现没有 PROJECT.md，请引导用户先去和 `studio` 聊需求。

## 前置条件

- `PROJECT.md` 已存在于仓库根，且用户已确认内容。
- 用户说"可以了""开始选型""开始搭建"。

## 核心原则

- **只问够定架构的信息**：PROJECT.md 里的非功能约束已经给了判断依据，不需要再聊功能细节。
- **给推荐、让用户选**：每个维度和用户确认选型，给出推荐理由和选项利弊。
- **每确认一项就记录一项**：不等到全部确认完再写。
- **确认后才动手**：所有选型用户都点头了才创建文件。

## 执行流程

### 阶段 1：技术选型咨询

逐项与用户确认。一次一项，每项给出推荐理由。参考 `rules/tech-stack-catalog.md`：

1. **架构**：单体 / 微服务。规模小、边界未清晰 → 推荐单体。
2. **module 前缀**：询问用户确认。推荐 `github.com/{org}/{project}`。
3. **数据库**：PostgreSQL / MySQL。
4. **消息队列**：需要异步 → Asynq / Kafka；否则留抽象不接。
5. **对象存储**：MinIO(开发) / OSS / S3。
6. **前端/客户端**（可多选）：Web / Admin / 移动端 / 纯 API。
7. **UI 设计 MCP**（可多选）：Magic / Figma / shadcn / Ardot。
8. **代码智能工具**：gopls（默认）/ Serena / CodeGraphContext。

每确认一项，回写到 `rules/tech-selection.md`。
全部确认后问用户："这些选型可以开始脚手架了吗？"

### 阶段 2：脚手架化

用户说"可以"后，bootstrap-project **不直接创建后端代码**。它把脚手架工作拆为任务，
由 tech-lead 派发给 backend-dev 执行。bootstrap-project 只做：

1. 产出 `rules/tech-selection.md`（全部选型记录）
2. 更新 `PROJECT.md` 的技术架构章节
3. 初始化 openspec
4. 写四端配置
5. 告诉用户：技术选型已完成，下一步由 tech-lead 派发 backend-dev 执行脚手架。

**后端骨架（clone 模板 + 替换 module + 写 .api + 生成代码）由 backend-dev agent 在
tech-lead 派发后执行，不在 bootstrap-project 中直接操作。**

### 阶段 3：收尾

把选型同步回写进 `PROJECT.md` 的技术架构章节。
告诉用户：
> "技术选型已确认。现在告诉 tech-lead 开始脚手架，他会派发 backend-dev 执行。"

## 后端脚手架任务（由 tech-lead 派发 backend-dev 执行，不在本技能直接操作）

backend-dev 收到 tech-lead 的派发后，执行以下步骤：

**Step 1：clone 后端模板**
```bash
git clone git@github.com:kruily/go-ai-backend-template.git backend
```

**Step 2：替换 GOAI_MODULE**
```bash
cd backend
find . -type f \( -name "*.go" -o -name "*.tpl" -o -name "go.mod" \) -exec sed -i '' 's/GOAI_MODULE/{module前缀}/g' {} \;
go build ./...
```

**Step 3：写 .api 定义**
在 `backend/{service_name}/api/desc/` 创建第一条接口。**不手动创建 handler/logic 目录。**

**Step 4：调用 gen-api.sh**
```bash
backend/scripts/gen-api.sh {service}/api/desc/import.api
go build ./...
```

**Step 5：按数据库选型补全驱动**
在 `shared/pkg/database/database.go` 按选型补 `openPostgres` 或 `openMySQL`。
