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

// 予算の消化状態（接近閾値未満 / 接近 / 超過）。
export type BudgetStatus = 'under' | 'near' | 'over'

// 予算進捗（api-spec §5.1）。category が null の場合は月全体予算。
export interface BudgetProgress {
  category: CategoryRef | null
  budget: number
  spent: number
  ratio: number
  status: BudgetStatus
}

export interface Summary {
  year_month: string
  totals: SummaryTotals
  expense_by_category: ExpenseByCategory[]
  budget_progress: BudgetProgress[]
}

// 予算（api-spec §4）。category が null の場合は月全体予算。
export interface Budget {
  id: number
  category: CategoryRef | null
  year_month: string
  amount: number
  created_at: string
  updated_at: string
}

// 予算の作成リクエスト本体。category_id 省略/null で月全体予算。
export interface BudgetInput {
  category_id?: number | null
  year_month: string
  amount: number
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
