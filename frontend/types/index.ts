// アプリ共通の型定義。バックエンド API（docs/api-spec.md）のスキーマに対応する。

export type TxType = 'income' | 'expense'

export interface CategoryRef {
  id: number
  name: string
}

export interface Category {
  id: number
  name: string
  type: TxType
  created_at: string
  updated_at: string
}

export interface Transaction {
  id: number
  occurred_on: string // YYYY-MM-DD
  amount: number
  type: TxType
  category: CategoryRef
  memo: string | null
  created_at: string
  updated_at: string
}

// カテゴリの作成/更新リクエスト本体。
export interface CategoryInput {
  name: string
  type: TxType
}

// 取引の作成/更新リクエスト本体。
export interface TransactionInput {
  occurred_on: string
  amount: number
  type: TxType
  category_id: number
  memo: string | null
}

// 集計サマリ（api-spec §5.1）。
export interface SummaryTotals {
  income: number
  expense: number
  balance: number
}

export interface ExpenseByCategory {
  category: CategoryRef
  total: number
}

export interface Summary {
  year_month: string
  totals: SummaryTotals
  expense_by_category: ExpenseByCategory[]
  // budget_progress は F-4 で利用。現状はバックエンドから空配列が返る。
  budget_progress: unknown[]
}

// 月別収支推移の1点。
export interface TrendPoint {
  year_month: string // YYYY-MM
  income: number
  expense: number
  balance: number
}

export interface Trend {
  months: TrendPoint[]
}

// API 共通エラー形式（api-spec §1.2）。
export interface ApiFieldError {
  field: string
  message: string
}

export interface ApiError {
  error: {
    code: string
    message: string
    details?: ApiFieldError[]
  }
}
