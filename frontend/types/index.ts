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

// 取引の作成/更新リクエスト本体。
export interface TransactionInput {
  occurred_on: string
  amount: number
  type: TxType
  category_id: number
  memo: string | null
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
