# frontend — ExpenseTracker

Nuxt 3 (Vue 3) で実装するフロントエンド。設計は [docs](../docs) を参照。

UI は Tailwind CSS、グラフは Chart.js（クライアント専用）で描画する。API 呼び出しは `composables/` のラッパ経由に統一している。

## 必要なツール

- Node.js 20 以上
- npm

## セットアップ

```bash
cp .env.example .env
npm install
npm run dev
# → http://localhost:3000
```

バックエンド（既定 `http://localhost:8080/api`）が起動している必要がある。未起動の場合は各画面でエラー表示になる。

### 環境変数（`.env`）

| 変数 | 既定値 | 説明 |
|------|--------|------|
| `NUXT_PUBLIC_API_BASE` | `http://localhost:8080/api` | バックエンド API のベース URL |

## スクリプト

```bash
npm run dev        # 開発サーバー（HMR）
npm run build      # 本番ビルド
npm run preview    # ビルド結果のプレビュー
```

## 画面

| パス | 画面 | 機能 |
|------|------|------|
| `/transactions` | 取引一覧（月・種別・カテゴリで絞り込み） | F-1 |
| `/transactions/new` | 取引登録 | F-1 |
| `/transactions/{id}` | 取引編集・削除 | F-1 |
| `/categories` | カテゴリ管理 | F-2 |
| `/dashboard` | サマリ・グラフ・予算進捗 | F-3 / F-4 |
| `/budgets` | 予算設定 | F-4 |

ヘッダーの月セレクタで選択した対象月（`year_month`）は `useMonth` で全画面に共有される。

## ディレクトリ

```text
pages/         画面（ファイルベースルーティング）
                 transactions/ categories/ budgets/ dashboard.vue
components/     UI コンポーネント
                 TransactionForm / SummaryCards / CategoryPieChart /
                 MonthlyTrendChart / BudgetProgress
composables/   共通ロジック
                 useApi（fetch ラッパ）/ useExpenseApi（型付き API クライアント）/
                 useMonth（対象月の共有状態）
layouts/       共通レイアウト（ヘッダー・月セレクタ・ナビ）
types/         API スキーマに対応する型定義
utils/         ユーティリティ（API エラー抽出など）
nuxt.config.ts 設定（runtimeConfig.public.apiBase）
```
