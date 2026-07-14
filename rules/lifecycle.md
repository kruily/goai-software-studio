# 开发生命周期

本工作室把一个想法推进到上线代码，走固定的六步。每步有对应的 agent、产物与评审门。

```
① 想法
   ↓
② 项目设定
   ├── ②a 需求访谈（studio 技能）→ PROJECT.md  ← Gate: 需求确认
   └── ②b 管道铺设（studio 技能）→ openspec init + 四端配置  ← Gate: 管道就绪
   ↓  ← Gate: 项目就绪（PROJECT.md 已确认）
③ 功能规划（per feature，OpenSpec propose）
   ├── ③a PM 需求细化 → proposal.md  ← Gate: PRD 签批
   ├── ③b UI Designer 设计 → 线框/高保真/原型  ← Gate: 设计签批
   ├── ③c Tech Lead 技术架构（含选型）→ design.md + API 契约 + tech-selection.md  ← Gate: API 冻结
   └── ③d Tech Lead 任务拆分 → tasks.md  ← Gate: Sprint 承诺
   ↓  ← Gate: 规划完成
④ 实现 + 测试（OpenSpec apply）
   ├── ④a Frontend Dev + Backend Dev 并行开发
   ├── ④b Code Reviewer 每次 PR 审查  ← Gate: CR 通过
   ├── ④c Test Engineer 按 spec 场景持续产出测试
   └── ④d PM UAT 验收  ← Gate: UAT 通过
   ↓  ← Gate: 实现完成
⑤ 部署 + 发布
   ├── ⑤a DevOps 构建部署 + CI
   ├── ⑤b Test Engineer 部署后测试（E2E + 负载）
   └── ⑤c PM + DevOps 发布决策  ← Gate: 发布
   ↓
⑥ 复盘 + Archive
   ├── ⑥a Tech Lead OpenSpec archive → specs 合并
   ├── ⑥b PM 回写 PROJECT.md + 变更记录
   └── ⑥c 团队回顾
   ↓（循环：下一个功能，回到 ③）
```

---

## ① 想法

一句话或一段话说明要做什么产品、解决什么问题、大致规模。信息不必完整——第 ② 步会补齐。

**Agent:** 用户 + PM（project-manager）
**Gate:** 无（想法清晰即可进入 setup）

---

## ② 项目设定（Project Setup）

拆为两步，中间有确认门。

**②a 需求访谈（studio 技能）：** 与用户纯聊需求，不讨论技术。产出 `PROJECT.md`。
**Gate: 需求确认。** 用户确认 PROJECT.md 内容可接受后才能进入功能规划。

**②b 管道铺设（studio 技能）：** 初始化 openspec、写四端配置、回写 PROJECT.md 时间戳。**不做技术选型、不创建代码、不推荐架构。** 技术选型由 tech-lead 在功能规划阶段（③c）完成。
**Gate: 管道就绪。** openspec 目录已建、四端配置已同步。确认后进入功能规划。

**Agents:** PM（project-manager）
**Gate: 项目就绪。** `PROJECT.md` 已确认、管道已铺设。确认后进入功能规划。

---

## ③ 功能规划（Feature Planning，per feature）

每个功能是一个 OpenSpec change。用 `spec-driven` 技能引导。

### ③a 需求细化（PM）

PM 读 PROJECT.md，与用户细化该功能的需求：
- 功能定位、边界
- 验收标准（GIVEN/WHEN/THEN）
- 成功指标

产出 `openspec/changes/<name>/proposal.md`。

**Gate: PRD 签批（PM）。** proposal.md 内容用户认可后进入设计。

### ③b UI/UX 设计（UI Designer）

UI Designer 用设计 MCP（Magic / Figma / shadcn / Ardot）产出：
- **线框图**（低保真，确认布局和信息架构）
- **高保真设计稿**（像素级，含颜色/字体/间距）
- **交互原型**（可点击的流程演示）
- **状态矩阵**（每个组件的 loading / empty / error / edge case 状态）

设计稿可通过 Figma 链接或导出的截图交付。
前端/后端实现**不等待设计全部完成**——架构和设计可以并行推进。

**Gate: 设计签批（PM + Tech Lead + UI Designer）。** 所有 screens 设计完成、状态覆盖完整，确认可实现。

### ③c 技术架构（Tech Lead）

Tech Lead 做技术选型推荐与确认，并产出：
- `rules/tech-selection.md`（技术选型记录：架构、数据库、队列、存储、前端、设计 MCP、module 前缀）
- `design.md`（系统的技术方案，含模块划分、组件交互、数据流）
- API 契约（`.api` / `.proto` 定义，遵循 gozero-add-api 规范）
- 数据模型（GORM model 定义，不嵌套 BaseModel）
- 重要架构决策记录（可写入 `design.md` 决策章节）

**Gate: API 冻结（Tech Lead）。** 契约确定后前后端可并行开发。

### ③d 任务拆分（Tech Lead）

Tech Lead 将设计拆解为可执行的任务：
- 标注依赖关系
- 标注可并行项
- 预估相对大小

产出 `tasks.md`。

**Gate: Sprint 承诺（PM + Tech Lead）。** 团队认可当前 scope 可实现。

---

## ④ 实现 + 测试（OpenSpec apply）

前后端并行实现。

### ④a 实现

