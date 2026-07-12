#!/usr/bin/env bash
# setup-agents.sh — 安装工作室 agent 工具链：代码智能 MCP + OpenSpec。
#
# 用法：
#   ./scripts/setup-agents.sh
#
# 说明：
#   - 幂等：已装的跳过。
#   - Serena、CodeGraphContext 为可选，默认只提示不强装。
set -euo pipefail

echo "==> 检查 Go 工具链"
if ! command -v go >/dev/null 2>&1; then
  echo "缺少 go，请先安装 Go 1.26+" >&2; exit 1
fi

echo "==> 安装 gopls（代码智能，默认 MCP）"
if ! command -v gopls >/dev/null 2>&1; then
  go install golang.org/x/tools/gopls@latest
fi
gopls version 2>/dev/null || true

echo "==> 安装 goctl（go-zero 代码生成）"
if ! command -v goctl >/dev/null 2>&1; then
  go install github.com/zeromicro/go-zero/tools/goctl@latest
fi

echo "==> 检查 OpenSpec（必选：功能级 spec 驱动）"
if ! command -v openspec >/dev/null 2>&1; then
  if command -v npm >/dev/null 2>&1; then
    echo "   正在安装 OpenSpec（/opsx:propose/apply/archive 必需）..."
    npm install -g @fission-ai/openspec@latest || {
      echo "   OpenSpec 安装失败。请手动: npm i -g @fission-ai/openspec@latest（需 Node 20.19+）" >&2
    }
  else
    echo "   未检测到 npm。请安装 Node 20.19+ 后: npm i -g @fission-ai/openspec@latest" >&2
  fi
else
  echo "   OpenSpec 已安装"
fi

cat <<'EOF'

==> 基础工具就绪。可选项（按需手动安装）：
    - Serena（符号级检索/编辑，跨语言）：uv tool install -p 3.13 serena-agent
      然后在 .mcp.json 把 serena 的 enabled 打开。
    - CodeGraphContext（调用/依赖图）：见 rules/code-intelligence.md。

==> 下一步：在 AI agent 里让 project-manager 启动 bootstrap-project。
EOF
