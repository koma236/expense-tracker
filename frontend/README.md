# frontend — ExpenseTracker

Nuxt 3 (Vue 3) で実装するフロントエンド。設計は [docs](../docs) を参照。

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

起動後、トップページでバックエンド（`/api/health`）への疎通を確認できる。
バックエンドが未起動の場合はエラー表示になる。

## ディレクトリ

```text
pages/         画面（ファイルベースルーティング）
components/     UI コンポーネント
composables/   API クライアントなど共通ロジック（useApi.ts）
app.vue        ルートコンポーネント
nuxt.config.ts 設定（runtimeConfig.public.apiBase）
```
