# CLAUDE.md — Claude Code プロジェクト指示書

Claude Code はセッション開始時にこのファイルを自動的に読み込む。ここに記載されたルールはすべてのセッションで必ず遵守すること。

---

## プロジェクト概要

- **リポジトリ:** koma236/expense-tracker
- **目的:** 家計簿アプリ（学習用課題）
- **技術スタック:** フロントエンド: Nuxt 3 (Vue 3) / バックエンド: Go / データベース: MySQL
- **デプロイ:** AWS EC2（アプリ）/ RDS（MySQL）

---

## GitHub 開発フロー（必須ルール）

### 鉄則：作業を始める前に必ず GitHub Issue を作成すること

**すべて強制：**

1. **Issue ファースト** — 対応する Issue が存在しない状態で作業を開始してはならない
2. **ブランチ必須** — Issue ごとに専用ブランチを切ること
3. **main への直接プッシュ禁止** — すべての変更は Pull Request を通じてマージすること
4. **PR は Issue を参照すること** — PR 本文に `Closes #N` を必ず含めること
5. **マージはユーザーが実施** — Claude（AI）は PR 作成（Step 4）で停止すること。ユーザーの明示的な指示がない限り `gh pr merge` を実行してはならない。レビューとマージはユーザー本人が行う。

---

## ステップバイステップ ワークフロー

### Step 1: GitHub Issue を作成する

```bash
gh issue create \
  --title "作業内容を簡潔に記述" \
  --label "feature"  # feature / bug / chore / documentation
```

Issue 番号（例: `#5`）を必ず控えておくこと。

### Step 2: ブランチを作成する

**ブランチ命名規則（英語）：**

| 種別 | 形式 | 例 |
|------|------|----|
| 新機能 | `feature/issue-{N}-{slug}` | `feature/issue-5-add-expense-api` |
| バグ修正 | `fix/issue-{N}-{slug}` | `fix/issue-7-fix-null-pointer` |
| ドキュメント | `docs/issue-{N}-{slug}` | `docs/issue-9-update-readme` |
| その他 | `chore/issue-{N}-{slug}` | `chore/issue-11-update-deps` |

- `{N}` は Issue 番号（数字のみ）
- `{slug}` は英数字・ハイフンのみ、20文字以内推奨

```bash
git checkout -b feature/issue-{N}-{slug}
```

### Step 3: 実装・コミットする

コミットメッセージ形式: `<type>: <要約> (#N)`

- type は英語: `feat` / `fix` / `docs` / `chore` / `refactor` / `test`

```bash
git add <ファイル名>
git commit -m "feat: 支出登録APIを追加 (#5)"
```

### Step 4: Pull Request を作成する

```bash
gh pr create \
  --title "変更内容の要約 (#N)" \
  --body "$(cat .github/PULL_REQUEST_TEMPLATE.md)" \
  --base main
```

PR 本文の `Closes #` 欄に Issue 番号を必ず記入すること。

### Step 5: マージ（ユーザーが実施）

**Claude は Step 4（PR 作成）で停止する。** マージはユーザーがレビューしたうえで手動で行う。Claude はユーザーの明示的な指示がない限り、以下のコマンドを実行してはならない。

```bash
# ↓ ユーザー向けの手順（Claude は実行しない）
gh pr merge <PR番号> --squash --delete-branch
```

---

## ブランチ保護ルール（GitHub 側で設定済み）

ソロ開発のため、レビュー必須は設定していない。

- `main` への直接プッシュ: **禁止**（PR 経由のみ）
- force push: **禁止**
- ブランチ削除（main）: **禁止**

---

## コミットの共同作成者（Cursor / Claude）

GitHub のコントリビューターとして Cursor / Claude を表示させるため、`.githooks/prepare-commit-msg` フックで全コミットに以下のトレーラーを自動付与する。

```
Co-Authored-By: Cursor Agent <cursoragent@cursor.com>
Co-Authored-By: Claude <noreply@anthropic.com>
```

**クローンごとに一度だけ有効化が必要**（git フックは共有されないため）:

```bash
git config core.hooksPath .githooks
chmod +x .githooks/prepare-commit-msg   # 新規クローン時は git の実行ビットで既に付与済みの場合あり
```

- 既に同じトレーラーがある場合は重複追加しない（冪等）
- `squash` / `merge` 由来のメッセージには追記しない
- 共同作成者として表示はされるが、クリック可能なプロフィールにするには各メールが GitHub アカウントに紐づいている必要がある

---

## よくある間違いと対処法

| 間違い | 正しい対応 |
|--------|-----------|
| Issue なしで作業を始めた | 今すぐ Issue を作成し、ブランチを切り直す |
| main で作業してしまった | `git stash` → ブランチ作成 → `git stash pop` |
| PR に `Closes #N` がない | PR 本文を編集して追加する |

---

## プロジェクト固有コマンド

モノレポ構成。`backend/`（Go: chi + sqlc）、`frontend/`（Nuxt 3）、ローカル DB は docker-compose（MySQL）。
詳細は各ディレクトリの README（[backend/README.md](backend/README.md) / [frontend/README.md](frontend/README.md)）と [アーキテクチャ設計書](docs/architecture.md) を参照。

必要ツール: Go 1.22+ / Node.js 20+ / Docker / golang-migrate / sqlc
（macOS: `brew install go node golang-migrate sqlc`）

```bash
# データベース（ローカル, 無料）— リポジトリルートで
docker compose up -d

# バックエンド（Go）
cd backend
cp .env.example .env
go mod tidy
migrate -path db/migrations -database "mysql://app:app@tcp(127.0.0.1:3306)/expense_tracker" up   # マイグレーション
docker exec -i expense-tracker-mysql mysql --default-character-set=utf8mb4 -uapp -papp expense_tracker < db/seeds/0001_categories.sql  # 初期データ（日本語文字化け防止に utf8mb4 指定が必須）
go run ./cmd/server          # 起動 → http://localhost:8080/api/health
sqlc generate                # クエリ変更時にコード再生成
go build ./... && go vet ./...

# フロントエンド（Nuxt 3）
cd frontend
cp .env.example .env
npm install
npm run dev                  # http://localhost:3000
npm run build
```

---

## 注意事項

- `.env` ファイルはコミットしないこと（`.gitignore` 済み）
- ビルド成果物（`build/` `dist/` `node_modules/` 等）はコミットしないこと（`.gitignore` 済み）
- 技術スタックは確定済み（Nuxt 3 / Go / MySQL、デプロイは AWS EC2 + RDS）。要件定義は次ステップで実施する
