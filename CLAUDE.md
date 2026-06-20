# CLAUDE.md — Claude Code プロジェクト指示書

Claude Code はセッション開始時にこのファイルを自動的に読み込む。ここに記載されたルールはすべてのセッションで必ず遵守すること。

---

## プロジェクト概要

- **リポジトリ:** koma236/expense-tracker
- **目的:** 家計簿アプリ（学習用課題）
- **技術スタック:** TBD（要件定義で確定する）

---

## GitHub 開発フロー（必須ルール）

### 鉄則：作業を始める前に必ず GitHub Issue を作成すること

**すべて強制：**

1. **Issue ファースト** — 対応する Issue が存在しない状態で作業を開始してはならない
2. **ブランチ必須** — Issue ごとに専用ブランチを切ること
3. **main への直接プッシュ禁止** — すべての変更は Pull Request を通じてマージすること
4. **PR は Issue を参照すること** — PR 本文に `Closes #N` を必ず含めること

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

### Step 5: マージする

```bash
gh pr merge <PR番号> --squash --delete-branch
```

---

## ブランチ保護ルール（GitHub 側で設定済み）

ソロ開発のため、レビュー必須は設定していない。

- `main` への直接プッシュ: **禁止**（PR 経由のみ）
- force push: **禁止**
- ブランチ削除（main）: **禁止**

---

## よくある間違いと対処法

| 間違い | 正しい対応 |
|--------|-----------|
| Issue なしで作業を始めた | 今すぐ Issue を作成し、ブランチを切り直す |
| main で作業してしまった | `git stash` → ブランチ作成 → `git stash pop` |
| PR に `Closes #N` がない | PR 本文を編集して追加する |

---

## プロジェクト固有コマンド

```text
TBD：技術スタック確定後に、ビルド・起動・テスト・ヘルスチェック等のコマンドを追記する。
```

---

## 注意事項

- `.env` ファイルはコミットしないこと（`.gitignore` 済み）
- ビルド成果物（`build/` `dist/` `node_modules/` 等）はコミットしないこと（`.gitignore` 済み）
- 技術スタック・要件は要件定義ステップで確定し、確定後に本ファイルとテンプレートを実スタックへ更新すること
