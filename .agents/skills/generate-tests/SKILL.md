---
name: generate-tests
description: 从 OpenSpec spec 中的 GIVEN/WHEN/THEN 场景和 Go 源码生成测试用例。生成 Go 单元测试（table-driven + testify）和集成测试（testcontainers-go），覆盖 happy path、error handling、edge cases。
---

# generate-tests

你从 OpenSpec 的 spec 场景、Go 源码的 interface/struct、以及 `testing-conventions.md` 规范出发，
生成测试代码。支持单测（logic 层）与集测（真实 DB/Redis）。

## 使用时机

- OpenSpec propose 产出 `spec.md` 后（含 GIVEN/WHEN/THEN 场景），需要为它生成测试。
- 新增/修改 Go interface 或 logic 函数后，需要补测试。
- 用户说"加个测试""跑测试""覆盖这个场景"。

## 测试规范（摘要）

详见 `rules/testing-conventions.md`。关键要点：

| 维度 | 规范 |
|------|------|
| 框架 | Go `testing` + `testify`（assert/require） |
| 单测 | table-driven + subtests `t.Run` |
| 集测 | testcontainers-go（PostgreSQL/Redis） |
| mock | `testify/mock` 轻量；复杂接口用 `gomock` |
| 上下文 | Go 1.24 `t.Context()` |
| 目录 | `internal/logic/*_test.go` 就近 |
| golden | 结构化输出用 `testdata/*.golden` |
| CI | `go test -race -count=1 -cover ./...` |

## 执行流程

### 1. 读 spec，提取场景

- 读 `openspec/changes/<name>/specs/` 或 `openspec/specs/<domain>/spec.md`。
- 从 spec 中提取：
  - `Requirement: ... SHALL ...` → 测试点：这条需求是否满足。
  - `Scenario: GIVEN ... WHEN ... THEN ...` → 具体测试用例。
- 自行补充边界：空值、超长、重复、状态机跳转、并发冲突。

### 2. 生成 Go 测试（logic 层）

- 用 **table-driven + subtests** 模式。
- 每个 table row = 一个场景：给输入名 + 输入值 + 期望值 + 期望错误。
- 测试文件放 `{module}/api/internal/logic/*_test.go`（与 logic 同包）。
- 示例骨架：
  ```go
  func TestGetUserProfile(t *testing.T) {
      tests := []struct {
          name    string
          userID  int64
          wantErr bool
          wantMsg string
      }{
          {name: "正常查询", userID: 1, wantErr: false},
          {name: "不存在的用户", userID: 999, wantErr: true, wantMsg: "user not found"},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              // arrange: mock DB / 准备数据
              // act: 调 logic
              // assert: assert.Equal / require.Error
          })
      }
  }
  ```

### 3. 生成集成测试（可选，按需）

- 用 `testcontainers-go` 启动真实 PostgreSQL/Redis。
- 测试文件用 `*_integration_test.go` 后缀（`go test -tags=integration`）。
- 每个集测文件有 `TestMain` 管理容器生命周期：
  ```go
  func TestMain(m *testing.M) {
      container, _ := testcontainers.PostgresContainer(context.Background())
      defer container.Terminate(context.Background())
      os.Exit(m.Run())
  }
  ```

### 4. 执行与确认

- 在 `backend/` 下运行：`go test -race -count=1 -cover ./...`。
- 检查覆盖率报告，识别未覆盖分支。
- 若有 error `expected: X, got: Y`→ 确认是测试用例问题还是生产代码问题。

## 完成后

报告：生成了几个测试文件、覆盖了哪些 spec 场景、构建与测试结果、剩余未覆盖区域。

## 禁止

- 不改生产代码（只生成测试文件）。
- 不生成不可运行/不编译的测试（生成后必须 `go build ./...` + `go test ./...` 验证）。
- 不过度依赖 mock 忽略真实环境差异（关键路径用集测而非 mock）。
