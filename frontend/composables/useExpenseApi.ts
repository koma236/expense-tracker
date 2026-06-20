import type { Category, Transaction, TransactionInput, TxType } from '~/types'

// 取引・カテゴリ API の型付きラッパ。各画面はこれを経由して呼び出す。
export function useExpenseApi() {
  const api = useApi()

  return {
    // カテゴリ一覧（type 指定で絞り込み）
    listCategories: (type?: TxType) => {
      const qs = type ? `?type=${type}` : ''
      return api.get<{ categories: Category[] }>(`/categories${qs}`).then((r) => r.categories)
    },

    // 取引一覧（対象月・種別・カテゴリで絞り込み）
    listTransactions: (params: { yearMonth: string; type?: TxType; categoryId?: number }) => {
      const q = new URLSearchParams({ year_month: params.yearMonth })
      if (params.type) q.set('type', params.type)
      if (params.categoryId) q.set('category_id', String(params.categoryId))
      return api
        .get<{ transactions: Transaction[] }>(`/transactions?${q.toString()}`)
        .then((r) => r.transactions)
    },

    getTransaction: (id: number) => api.get<Transaction>(`/transactions/${id}`),
    createTransaction: (input: TransactionInput) => api.post<Transaction>('/transactions', input),
    updateTransaction: (id: number, input: TransactionInput) =>
      api.put<Transaction>(`/transactions/${id}`, input),
    deleteTransaction: (id: number) => api.delete<{ deleted: boolean }>(`/transactions/${id}`),
  }
}
