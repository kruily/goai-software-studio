# pi agent 附加系统提示

本文件内容会**追加**到 pi 的默认系统提示之后。规范真源见仓库根 `AGENTS.md`（pi 会自动读取）。

- 遵循 `AGENTS.md` 的行为约束：最小改动、默认不跑测试、构建即验收、指令不清先澄清。
- 后端遵循 go-zero 约定：不手写 goctl 生成物，基础设施走 `backend/shared/pkg` 抽象。
- 开发生命周期见 `rules/lifecycle.md`。
