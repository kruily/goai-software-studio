---
name: gorm-add-model
description: 当需要新增 GORM 数据模型时使用。在 model/{域}/ 建模型、按需定义 JSON 自定义类型、在 model/migrate.go 注册迁移、使用软删除。支持项目选定的任意 GORM 数据库（PostgreSQL/MySQL/SQLite 等），遵循工作室 GORM 规范。
---

# gorm-add-model

你为后端新增 GORM 模型。遵循工作室 GORM 规范（见 `rules/backend-spec.md`）。
数据库类型由项目选型决定（见 `rules/tech-selection.md`），本技能对各数据库通用。

## 使用时机

- 新增一张表 / 一个领域模型。
- 用户说"加个模型""建表""存储 xxx 数据"。

## 核心规范

- **模型位置**：`backend/model/{域}/`（顶层 model，不在 internal 内）。
- **连接经抽象**：DB 句柄来自 `shared/pkg/database`（按 `Config.Driver` 分发 pg/mysql/sqlite），model 只拿 `*gorm.DB`，不感知底层数据库。
- **结构化 JSON 字段**：定义专用结构体于 `model/custom_type/`，实现 `GormDataType()`、`GormDBDataType()`、`Value()`、`Scan()`；**不用 `map[string]interface{}` 或 `datatypes.JSON`**。
  - **跨数据库注意**：`GormDataType()` 返回逻辑类型；用 `GormDBDataType(db, field)` 按方言返回物理类型（PostgreSQL→`jsonb`，MySQL→`json`，SQLite→`text`），保证同一模型可迁移到不同数据库。
- **迁移注册**：新模型必须在 `model/migrate.go` 注册（统一迁移入口）。
- **软删除**：需要软删的表用 `gorm.DeletedAt`。
- **数据访问只在 model 层**：handler/logic 不直接拼 SQL。
- **可移植性**：避免特定数据库方言（如 PG 专有函数、MySQL 专有语法）；确需时在 `rules/tech-selection.md` 记录该表锁定的数据库。

## 执行流程

### 1. 建模型

- 在 `model/{域}/{name}.go` 定义 struct，字段带 gorm tag 与 json tag。
- 主键、时间戳（CreatedAt/UpdatedAt）、需要时 DeletedAt。
- 表名用 `TableName()` 显式指定（避免复数化歧义）。

### 2. 结构化 JSON 自定义类型（若有 JSON 字段）

- 在 `model/custom_type/{type}.go` 定义结构体，实现：
  - `GormDataType() string` → 返回逻辑类型 `"json"`
  - `GormDBDataType(db *gorm.DB, field *schema.Field) string` → 按方言返回物理类型：PostgreSQL `jsonb`、MySQL `json`、SQLite `text`
  - `Value() (driver.Value, error)` → json.Marshal
  - `Scan(value interface{}) error` → json.Unmarshal（兼容 []byte / string）
- 模型字段引用该类型。这样同一模型可在不同数据库间迁移。

### 3. 注册迁移

- 在 `model/migrate.go` 的迁移列表追加新模型，确保 `AutoMigrate` 覆盖。

### 4. 生成 CRUD（可选）

- 如需从已有表生成 model，用脚本（内部走改造版模板，driver 需与选型一致），从仓库根运行：
  ```bash
  backend/scripts/gen-model.sh pg "postgres://..." user model/user
  backend/scripts/gen-model.sh mysql "user:pass@tcp(...)/db" user model/user
  ```
- 否则手写领域数据访问方法（FindOne/List/Create/Update）。

### 5. 装配与验证

- 若 logic 需要，在 `ServiceContext` 注入 model。
- `go build ./...` 验证。

## 完成后

报告：新增的模型、JSONB 自定义类型（若有）、migrate 是否已注册、构建结果。

## 禁止

- 用 `map[string]interface{}` / `datatypes.JSON` 存结构化 JSON（用 custom_type 自定义类型）。
- 新模型不注册 migrate。
- 在 handler/logic 直接写 SQL 绕过 model 层。
- 无必要地使用特定数据库方言破坏可移植性（确需时记入 tech-selection.md）。
