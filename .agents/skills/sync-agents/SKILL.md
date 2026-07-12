---
name: sync-agents
description: 当 AGENTS.md 规范真源、.mcp.json 或 .agents/agents/ 的 agent 定义变更后使用，把规范、MCP、专职 agent 定义同步/校验到四端（Claude Code、Codex、opencode、pi），确保四端读到一致内容。
---

# sync-agents

你负责保持工作室四端 agent 配置的一致性。真源有三处：`AGENTS.md`（规范）、`.mcp.json`（MCP）、
`.agents/agents/`（专职 agent 定义）。本技能确保其余入口与之不漂移。

## 使用时机

- 修改了 `AGENTS.md`（规范、技能清单、行为约束等）之后。
- 新增/修改 MCP 服务器（如 bootstrap 选定的设计 MCP）。
- 新增/修改 `.agents/agents/` 下的专职 agent 定义。
- 例行校验四端是否漂移。

## 桥接模型（务必遵守）

真源单一，各端只是不同入口，**不复制规范正文**：

| 端 | 入口 | 与真源的关系 |
|----|------|-------------|
| Claude Code | `CLAUDE.md` | 首行 `@AGENTS.md` 导入真源；下方仅放 Claude 专属 |
| Codex | `AGENTS.md`（原生） + `.codex/config.toml` | 自动读真源；toml 只放行为/MCP 设置 |
| opencode | `opencode.json` 的 `instructions` | 指向 AGENTS.md 等文件；不复制内容 |
| pi | `AGENTS.md`（原生） + `.pi/` | 自动读真源；APPEND_SYSTEM.md 仅追加要点 |
| 技能 | `.agents/skills/`（四端） | Claude 经 `.claude/skills` symlink |
| 专职 agent | `.agents/agents/`（真源） | 生成到各端（见下，格式不同不能 symlink） |
| MCP | `.mcp.json` | 各端配置引用同一批服务器 |

## 专职 agent 的四端生成（重点）

agent 定义没有跨端中立标准（不像 skills），格式各异，**必须从 `.agents/agents/{name}.md` 真源生成**：

| 端 | 目标路径 | 格式 | 转换 |
|----|---------|------|------|
| Claude Code | `.claude/agents/{name}.md` | MD+frontmatter | tools 首字母大写（read→Read）；保留 `skills:` 预加载 |
| opencode | `.opencode/agents/{name}.md` | MD+frontmatter | tools 小写；技能作用域用 `permission.skill` glob |
| pi | `.pi/agents/{name}.md` | MD+frontmatter | tools 小写；基本原样 |
| Codex | `.codex/agents/{name}.toml` | **TOML** | 正文 body → `developer_instructions`；`name`/`description`/`model` 转字段；无 skills 概念则并入指令 |

生成规则：
- **真源用小写 tools**（read/write/edit/grep/glob/bash），Claude 版转 Titlecase。
- **per-agent 技能绑定**：真源的 `skills:` 列表 →
  - Claude：`.claude/agents` 的 `skills:` frontmatter（预加载）。
  - opencode：`permission.skill` 白名单（如 `{"*":"deny","gozero-*":"allow"}`）。
  - Codex/pi：无一等字段，可在 developer_instructions/正文里说明该 agent 应使用哪些技能。
- 生成后逐个报告：写了哪些文件、做了哪些格式转换。

## 执行流程

1. **读真源**：`AGENTS.md`、`.mcp.json`。
2. **校验 Claude**：`CLAUDE.md` 首行必须是 `@AGENTS.md`；正文不得重复真源规范（只留 Claude 专属）。
3. **校验 opencode**：`opencode.json` 的 `instructions` 数组指向的文件都存在；`mcp` 块与 `.mcp.json` 的服务器集合一致（命令/参数对应）。
4. **校验 Codex**：`.codex/config.toml` 的 `mcp_servers` 与 `.mcp.json` 一致；不含与 AGENTS.md 重复的规范正文。
5. **校验 pi**：`.pi/settings.json` 的 mcp 与 `.mcp.json` 一致；`.pi/APPEND_SYSTEM.md` 只含指向真源的要点。
6. **校验 symlink**：`.claude/skills -> ../.agents/skills` 存在且有效。
7. **生成/校验专职 agent**：对 `.agents/agents/` 下每个 agent，确认四端目标文件存在且与真源一致（按上表转换）；漂移则重新生成。
   - 运行 `scripts/sync-agents.py` 一键执行全量同步：
   ```bash
   python3 .agents/skills/sync-agents/scripts/sync-agents.py
   ```
8. **修正漂移**：发现不一致时，**以 AGENTS.md / .mcp.json / .agents/agents 为准**修正各端入口文件，逐项报告改了什么。

## 完成后

报告：校验了哪些端、发现并修正了哪些漂移、当前是否一致。

## 禁止

- 不把 AGENTS.md 的规范正文复制进任何其他文件（只用导入/引用）。
- 不在同步时改动 AGENTS.md 本身的规范内容（真源变更应由用户或对应技能发起）。
- 不删除各端的专属设置（只对齐共享部分）。
