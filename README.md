# ExpenseTracker

家計簿アプリケーション（学習用課題）。

収入・支出の記録、カテゴリ管理、月次の集計・グラフ表示、予算管理までを備えた家計簿アプリ。
学習目的のため、これまで使ったことのない技術を採用している。構成は **Nuxt 3 フロントエンド ↔ Go の REST API ↔ MySQL** の3層。

## 主な機能

| ID | 機能 | 概要 |
|----|------|------|
| F-1 | 収支記録 | 収入・支出の登録／一覧／編集／削除（CRUD）。月・種別・カテゴリで絞り込み |
| F-2 | カテゴリ管理 | 収入／支出カテゴリの追加・編集・削除。使用中カテゴリは削除をブロック |
| F-3 | 集計・グラフ表示 | 当月サマリ、カテゴリ別支出（円グラフ）、月別収支推移（折れ線） |
| F-4 | 予算管理 | 月全体／カテゴリ別の予算設定と、消化率に応じた接近・超過アラート |

## 技術スタック

- フロントエンド: Nuxt 3 (Vue 3) + Tailwind CSS + Chart.js
- バックエンド: Go（[chi](https://github.com/go-chi/chi) ルーター + [sqlc](https://sqlc.dev/) による型安全な SQL）
- データベース: MySQL 8.0（ローカルは Docker、本番は AWS RDS）
- デプロイ: AWS EC2（アプリ）/ RDS（MySQL）

## リポジトリ構成

```text
.
├── backend/          Go の REST API（→ backend/README.md）
├── frontend/         Nuxt 3 のフロントエンド（→ frontend/README.md）
├── docs/             設計ドキュメント
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

## デプロイ

本番は **AWS EC2（アプリ）+ RDS（MySQL）** 構成。学習用課題のためコストを抑え、**EC2 は使うときだけ起動する「都度起動」運用**とする。EC2 上では Docker を使わず、Go はバイナリ＋systemd、Nuxt は静的SPA を nginx で配信し、`/api` を Go API へリバースプロキシする（同一オリジンのため CORS 不要）。

手順の詳細・コスト方針・起動/停止の運用は [デプロイ手順書](docs/deploy.md) を参照。デプロイ用ファイルは [deploy/](deploy/) にある。

## ドキュメント

- [要件定義書](docs/requirements.md)
- [API 仕様書](docs/api-spec.md)
- [DB 設計書](docs/db-design.md)
- [画面設計書](docs/screen-spec.md)
- [アーキテクチャ設計書](docs/architecture.md)

## 開発フロー

Issue ファースト・ブランチ必須・PR 経由マージを徹底している。詳細は [CLAUDE.md](CLAUDE.md) を参照。

1. 作業前に GitHub Issue を作成
2. `feature/issue-{N}-{slug}` などのブランチを切る
3. 実装・コミット（`<type>: <要約> (#N)`）
4. `main` への直接プッシュは禁止。Pull Request を作成し、本文に `Closes #N` を記載
5. レビューのうえマージ