**Backend Dev** 按 `tasks.md` 和 `.api` 契约，调用 `gozero-add-api` / `gorm-add-model` / `add-worker-task` / `add-infra-adapter` 技能实现后端逻辑。

**Frontend Dev** 按设计稿和 API 契约，调用 `scaffold-frontend` 技能初始化前端，实现页面与组件，对接 API。

**并行前提：** API 契约已冻结（③c 门通过），设计稿已可用（③b 门通过）。

### ④b Code Review

**Code Reviewer** 审查每一个 PR：
- 分层边界（handler 不写业务逻辑）
- 基础设施抽象（不走直连 SDK）
- 统一响应体不入 .api
- 结构化 JSON 用自定义类型
- 无硬编码密钥

**Gate: Code Review 通过。** 每个 PR 必须通过 code-reviewer 审查。架构级变更需 Tech Lead 额外确认。

### ④c 持续测试（Test Engineer）

**Test Engineer** 在 spec propose 后就介入，按 proposal.md 的 GIVEN/WHEN/THEN 场景：
- 生成 Go 单测/集测（`generate-tests`）
- 生成 E2E 测试（`e2e-runner`）
- 按需生成负载测试（`load-test`）

测试随实现持续产出和执行，不等到实现全部完成才测。
Bug 报告直接反馈给对应 developer agent，修复后重新验证。

### ④d UAT（PM）

PM 验证实现是否符合 proposal.md 的验收标准。
可运行交付物（或 staging 环境）进行验收。

**Gate: UAT 通过（PM）。** 所有 P0/P1 bug 关闭，验收标准满足。

---

## ⑤ 部署 + 发布

### ⑤a 构建部署（DevOps）

DevOps 调用 `scaffold-deploy` 技能生成部署物料：
- Dockerfile（多阶段构建）
- docker-compose 编排
- nginx 转发配置（微服务时）
- CI 流水线

构建镜像，部署到 staging。

### ⑤b 部署后测试（Test Engineer）

- 冒烟测试：关键链路是否正常
- E2E 测试：核心用户流程
- 负载测试（按需）：验证性能阈值

### ⑤c 发布决策（PM + DevOps）

确认条件：
- UAT 通过
- 所有 P0/P1 bug 关闭
- 部署后测试通过
- 监控仪表盘就绪
- 回滚预案确认

**Gate: 发布门（PM + DevOps）。** 满足条件即发布。

---

## ⑥ 复盘 + Archive

### ⑥a Archive（Tech Lead）

```bash
/opsx:archive
```
- 本次功能的 delta specs 合并到 `openspec/specs/`（系统行为真源）
- `openspec/changes/<name>/` 标记为已归档

### ⑥b 回写 PROJECT.md（PM）

- 勾选完成的功能项
- 追加一行变更记录
- 更新最后更新时间

### ⑥c 回顾（团队）

- 本次迭代哪些做得好
- 哪些可以改进
- 下一轮开始前调整流程

**循环：** 下一个功能回到第 ③ 步。

---

## 产物对照

| 产物 | 粒度 | 阶段 | 何时更新 |
|------|------|------|----------|
| `PROJECT.md` | 全局产品蓝图 | ② | studio 创建，每次 archive 回写 |
| `rules/tech-selection.md` | 技术选型记录 | ② | studio 管道铺设时创建，选型变更时更新 |
| `openspec/changes/<name>/proposal.md` | 单功能需求 | ③a | 每次 propose |
| `openspec/changes/<name>/design.md` | 技术方案 | ③c | 每次 propose |
| `openspec/changes/<name>/tasks.md` | 任务清单 | ③d | 每次 propose |
| `openspec/specs/` | 系统当前行为真源 | ③-⑥ | 每次 archive |
| 代码 | 实现 | ④ | apply / 派发开发 |
| 测试代码 | 验证 | ④c | 持续 |
| deploy 物料 | 部署 | ⑤ | 发布前 |
| 回顾笔记 | 改进 | ⑥ | 每次迭代结束 |

## 评审门汇总

| 门 | 条件 | 签字方 | 进入 |
|----|------|-------|------|
| 项目设定 | PROJECT.md + tech-selection.md 已确认 | 用户 | ③ |
| PRD 签批 | proposal.md 内容已认可 | PM | ③b |
| 设计签批 | 所有 UI 状态覆盖+可实施 | PM + Tech Lead + UI Designer | ③c |
| API 冻结 | API 契约已定义 | Tech Lead | ④ |
| Sprint 承诺 | tasks scope 团队认可 | PM + Tech Lead | ④ |
| CR 通过 | 每 PR 审查通过 | Code Reviewer | ④d |
| UAT 通过 | 验收标准满足 | PM | ⑤ |
| 发布门 | 测试通过+监控就绪+预案确认 | PM + DevOps | ⑥ |

## 技能与阶段对照

| 技能 | 阶段 |
|------|------|
| `studio` | ② |
| `spec-driven` | ③④⑤⑥（贯穿）|
| `dispatch-dev` | ④a |
| `gozero-add-api` / `gorm-add-model` / `add-worker-task` / `add-infra-adapter` | ④a |
| `code-reviewer`（agent 非 skill） | ④b |
| `generate-tests` / `e2e-runner` / `load-test` | ④c⑤b |
| `scaffold-frontend` | ④a |
| `scaffold-deploy` | ⑤a |
