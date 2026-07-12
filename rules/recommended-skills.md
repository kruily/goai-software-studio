# 推荐技能清单（按 agent）

工作室**自带的工具型技能**在 `.agents/skills/`（见 AGENTS.md）。本文件列出**可选的外部社区/官方技能**，
按 agent 角色推荐，供你按需引入。

> ⚠️ **供应链与 license 提醒**：技能可能含可执行 `scripts/`。引入前请：
> 1. 读一遍要引入的 `SKILL.md` 与其 scripts；
> 2. 固定到具体 commit；
> 3. 保留来源仓库的 LICENSE，并在技能目录记录来源；
> 4. **不要**批量捆绑"上千技能"合集。
>
> Anthropic 的 **文档技能（docx/pdf/pptx/xlsx）是 source-available、不可再分发**——用 `/plugin install document-skills@anthropic-agent-skills` 安装，别拷进本仓库。

## 引入方式

技能是「拷贝文件夹」而非 pip/npm 安装。放到 `.agents/skills/{name}/`（四端通用）：

```bash
# 示例：从来源仓库拷单个技能目录到工作室
git clone --depth 1 https://github.com/<repo> /tmp/src
cp -R /tmp/src/skills/<skill-name> .agents/skills/<skill-name>
# 记录来源与 commit 到该目录的 SOURCE.md
```

per-agent 绑定：在 `.agents/agents/{agent}.md` 的 `skills:` 列表加上技能名（Claude 预加载；opencode 用 permission.skill）。

## 按 agent 推荐（均需自行校验后引入）

| Agent | 推荐外部技能 | 来源 | 说明 |
|-------|------------|------|------|
| project-manager | `doc-coauthoring` | anthropics/skills (Apache-2.0) | 结构化协作写 spec/proposal |
| | `requirements-clarity` | softaworks/agent-toolkit | 澄清欠明确的需求 |
| tech-lead | `c4-architecture`, `mermaid-diagrams` | softaworks/agent-toolkit | 架构文档与图 |
| | `backend-to-frontend-handoff-docs` | softaworks/agent-toolkit | 前后端接口契约 |
| backend-dev | `golang-code-style`, `golang-concurrency`, `golang-error-handling`, `golang-context` | samber/cc-skills-golang | 通用 Go 最佳实践（补充工作室后端技能） |
| | `golang-grpc`, `golang-database` | samber/cc-skills-golang | 微服务/数据相关 |
| frontend-dev | `webapp-testing` | anthropics/skills (Apache-2.0) | Playwright 前端验证 |
| | `openapi-to-typescript` | softaworks/agent-toolkit | 从后端 OpenAPI 生成 TS 客户端 |
| ui-designer | `frontend-design`, `theme-factory` | anthropics/skills (Apache-2.0) | 设计方向与主题 |
| | `brand-guidelines` | anthropics/skills | 需改写为客户品牌 |
| devops | `terraform-*` | hashicorp/agent-skills (官方) | 若用 Terraform/HCP |
| code-reviewer | `golang-security`, `golang-safety`, `golang-lint` | samber/cc-skills-golang | Go 审查视角 |
| | `static-analysis`, `differential-review` | trailofbits/skills (官方) | 安全审查（偏安全/加密） |

## 应自行编写（无良好现成品）

这些是工作室/项目特有的，已由工作室自带技能覆盖或应用 `author-skill` 编写：

- **go-zero 后端**：已由 `gozero-add-api` / `gorm-add-model` / `add-worker-task` 覆盖。
- **docker/nginx 部署**：devops 的 `deploy/` 工作流，特有，建议自行沉淀。
- **OpenSpec 提案 / 任务派发**：已由 `spec-driven` / `dispatch-dev` 覆盖。

## 来源仓库

- [anthropics/skills](https://github.com/anthropics/skills) — 官方，example skills 为 Apache-2.0。
- [samber/cc-skills-golang](https://github.com/samber/cc-skills-golang) — ~40 个 Go 技能。
- [softaworks/agent-toolkit](https://github.com/softaworks/agent-toolkit) — 工作流/图表/文档。
- [hashicorp/agent-skills](https://github.com/hashicorp/agent-skills) — 官方 Terraform/Packer。
- [trailofbits/skills](https://github.com/trailofbits/skills) — 官方安全审查。
- [VoltAgent/awesome-agent-skills](https://github.com/VoltAgent/awesome-agent-skills) — 发现索引。
