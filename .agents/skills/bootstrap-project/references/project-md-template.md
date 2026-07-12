# PROJECT.md 模板

bootstrap-project 在开项目时据此在**仓库根**创建 `PROJECT.md`。它是项目的全局产品蓝图（活文档），
每次功能迭代（OpenSpec archive）后由 spec-driven 回写更新。

与 openspec/ 的分工：
- `PROJECT.md`：粗粒度、稳定、少变的产品蓝图。
- `openspec/`：细粒度、频繁变化的功能级增量。

将下面内容写入仓库根 `PROJECT.md`，把 `{{PROJECT_NAME}}` 替换为项目名，逐项填充 TODO：

```markdown
# {{PROJECT_NAME}} — 项目设定文档

> 状态：草稿 · 最后更新：TODO · 维护：bootstrap-project / spec-driven 自动回写

## 一句话定位

TODO：这个产品是什么、为谁解决什么问题。

## 目标用户

- TODO：主要用户群体与场景。

## 核心功能清单

> 粗粒度功能列表。每项功能的细节在 `openspec/changes/` 中展开。

- [ ] TODO：功能 A
- [ ] TODO：功能 B

## 非功能约束

- 规模量级：TODO（并发/数据量/预期用户数）
- 实时性：TODO（是否需要 SSE/WebSocket 流式）
- 异步任务：TODO（是否需要队列/定时任务）
- AI 能力：TODO（是否含大模型/多模态，走哪种集成）

## 技术架构（由 bootstrap-project 写入）

| 维度 | 选型 |
|------|------|
| 架构 | TODO：单体 / 微服务 |
| module 前缀 | TODO：github.com/{organization}/{project} |
| 数据库 | TODO |
| 消息队列 | TODO |
| 对象存储 | TODO |
| 前端/客户端 | TODO |
| UI 设计 MCP | TODO |
| 代码智能工具 | gopls（默认） + TODO |

## 里程碑

- [ ] M1 — TODO
- [ ] M2 — TODO

## 变更记录

> 每次 OpenSpec archive 后由 spec-driven 追加一行。

- TODO：初始化
```
