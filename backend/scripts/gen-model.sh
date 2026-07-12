#!/usr/bin/env bash
# gen-model.sh — 用工作室改造版 goctl 模板从数据库生成 GORM/SQL 模型代码。
#
# 用法：
#   ./scripts/gen-model.sh <driver> <datasource> <table> <output_dir>
# 例（PostgreSQL）：
#   ./scripts/gen-model.sh pg "postgres://user:pass@localhost:5432/db?sslmode=disable" user model/user
# 例（MySQL）：
#   ./scripts/gen-model.sh mysql "user:pass@tcp(localhost:3306)/db" user model/user
#
# 约定：driver 需与 rules/tech-selection.md 的数据库选型一致；使用改造版模板。
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
cd "$BACKEND_DIR"

DRIVER="${1:?用法: gen-model.sh <pg|mysql> <datasource> <table> <output_dir>}"
DSN="${2:?缺少 datasource}"
TABLE="${3:?缺少 table}"
OUT="${4:?缺少 output_dir}"
HOME_TPL="./shared/goctl"

if ! command -v goctl >/dev/null 2>&1; then
  echo "错误：未安装 goctl。" >&2; exit 1
fi

echo "==> 生成 model（driver=$DRIVER，table=$TABLE，输出：$OUT）"
goctl model "$DRIVER" datasource \
  --url "$DSN" --table "$TABLE" \
  --dir "$OUT" --home "$HOME_TPL" --style go_zero

echo "==> 完成。记得在 model/migrate.go 注册新模型。"
