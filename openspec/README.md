# openspec/

本目录承载 **OpenSpec** 规范驱动开发的产物。OpenSpec（Fission-AI）把功能需求沉淀成
可评审的规范，再推进实现。是工作室生命周期第 ④ 步的载体，由 `spec-driven` 技能引导。

## 目录

```
openspec/
├── specs/            # 系统当前行为的真源（每次 archive 合并进来）
│   └── {capability}/spec.md
├── changes/          # 进行中的变更提案，一个功能一个目录
│   └── {change-name}/
│       ├── proposal.md   # 为什么做、改什么（project-manager 产出）
│       ├── design.md     # 技术方案（tech-lead 产出，遵循后端规范）
│       ├── tasks.md      # 可执行任务清单（标注可并行，供 dispatch-dev/tech-lead 派发）
│       └── specs/        # 该变更的 delta（增/改/删）
│       └── archive/      # 已完成变更归档
└── AGENTS.md         # 由 openspec init 生成的 agent 指引（勿手改）
```

## 安装与初始化

模板**不预装** OpenSpec（它是 npm 工具）。由 `bootstrap-project` 在第 ③ 步执行，或手动：

```bash
npm install -g @fission-ai/openspec@latest   # 需 Node 20.19+
openspec init                                 # 生成结构 + 各 agent 的 /opsx 命令
```

`openspec init` 会为已装的 agent（Claude Code / Codex / opencode 等）生成 `/opsx:*` 斜杠命令，
并写入 `openspec/AGENTS.md`。**这些命令与 AGENTS.md 由 OpenSpec 维护，不要手写或手改。**

## 工作循环（详见 spec-driven 技能）

```
/opsx:propose <功能>   # project-manager 产出 proposal，tech-lead 补 design/tasks
/opsx:apply            # tech-lead 派发 backend-dev/frontend-dev 实现
/opsx:archive          # 合并 delta 到 specs/，并回写根目录 PROJECT.md
```

## 与 PROJECT.md 的分工

- `PROJECT.md`（仓库根，由 bootstrap 生成）：全局产品蓝图，粗粒度、稳定。
- `openspec/`：功能级增量，细粒度、频繁。archive 时增量回流，两者互补。
