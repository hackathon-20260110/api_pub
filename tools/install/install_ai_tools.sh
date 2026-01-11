#! /bin/bash

set -euo pipefail
export DEBIAN_FRONTEND=noninteractive

# ---------------------------------------------
# install_ai_tools.sh
#
# 目的:
# - curl / node / npm を用意
# - Claude Code / Codex / Gemini CLI を導入（基本はnpmグローバル）
#
# 使い方:
#   bash install_ai_tools.sh
#
# 注意:
# - CLI導入コマンドはユーザー指定の固定コマンドのみを使用します（任意コマンド/候補探索はしません）
# ---------------------------------------------

log() { echo -e "[install-ai-tools] $*"; }
warn() { echo -e "[install-ai-tools][WARN] $*" >&2; }
die() { echo -e "[install-ai-tools][ERROR] $*" >&2; exit 1; }

have_cmd() { command -v "$1" >/dev/null 2>&1; }

as_root() {
  if [ "$(id -u)" -eq 0 ]; then
    "$@"
  elif have_cmd sudo; then
    sudo "$@"
  else
    die "root権限が必要です（sudoも見つかりません）。rootで実行してください。"
  fi
}

apt_update_once() {
  # 多重実行で遅くならないように、1回だけupdateする
  if [ -z "${_APT_UPDATED:-}" ]; then
    as_root apt-get update -y
    _APT_UPDATED=1
  fi
}

apt_install() {
  apt_update_once
  as_root apt-get install -y --no-install-recommends "$@"
}

ensure_cmd_or_install_apt() {
  local cmd="$1"; shift
  if have_cmd "$cmd"; then
    return 0
  fi
  apt_install "$@"
  have_cmd "$cmd" || die "コマンド '$cmd' の用意に失敗しました。"
}

ensure_base_deps() {
  if ! have_cmd apt-get; then
    die "このスクリプトは apt-get が使える環境（Debian/Ubuntu系）を想定しています。"
  fi

  # TLS/ダウンロード系の最低限
  ensure_cmd_or_install_apt curl curl ca-certificates
  ensure_cmd_or_install_apt git git

  # npmのネイティブ依存が必要になるケースに備えて
  apt_install build-essential python3
}

ensure_node_npm() {
  if have_cmd node && have_cmd npm; then
    log "node/npm は既にあります: node=$(node --version 2>/dev/null || true), npm=$(npm --version 2>/dev/null || true)"
    return 0
  fi

  log "node/npm を apt でインストールします..."
  apt_install nodejs npm

  have_cmd node || die "node のインストールに失敗しました。"
  have_cmd npm || die "npm のインストールに失敗しました。"
  log "node/npm をインストールしました: node=$(node --version), npm=$(npm --version)"
}

install_with_npm() {
  # usage: install_with_npm "tool name" "@scope/pkg"
  local tool="$1"
  local pkg="$2"
  log "$tool: npmでインストールします: npm install -g $pkg"
  npm install -g "$pkg"
}

print_summary() {
  log "インストール結果:"
  for bin in curl git node npm claude codex gemini; do
    if have_cmd "$bin"; then
      local p
      p="$(command -v "$bin")"
      log " - $bin: OK ($p)"
    else
      log " - $bin: MISSING"
    fi
  done
}

main() {
  log "開始: AIツールのセットアップ"
  ensure_base_deps
  ensure_node_npm
  install_with_npm "Codex" "@openai/codex"
  install_with_npm "Claude Code" "@anthropic-ai/claude-code"
  install_with_npm "Gemini CLI" "@google/gemini-cli"

  print_summary

  log "完了"
  log "補足: 各CLIは認証情報(APIキー等)の設定が別途必要です。"
}

main "$@"
