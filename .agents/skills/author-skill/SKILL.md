---
name: author-skill
description: 当需要为当前项目编写一个新的业务技能（放在项目的 skills/ 目录）时使用。指导按 Agent Skills 标准格式创建 SKILL.md 与 references/，保证四端（Claude Code、Codex、opencode、pi）通用。用于把项目特有的领域工作流沉淀成可复用技能。
---

# author-skill

你指导用户为**当前项目**编写业务技能（区别于工作室自带的工具型技能）。
业务技能沉淀项目特有的领域工作流，放项目的 `skills/` 目录。

## 使用时机

- 用户想把一段重复的领域流程固化成技能（如某内容生产流程、某审批流程）。
- 用户说"写一个技能""加个 skill"。

## 工具型技能 vs 业务技能

| | 工具型技能 | 业务技能 |
|--|-----------|---------|
| 位置 | `.agents/skills/`（工作室自带） | 项目的 `skills/`（项目自建） |
| 内容 | 开发流程（bootstrap、加 API 等） | 领域工作流（项目特有） |
| 归属 | 模板 | 各项目 |

## Agent Skills 标准格式（务必遵守，保证四端通用）

```
skills/{skill-name}/
├── SKILL.md              # 入口，含 frontmatter
└── references/           # 按需读取的细分规则（可选）
    └── *.md
```

frontmatter **只用严格子集**：

```yaml
---
name: {skill-name}        # 必须匹配目录名，小写字母/数字/连字符
description: 一句话说明「做什么 + 何时用」，1-1024 字。四端靠它匹配触发。
---
```

`SKILL.md` 正文建议 < 500 行；大块细则拆到 `references/*.md`，正文里说明"何时读哪个 reference"（渐进式披露）。

## 执行流程

1. **理解意图**：这个技能解决什么、何时触发、涉及哪些步骤。
2. **定 name**：小写连字符，与目录名一致。
3. **写 description**：说清"做什么 + 何时用"——这是四端触发匹配的关键，务必具体。
4. **写正文**：使用时机、核心原则、执行流程（分步）、完成后、禁止事项。用自然语言，不塞 JSON schema。
5. **拆 references**（若内容多）：领域模型、任务目录、质量规则等放 `references/`，正文引用。
6. **校验**：name 匹配目录名；description 非空；正文精简。

## 完成后

报告技能路径与触发方式（Claude 自动匹配或 `/name`；opencode `skill` 工具；pi `/skill:name`）。

## 禁止

- 不用 Claude 专属 frontmatter 字段（`context: fork` 等）于需要四端通用的业务技能。
- 不把工具由技能静态声明（工具由运行时注入）。
- 不把 prompt 模板塞进技能（那属于项目的 prompts/ 目录）。
