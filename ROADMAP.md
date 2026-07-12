# 工作室能力清单与待办

## ✅ 已完整交付

### 规范层
- [x] AGENTS.md（唯一真源，含行为约束）
- [x] CLAUDE.md（@AGENTS.md 桥接）
- [x] opencode.json / .codex/config.toml / .pi/（四端配置）
- [x] rules/ × 12 份规范文档（lifecycle / backend-spec / frontend-spec / api-proto-conventions / agent-guide / code-intelligence / testing-conventions / template-boundary / service-split-patterns / frontend-backend-spec / tech-stack-catalog / recommended-skills）
- [x] 行为约束（最小改动、构建即验收、默认不跑测试）

### Agent 层
- [x] 8 个 agent 定义 × 5 端同步（40 文件）
- [x] 真源 `.agents/agents/` → 自动生成到 `.claude/` / `.codex/` / `.opencode/` / `.pi/`
- [x] sync-agents.py 同步脚本（置于 sync-agents 技能 scripts/ 下）

### Skills 层
- [x] 13 个技能：
  - 项目启动：bootstrap-project
  - 流程编排：spec-driven / dispatch-dev / sync-agents / author-skill
  - 后端开发：gozero-add-api / gorm-add-model / add-worker-task / add-infra-adapter
  - 前端：scaffold-frontend
  - 测试：generate-tests / e2e-runner / load-test

### 后端代码层
- [x] backend/shared/pkg/（storage / taskqueue / database / jwtx / redisx 接口 + 实现）
- [x] backend/shared/utils/（response / errorx / contexts / consts）
- [x] backend/shared/goctl/（60 个改造版模板）
- [x] backend/model/base/（BaseRepo 泛型 CRUD）
- [x] backend/scripts/（gen-api / gen-rpc / gen-model / setup-agents / sync-agents）
- [x] go build + go vet 通过

### 模板工程
- [x] .gitignore / LICENSE(MIT) / README.md
- [x] .mcp.json / deploy/ / openspec/ 骨架
- [x] GOAI_MODULE 占位
- [x] rules/template-boundary.md（三层模型 + 设计哲学）

---

## 🔶 已分析待办（按推荐优先级）

### P1 — 核心能力缺口

| # | 事项 | 说明 |
|---|------|------|
| 1 | **bootstrap-project 阶段 3 细粒度化** | 当前阶段 3（脚手架化）停留在"用 goctl 生成"的笼统描述，应拆成 AI 可逐条执行的清单（创建目录 → 写 .api → 调 gen-api.sh → 补模型 → 验证构建） |
| 2 | **test-engineer 与 dispatch-dev 的集成** | dispatch-dev 派发完成后，目前没有"自动触发 test-engineer"的机制。应在技能中写明"提交后 → code-reviewer → test-engineer"的编排逻辑 |
| 3 | **generate-tests 技能缺少 Go 测试框架的 task 模板** | 当前技能描述概念正确，但没有可以直接复用的测试 task 格式模板（table-driven 的标准样例在技能中引用了 testing-conventions.md，但那里也是描述性的） |

### P2 — 增强

| # | 事项 | 说明 |
|---|------|------|
| 4 | **deploy 的 nginx 配置模板** | deploy 现在是纯说明，但微服务场景下 nginx 是必要的。可以加一份 `deploy/nginx/backend.conf.example` |
| 5 | **CI 流水线模板** | `.github/workflows/ci.yml` 或类似，让项目克隆后直接有测试/构建/ lint 流水线 |
| 6 | **backend/model/custom_type/** | 自定义 JSON 类型的目录存在，但缺少一个 `json_struct.go` 示例（实现 GormDBDataType 按方言返回 jsonb/json/text） |
| 7 | **.env.example 移入 references** | 已完成（移入 scaffold-deploy/references/env-template.md） |

### P3 — 长期

| # | 事项 | 说明 |
|---|------|------|
| 8 | **opencode 的 command 对接** | 当前 opencode 通过 `.opencode/agents/` 读 agent，但 opencode 的 native command（`.opencode/command/*.md`）尚未对接技能 |
| 9 | ***_test.go 生成脚本** | `gen-test.sh` 或类似，自动化生成 logic 测试骨架 |
| 10 | **模块之间的 proto 定义管理** | 微服务下 proto 文件的共享/版本/生成流程 |
