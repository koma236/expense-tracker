#!/usr/bin/env bash
#
# ExpenseTracker 再デプロイスクリプト（EC2 上で実行）。
#
# 役割: 最新コードを取得し、バックエンド（Go バイナリ）とフロントエンド（Nuxt 静的SPA）を
#       ビルドして配置し、サービスを再起動する。初回セットアップ後の更新作業用。
#
# 前提:
#   - リポジトリが REPO_DIR にクローン済み
#   - Go / Node.js 20 / nginx / systemd ユニット（expense-tracker-api）が導入・有効化済み
#   - /etc/expense-tracker/api.env が設定済み
#   詳細は docs/deploy.md を参照。
#
# 使い方:
#   cd /opt/expense-tracker && ./deploy/deploy.sh

set -euo pipefail

# ---- 設定（必要に応じて環境変数で上書き可能）---------------------------------
REPO_DIR="${REPO_DIR:-/opt/expense-tracker}"
WEB_ROOT="${WEB_ROOT:-/var/www/expense-tracker}"
BIN_PATH="${BIN_PATH:-/usr/local/bin/expense-tracker-api}"
API_SERVICE="${API_SERVICE:-expense-tracker-api}"
# フロントは同一オリジンの nginx 経由で /api を叩く
export NUXT_PUBLIC_API_BASE="${NUXT_PUBLIC_API_BASE:-/api}"

echo "==> リポジトリを更新: ${REPO_DIR}"
cd "${REPO_DIR}"
git pull --ff-only

# ---- バックエンド（Go）-------------------------------------------------------
echo "==> バックエンドをビルド"
cd "${REPO_DIR}/backend"
go build -o /tmp/expense-tracker-api ./cmd/server
sudo install -m 0755 /tmp/expense-tracker-api "${BIN_PATH}"
rm -f /tmp/expense-tracker-api

echo "==> API サービスを再起動: ${API_SERVICE}"
sudo systemctl restart "${API_SERVICE}"

# ---- フロントエンド（Nuxt 静的SPA）------------------------------------------
echo "==> フロントエンドをビルド（nuxt generate）"
cd "${REPO_DIR}/frontend"
npm ci
npm run generate

echo "==> 静的ファイルを配置: ${WEB_ROOT}"
sudo mkdir -p "${WEB_ROOT}"
sudo rsync -a --delete "${REPO_DIR}/frontend/.output/public/" "${WEB_ROOT}/"

echo "==> nginx をリロード"
sudo nginx -t
sudo systemctl reload nginx

echo "==> 完了。動作確認: curl -fsS http://localhost/api/health"
