# 代码智能工具

让 AI agent「看懂」代码结构的工具。通过 `.mcp.json`（四端共享）配置，由 `bootstrap-project` 按需启用。
**导航优先用这些语义工具，而非全仓 grep。**

## 三层

| 工具 | 定位 | 默认 | 依赖 |
|------|------|------|------|
| **gopls MCP** | Go 官方 LSP 的 MCP 模式：定义/引用/包 API/诊断 | ✅ 默认启用 | `go install golang.org/x/tools/gopls@latest`（v0.20+） |
| **Serena** | 符号级检索与编辑，跨语言（含前端） | 推荐（默认关） | `uv tool install -p 3.13 serena-agent` |
| **CodeGraphContext** | 调用/依赖图："谁调用了 X""改动影响范围" | 可选加装 | pip 安装 + 图数据库（较重） |

## gopls MCP（默认）

Go 语义导航，官方、零额外运行时、随代码变化保持正确。

```bash
go install golang.org/x/tools/gopls@latest
# 可选：导出用法指引给 agent 参考
# gopls mcp -instructions > rules/gopls-instructions.md
```

已在 `.mcp.json`、`opencode.json`、`.codex/config.toml`、`.pi/settings.json` 配置为 `gopls mcp`。

## Serena（推荐）

符号级检索 + 编辑 + 记忆，覆盖非 Go 部分（前端 TS 等）。

```bash
uv tool install -p 3.13 serena-agent
```

`.mcp.json` 中已配置（默认 `enabled:false`）；需要时启用。拉入 uv/Python 依赖，故默认关。

## CodeGraphContext（可选加装）

真正的代码图（函数/调用/依赖），适合大范围影响分析。pre-1.0 且需图数据库（Neo4j/FalkorDB），
**不做默认**，仅在确需调用图查询时按其文档加装。

## 何时用哪个

- 找定义/引用/包 API、看类型/诊断 → **gopls MCP**。
- 跨语言、符号级安全编辑、需要跨会话记忆 → **Serena**。
- "改这个函数会影响哪些地方""谁调用了它"的全局图查询 → **CodeGraphContext**。
