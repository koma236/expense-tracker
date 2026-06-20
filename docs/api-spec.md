# API 仕様書 — ExpenseTracker

[要件定義書](requirements.md)の「REST API 概要」を基に、各エンドポイントのリクエスト/レスポンススキーマ・ステータスコード・エラー形式・バリデーションを定義する。

- **ベースURL:** `/api`
- **形式:** リクエスト/レスポンスともに JSON（`Content-Type: application/json`）
- **文字コード:** UTF-8
- **金額:** 整数（円）
- **日付:** `YYYY-MM-DD`（取引日）／対象月は `YYYY-MM`
- **認証:** なし（単一ユーザー前提）

---

## 1. 共通仕様

### 1.1 ステータスコード

| コード | 用途 |
|--------|------|
| 200 OK | 取得・更新・削除の成功 |
| 201 Created | 新規作成の成功 |
| 400 Bad Request | バリデーションエラー、不正なパラメータ |
| 404 Not Found | リソースが存在しない |
| 409 Conflict | 一意制約違反（重複登録など） |
| 500 Internal Server Error | サーバー内部エラー |

### 1.2 エラーレスポンス形式

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力内容に誤りがあります",
    "details": [
      { "field": "amount", "message": "金額は1以上で入力してください" }
    ]
  }
}
```

- `code`: 機械可読なエラー種別（`VALIDATION_ERROR` / `NOT_FOUND` / `CONFLICT` / `INTERNAL_ERROR`）
- `details`: フィールド単位のエラー（バリデーション時のみ）

### 1.3 共通の日時フィールド

各リソースのレスポンスには `created_at` / `updated_at`（ISO 8601, 例 `2026-06-20T10:00:00+09:00`）を含む。

---

## 2. Transactions（取引）

### 2.1 一覧取得 — `GET /api/transactions`

**クエリパラメータ**

| 名前 | 必須 | 説明 |
|------|------|------|
| `year_month` | 任意 | `YYYY-MM`。指定時は当月で絞り込み。未指定時は当月をデフォルト |
| `type` | 任意 | `income` / `expense` で絞り込み |
| `category_id` | 任意 | カテゴリで絞り込み |

**レスポンス 200**

```json
{
  "transactions": [
    {
      "id": 12,
      "occurred_on": "2026-06-18",
      "amount": 1280,
      "type": "expense",
      "category": { "id": 1, "name": "食費" },
      "memo": "ランチ",
      "created_at": "2026-06-18T12:30:00+09:00",
      "updated_at": "2026-06-18T12:30:00+09:00"
    }
  ]
}
```

### 2.2 取得 — `GET /api/transactions/{id}`

- **200**: 単一の取引オブジェクト（2.1 の要素と同形）
- **404**: 該当なし

### 2.3 作成 — `POST /api/transactions`

**リクエスト**

```json
{
  "occurred_on": "2026-06-18",
  "amount": 1280,
  "type": "expense",
  "category_id": 1,
  "memo": "ランチ"
}
```

**バリデーション**

| フィールド | ルール |
|------------|--------|
| `occurred_on` | 必須・有効な日付（`YYYY-MM-DD`） |
| `amount` | 必須・整数・1以上 |
| `type` | 必須・`income` または `expense` |
| `category_id` | 必須・存在するカテゴリ・**`type` が一致**すること |
| `memo` | 任意・255文字以内 |

- **201**: 作成された取引
- **400**: バリデーションエラー（`type` とカテゴリ種別の不一致を含む）

### 2.4 更新 — `PUT /api/transactions/{id}`

- リクエスト・バリデーションは 2.3 と同じ（全項目指定）
- **200**: 更新後の取引 / **400** / **404**

### 2.5 削除 — `DELETE /api/transactions/{id}`

- **200**: `{ "deleted": true }` / **404**

---

## 3. Categories（カテゴリ）

### 3.1 一覧取得 — `GET /api/categories`

**クエリパラメータ**

| 名前 | 必須 | 説明 |
|------|------|------|
| `type` | 任意 | `income` / `expense` で絞り込み |

**レスポンス 200**

```json
{
  "categories": [
    { "id": 1, "name": "食費", "type": "expense",
      "created_at": "...", "updated_at": "..." }
  ]
}
```

### 3.2 作成 — `POST /api/categories`

**リクエスト**

```json
{ "name": "サブスク", "type": "expense" }
```

**バリデーション**

| フィールド | ルール |
|------------|--------|
| `name` | 必須・50文字以内・同一 `type` 内で重複不可 |
| `type` | 必須・`income` または `expense` |

- **201** / **400** / **409**（同一種別での名称重複）

### 3.3 更新 — `PUT /api/categories/{id}`

- リクエストは 3.2 と同じ
- **200** / **400** / **404** / **409**

### 3.4 削除 — `DELETE /api/categories/{id}`

- 紐づく取引または予算が存在する場合は削除不可
- **200**: `{ "deleted": true }` / **404** / **409**（使用中のため削除不可）

---

## 4. Budgets（予算）

### 4.1 一覧取得 — `GET /api/budgets`

**クエリパラメータ**

| 名前 | 必須 | 説明 |
|------|------|------|
| `year_month` | 任意 | `YYYY-MM`。未指定時は当月 |

**レスポンス 200**

```json
{
  "budgets": [
    { "id": 3, "category": { "id": 1, "name": "食費" },
      "year_month": "2026-06", "amount": 40000,
      "created_at": "...", "updated_at": "..." },
    { "id": 4, "category": null,
      "year_month": "2026-06", "amount": 150000,
      "created_at": "...", "updated_at": "..." }
  ]
}
```

- `category` が `null` の場合は「月全体予算」。

### 4.2 作成 — `POST /api/budgets`

**リクエスト**

```json
{ "category_id": 1, "year_month": "2026-06", "amount": 40000 }
```

- `category_id` を省略または `null` にすると月全体予算。

**バリデーション**

| フィールド | ルール |
|------------|--------|
| `category_id` | 任意・指定時は存在する支出カテゴリ |
| `year_month` | 必須・`YYYY-MM` 形式 |
| `amount` | 必須・整数・1以上 |
| （重複） | 同一 `(category_id, year_month)` の予算は1件まで（月全体含む） |

- **201** / **400** / **409**（重複）

### 4.3 更新 — `PUT /api/budgets/{id}`

- **200** / **400** / **404** / **409**

### 4.4 削除 — `DELETE /api/budgets/{id}`

- **200**: `{ "deleted": true }` / **404**

---

## 5. Summary（集計）

### 5.1 取得 — `GET /api/summary`

**クエリパラメータ**

| 名前 | 必須 | 説明 |
|------|------|------|
| `year_month` | 任意 | `YYYY-MM`。未指定時は当月 |

**レスポンス 200**

```json
{
  "year_month": "2026-06",
  "totals": {
    "income": 250000,
    "expense": 182300,
    "balance": 67700
  },
  "expense_by_category": [
    { "category": { "id": 1, "name": "食費" }, "total": 41200 },
    { "category": { "id": 3, "name": "交通費" }, "total": 12800 }
  ],
  "budget_progress": [
    {
      "category": { "id": 1, "name": "食費" },
      "budget": 40000,
      "spent": 41200,
      "ratio": 1.03,
      "status": "over"
    },
    {
      "category": null,
      "budget": 150000,
      "spent": 182300,
      "ratio": 1.22,
      "status": "over"
    }
  ]
}
```

- `budget_progress[].status`: 予算に対する状態
  - `under`（接近閾値未満） / `near`（接近: 既定80%以上） / `over`（超過: 100%超）
- `ratio`: `spent / budget`（小数）。予算未設定のカテゴリは `budget_progress` に含めない。

### 5.2 通知（F-4）の扱い

- 通知（接近/超過アラート）は専用エンドポイントを設けず、`/api/summary` の `budget_progress[].status` をフロントが解釈してアプリ内表示する。
- 接近閾値（既定80%）はフロント定数として持つ（将来は設定化を検討）。
