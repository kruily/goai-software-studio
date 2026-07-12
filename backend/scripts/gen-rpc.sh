#!/usr/bin/env bash
# gen-rpc.sh — 用工作室改造版 goctl 模板生成 go-zero zRPC 代码（微服务模块）。
#
# 用法：
#   ./scripts/gen-rpc.sh <proto_path> <output_dir>
# 例：
#   ./scripts/gen-rpc.sh services/user/rpc/user.proto services/user/rpc
#
# 约定：使用改造版模板 --home ./shared/goctl；生成后 pb 不手改。
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
cd "$BACKEND_DIR"

PROTO="${1:?用法: gen-rpc.sh <proto_path> <output_dir>}"
OUT="${2:?用法: gen-rpc.sh <proto_path> <output_dir>}"
HOME_TPL="./shared/goctl"

if ! command -v goctl >/dev/null 2>&1; then
  echo "错误：未安装 goctl。" >&2; exit 1
fi

echo "==> 生成 RPC 代码（模板：$HOME_TPL，输出：$OUT）"
goctl rpc protoc "$PROTO" \
  --go_out="$OUT" --go-grpc_out="$OUT" --zrpc_out="$OUT" \
  --home "$HOME_TPL" --style go_zero

echo "==> 完成。请在 $OUT/internal/logic/ 下实现业务逻辑。"
