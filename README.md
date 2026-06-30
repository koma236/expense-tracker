# ExpenseTracker

家計簿アプリケーション（学習用課題）。

収入・支出の記録、カテゴリ管理、月次の集計・グラフ表示、予算管理までを備えた家計簿アプリ。
学習目的のため、これまで使ったことのない技術を採用している。構成は **Nuxt 3 フロントエンド ↔ Go の REST API ↔ MySQL** の3層で、AWS（CloudFront + S3 + EC2 + RDS）にデプロイしている。

## デモ

収入・支出の登録から、月次サマリ・カテゴリ別グラフ・予算アラートまでの一連の操作デモ（本番環境）。

https://github.com/user-attachments/assets/8b7b885f-05fc-4cc4-8181-684482255665

## 主な機能

| ID | 機能 | 概要 |
|----|------|------|
| F-1 | 収支記録 | 収入・支出の登録／一覧／編集／削除（CRUD）。月・種別・カテゴリで絞り込み |
| F-2 | カテゴリ管理 | 収入／支出カテゴリの追加・編集・削除。使用中カテゴリは削除をブロック |
| F-3 | 集計・グラフ表示 | 当月サマリ、カテゴリ別支出（円グラフ）、月別収支推移（棒グラフ） |
| F-4 | 予算管理 | 月全体／カテゴリ別の予算設定と、消化率に応じた接近・超過アラート |

## 技術スタック

- フロントエンド: Nuxt 3 (Vue 3) + Tailwind CSS + Chart.js（本番は `ssr: false` の静的SPA）
- バックエンド: Go（[chi](https://github.com/go-chi/chi) ルーター + [sqlc](https://sqlc.dev/) による型安全な SQL）
- データベース: MySQL 8.0（ローカルは Docker、本番は AWS RDS）
- インフラ: AWS（CloudFront / S3 / EC2 / RDS）を Terraform で構築
- 配信: CloudFront を前段に置き、フロントは S3、API は EC2 へルーティング（HTTPS・同一オリジン）

## アーキテクチャ / デプロイ構成

```text
Browser ─HTTPS─▶ CloudFront (*.cloudfront.net)
                  ├─ /*     → S3（フロント静的SPA, 非公開・OAC 経由のみ）
                  └─ /api/* → EC2:80(nginx) → Go API :8080 (systemd) → RDS(MySQL 8.0)
```

- **フロント**は S3 に配置（完全非公開）、CloudFront(OAC) からのみ配信。SPA ルーティングは CloudFront Function で `/index.html` にフォールバック。
- **API** は EC2 上の Go（systemd 常駐）。nginx が `/api` を Go にリバースプロキシ。EC2 の 80 番は CloudFront からのみ許可（直接アクセス不可）。
- **DB** は RDS（MySQL 8.0、パブリックアクセス無効、EC2 のセキュリティグループからのみ 3306 許可）。
- フロントと API は CloudFront 上で同一オリジンになるため **CORS 不要**。公開URLは HTTPS。
- 学習用のためコストを抑え、**EC2 は使うときだけ起動する「都度起動」**運用。EC2 に Elastic IP を付与し、停止/起動しても公開URL（CloudFront）とオリジンは不変。

詳細・構築/デプロイ手順は以下を参照:

- インフラ構築（Terraform）: [terraform/README.md](terraform/README.md)
- アプリのデプロイ手順・運用: [docs/deploy.md](docs/deploy.md)
- 設計の背景: [アーキテクチャ設計書](docs/architecture.md)

## リポジトリ構成

```text
.
├── backend/          Go の REST API（→ backend/README.md）
├── frontend/         Nuxt 3 のフロントエンド（→ frontend/README.md）
├── terraform/        AWS インフラ（IaC）（→ terraform/README.md）
├── deploy/           デプロイ用ファイル（systemd / nginx / スクリプト）
├── docs/             設計・デプロイドキュメント
├── compose.yml       ローカル MySQL（Docker Compose）
└── CLAUDE.md         開発フロー・運用ルール
```

## クイックスタート（ローカル）

必要ツール: **Go 1.22+ / Node.js 20+ / Docker / golang-migrate / sqlc**
（macOS: `brew install go node golang-migrate sqlc`）

```bash
# 1. データベース（ローカル MySQL）をリポジトリルートで起動
docker compose up -d

# 2. バックエンド
cd backend
cp .env.example .env
go mod tidy
migrate -path db/migrations \
  -database "mysql://app:app@tcp(127.0.0.1:3306)/expense_tracker" up   # マイグレーション
docker exec -i expense-tracker-mysql \
  mysql --default-character-set=utf8mb4 -uapp -papp expense_tracker < db/seeds/0001_categories.sql  # 初期カテゴリ投入
go run ./cmd/server     # → http://localhost:8080/api/health

# 3. フロントエンド（別ターミナル）
cd frontend
cp .env.example .env
npm install
npm run dev             # → http://localhost:3000
```

> 初期データ投入時は、日本語が文字化けしないよう `--default-character-set=utf8mb4` を必ず付ける。

詳しい手順・各層の構成は [backend/README.md](backend/README.md) / [frontend/README.md](frontend/README.md) を参照。

## デプロイ（AWS）

Terraform でインフラを構築し、アプリを配置する。概要のみ示す（詳細は各ドキュメント参照）。

```bash
# 1. インフラ構築（S3 / CloudFront / EC2 / RDS / EIP）
cd terraform
cp terraform.tfvars.example terraform.tfvars   # ssh_cidr, db_password を設定
terraform init && terraform apply

# 2. バックエンドを EC2 へ配置（SSH して実施）
#    → docs/deploy.md の手順（マイグレーション・systemd・nginx）

# 3. フロントを S3 へ配信（リポジトリルートで）
./deploy/deploy-frontend.sh                     # build → S3 sync → CloudFront invalidation
```

- インフラの詳細: [terraform/README.md](terraform/README.md)
- アプリ配置・コスト方針・都度起動の運用: [docs/deploy.md](docs/deploy.md)

## ドキュメント

- [要件定義書](docs/requirements.md)
- [API 仕様書](docs/api-spec.md)
- [DB 設計書](docs/db-design.md)
- [画面設計書](docs/screen-spec.md)
- [アーキテクチャ設計書](docs/architecture.md)
- [デプロイ手順書](docs/deploy.md)

## 開発フロー

Issue ファースト・ブランチ必須・PR 経由マージを徹底している。詳細は [CLAUDE.md](CLAUDE.md) を参照。

1. 作業前に GitHub Issue を作成
2. `feature/issue-{N}-{slug}` などのブランチを切る
3. 実装・コミット（`<type>: <要約> (#N)`）
4. `main` への直接プッシュは禁止。Pull Request を作成し、本文に `Closes #N` を記載
5. レビューのうえマージ
