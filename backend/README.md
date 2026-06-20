# backend — ExpenseTracker API

Go (chi + sqlc) で実装する REST API。設計は [docs](../docs) を参照。

## 必要なツール

- Go 1.22 以上
- [golang-migrate](https://github.com/golang-migrate/migrate)（マイグレーション CLI）
- [sqlc](https://docs.sqlc.dev/)（SQL から型安全な Go コードを生成）
- Docker（ローカル MySQL 用）

macOS（Homebrew）の例:

```bash
brew install go golang-migrate sqlc
```

## セットアップ

```bash
# 1. ローカル MySQL を起動（リポジトリルートで）
cd .. && docker compose up -d && cd backend

# 2. 環境変数
cp .env.example .env

# 3. 依存解決
go mod tidy

# 4. マイグレーション（テーブル作成）
migrate -path db/migrations \
  -database "mysql://app:app@tcp(127.0.0.1:3306)/expense_tracker" up

# 5. 初期データ投入（シード）
#    日本語が文字化けしないよう --default-character-set=utf8mb4 を必ず付ける
docker exec -i expense-tracker-mysql \
  mysql --default-character-set=utf8mb4 -uapp -papp expense_tracker < db/seeds/0001_categories.sql

# 6. 起動
go run ./cmd/server
# → http://localhost:8080/api/health
```

## sqlc

`db/queries/*.sql` のクエリから型安全なコードを生成する。生成物は `internal/repository/db` に出力され、コミット対象とする。

```bash
sqlc generate
```

> 土台の段階では `ListCategories` のみ用意している。機能実装フェーズでクエリを追加していく。

## ディレクトリ

```text
cmd/server/        起動エントリポイント
internal/handler/  HTTP ハンドラ
internal/service/  業務ロジック（雛形）
internal/repository/ データアクセス（sqlc 生成物のラッパ）
internal/model/    ドメイン型（雛形）
db/migrations/     マイグレーション SQL
db/queries/        sqlc 入力クエリ
db/seeds/          初期データ
```
