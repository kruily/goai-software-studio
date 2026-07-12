#!/usr/bin/env bash
# gen-api.sh — 用工作室改造版 goctl 模板生成 go-zero REST API 代码。
#
# 用法（从仓库根执行）：
#   backend/scripts/gen-api.sh                      # 默认 user/api/desc/import.api（单体）
#   backend/scripts/gen-api.sh order/api/desc/import.api  # 其他模块
#
# 约定：
#   - 一律使用改造版模板 --home ./shared/goctl（走 response.Success/Error、透传错误码）。
#   - 生成后不手改 types.go / routes；业务逻辑写在 logic/。
#   - 由 gozero-add-api 技能调用；也可手动运行。
set -euo pipefail

# 切到 backend 根（脚本所在目录的上一级）。支持从仓库根或 backend/ 执行。
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
cd "$BACKEND_DIR"

API_DESC="${1:-user/api/desc/import.api}"
HOME_TPL="./shared/goctl"

# 解析 .api 文件所在模块的 api/ 目录（相对 backend/）
DESC_DIR="$(dirname "$API_DESC")"       # 例：user/api/desc
API_DIR="$(dirname "$DESC_DIR")"        # 例：user/api

if [[ ! -f "$API_DESC" ]]; then
  echo "错误：找不到 .api 入口：$API_DESC" >&2
  exit 1
fi
if ! command -v goctl >/dev/null 2>&1; then
  echo "错误：未安装 goctl。安装：go install github.com/zeromicro/go-zero/tools/goctl@latest" >&2
  exit 1
fi

# 生成目标目录 = .api 入口所在 dirname 的上一级（如 user/api/desc → user/api）
API_DIR="$(dirname "$DESC_DIR")"

echo "==> 格式化 .api（目录：$DESC_DIR）"
goctl api format --dir "$DESC_DIR"

echo "==> 生成 API 代码（模板：$HOME_TPL，输出：$API_DIR）"
goctl api go -api "$API_DESC" --dir "$API_DIR" --home "$HOME_TPL" --style go_zero

# 可选：结构体校验插件（若 .api 使用了 validate tag）。
if command -v goctl-validate >/dev/null 2>&1; then
  echo "==> 生成校验代码"
  goctl api plugin -p goctl-validate="validate --translator" --api "$API_DESC" --dir "$API_DIR"
fi

echo "==> 完成。请在 $API_DIR/internal/logic/ 下实现业务逻辑（生成物在 api/ 下）。"
