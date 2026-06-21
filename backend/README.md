# backend — ExpenseTracker API

Go（chi + sqlc）で実装する REST API。設計は [docs](../docs) を参照。

レイヤ構成は **handler（HTTP）→ service（業務ロジック・バリデーション）→ repository（sqlc 生成クエリ）**。
JSON のリクエスト／レスポンス形やエラー形式は [API 仕様書](../docs/api-spec.md) に準拠する。

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

### 環境変数（`.env`）

| 変数 | 既定値 | 説明 |
|------|--------|------|
| `APP_PORT` | `8080` | HTTP リッスンポート |
| `DB_HOST` / `DB_PORT` | `127.0.0.1` / `3306` | MySQL 接続先 |
| `DB_NAME` | `expense_tracker` | データベース名 |
| `DB_USER` / `DB_PASSWORD` | `app` / `app` | 接続ユーザー |
| `CORS_ALLOW_ORIGIN` | `http://localhost:3000` | CORS 許可オリジン（フロントエンド） |

## API エンドポイント

ベースパスは `/api`。詳細は [API 仕様書](../docs/api-spec.md) を参照。

| メソッド | パス | 説明 | 機能 |
|----------|------|------|------|
| GET | `/health` | ヘルスチェック（DB 疎通含む） | — |
| GET | `/categories` | カテゴリ一覧（`type` で絞り込み） | F-2 |
| POST | `/categories` | カテゴリ作成 | F-2 |
| PUT | `/categories/{id}` | カテゴリ更新 | F-2 |
| DELETE | `/categories/{id}` | カテゴリ削除（使用中は 409） | F-2 |
| GET | `/transactions` | 取引一覧（`year_month` / `type` / `category_id`） | F-1 |
| POST | `/transactions` | 取引作成 | F-1 |
| GET | `/transactions/{id}` | 取引取得 | F-1 |
| PUT | `/transactions/{id}` | 取引更新 | F-1 |
| DELETE | `/transactions/{id}` | 取引削除 | F-1 |
| GET | `/summary` | 当月サマリ（収支・カテゴリ別支出・予算進捗） | F-3 / F-4 |
| GET | `/summary/trend` | 月別収支推移 | F-3 |
| GET | `/budgets` | 予算一覧（`year_month`） | F-4 |
| POST | `/budgets` | 予算作成（重複は 409） | F-4 |
| PUT | `/budgets/{id}` | 予算（金額）更新 | F-4 |
| DELETE | `/budgets/{id}` | 予算削除 | F-4 |

## sqlc

`db/queries/*.sql` のクエリから型安全なコードを生成する。生成物は `internal/repository/db` に出力され、コミット対象とする。クエリを追加・変更したら再生成すること。

```bash
sqlc generate
```

## ビルド・検証

```bash
go build ./...
go vet ./...
```

## ディレクトリ

```text
cmd/server/          起動エントリポイント・ルーティング・DB 接続
internal/handler/    HTTP ハンドラ（リクエスト/レスポンス変換、エラー→HTTP 変換）
                     category / transaction / summary / budget / health / respond
internal/service/    業務ロジック・バリデーション・業務エラー
                     category / transaction / summary / budget
internal/repository/db/  sqlc 生成物（型安全なクエリ実装）
internal/model/      ドメイン型（必要に応じて追加）
db/migrations/       マイグレーション SQL（categories / transactions / budgets）
db/queries/          sqlc 入力クエリ
db/seeds/            初期データ（初期カテゴリ）
```
