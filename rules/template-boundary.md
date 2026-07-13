# 模板边界说明

本工作室的目录和文件分为三层。了解这三层的区别能让你一眼看出
"什么东西属于模板永远存在、什么东西需要我改、什么东西是启动项目后才生成"。

## 设计哲学：不是给小白用的

bootstrap 之后，项目不会立刻"跑起来"给你看。
这个工作室**需要有一定编程能力的人来驱动**——纯小白无法控制。

这**正是目的**：我们不做一键生成全栈应用的脚手架。模板提供的是规范、流程、约束和 agent 协作能力，
让有经验的开发者用 AI agent 更高效、更一致地把想法推进到上线。如果你不了解 go-zero 的项目结构、
不知道 go mod tidy 是做什么的、没有部署过 Go 服务——你会觉得它不好用，**但这不是 bug**。

以下是三层模型：

```
┌──────────────────────────────────────────────────────┐
│  Layer 1：永不改变（模板的规范层）                     │
│  克隆后不需要动，只需 AI agent 读取它们来约束行为      │
├──────────────────────────────────────────────────────┤
│  AGENTS.md / CLAUDE.md          规范真源 + Claude 桥接│
│  .agents/agents/* （8 个 agent 定义）                  │
│  .agents/skills/* （15 个工具型技能）                       │
│  rules/*（12 份规范文档）                                   │
│  .mcp.json / opencode.json / .codex/* / .pi/*（四端配置）│
│  LICENSE / .gitignore                 │
├──────────────────────────────────────────────────────┤
│  Layer 2：来自后端模板（clone 后存在）                        │
│  这些文件来自独立模板 `github.com/kruily/go-ai-backend-template`，  │
│  由 bootstrap-project 在技术选型确认后 clone 到项目 backend/ 下    │
├──────────────────────────────────────────────────────┤
│  backend/go.mod                  module 占位 GOAI_MODULE│
│  backend/shared/pkg/             基础设施抽象接口（接口不动，实现按选型）│
│  backend/shared/utils/           工具包（直接可用）    │
│  backend/shared/goctl/           改造版 goctl 模板     │
│  backend/model/base/             BaseRepo 泛型 CRUD    │
│  backend/scripts/gen-*.sh        生成脚本              │
└──────────────────────────────────────────────────────┘

注意：
- Layer 2 来自独立模板仓库 `github.com/kruily/go-ai-backend-template`，
  不在 harness 层中预建。`go build` 验证在 clone 后做。
- GOAI_MODULE 是一个显式占位，clone 后用 sed 替换为实际 module 前缀。

---

## 其他文件

不在三层列表中的文件：

| 文件 | 归属 |
|------|------|
| `ROADMAP.md` | 工作室自身的开发计划，不进 agent 规范 |
| `PROJECT.md` | bootstrap 后由 AI 生成的项目设定文件，不存在于模板中 |
| `docs/` | bootstrap 后项目产出文档（如 tech-selection.md 等） |
| `deploy/` | 仅 deploy/README.md（约束说明），具体物料由 scaffold-deploy 技能生成 |
| `backend/` 下的模块目录 | bootstrap 按选型由 goctl 生成，代码来自 `go-ai-backend-template` |


> 如果你在 rules/ 中看到了 ROADMAP.md，知道它已移到仓库根。
> 现在 rules/ 只放 agent 可读的规范文档。

