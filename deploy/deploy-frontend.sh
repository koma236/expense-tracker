#!/usr/bin/env bash
#
# フロントエンド（Nuxt 静的SPA）を S3 へ配信し、CloudFront のキャッシュを無効化する。
# ローカル（AWS CLI 認証済み）で実行する。バケット名・ディストリビューションIDは
# terraform output から取得する。
#
# 前提: terraform apply 済み（S3 / CloudFront 作成済み）。
# 使い方: リポジトリルートで  ./deploy/deploy-frontend.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TF_DIR="${REPO_ROOT}/terraform"

BUCKET="$(terraform -chdir="${TF_DIR}" output -raw s3_frontend_bucket)"
DIST_ID="$(terraform -chdir="${TF_DIR}" output -raw cloudfront_distribution_id)"
CF_URL="$(terraform -chdir="${TF_DIR}" output -raw cloudfront_url)"

echo "==> フロントビルド（静的SPA, /api 相対）"
cd "${REPO_ROOT}/frontend"
export NUXT_PUBLIC_API_BASE=/api
npm ci --no-audit --no-fund
npm run generate

echo "==> S3 へ同期: s3://${BUCKET}"
aws s3 sync .output/public/ "s3://${BUCKET}/" --delete
# index.html はキャッシュさせない（デプロイを即時反映）
aws s3 cp "s3://${BUCKET}/index.html" "s3://${BUCKET}/index.html" \
  --cache-control "no-cache, no-store, must-revalidate" \
  --content-type "text/html" --metadata-directive REPLACE

echo "==> CloudFront キャッシュ無効化: ${DIST_ID}"
aws cloudfront create-invalidation --distribution-id "${DIST_ID}" --paths "/*" >/dev/null

echo "==> 完了。公開URL: ${CF_URL}"
