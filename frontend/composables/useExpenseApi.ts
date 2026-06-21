import type {
  Budget,
  BudgetInput,
  Category,
  CategoryInput,
  Summary,
  Transaction,
  TransactionInput,
  Trend,
  TrendPoint,
  TxType,
} from '~/types'

// 取引・カテゴリ API の型付きラッパ。各画面はこれを経由して呼び出す。
export function useExpenseApi() {
  const api = useApi()

  return {
    // カテゴリ一覧（type 指定で絞り込み）
    listCategories: (type?: TxType) => {
      const qs = type ? `?type=${type}` : ''
      return api.get<{ categories: Category[] }>(`/categories${qs}`).then((r) => r.categories)
    },

    createCategory: (input: CategoryInput) => api.post<Category>('/categories', input),
    updateCategory: (id: number, input: CategoryInput) =>
      api.put<Category>(`/categories/${id}`, input),
    deleteCategory: (id: number) => api.delete<{ deleted: boolean }>(`/categories/${id}`),

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

    // 予算一覧（対象月）
    listBudgets: (yearMonth: string) =>
      api
        .get<{ budgets: Budget[] }>(`/budgets?year_month=${yearMonth}`)
        .then((r) => r.budgets),

    createBudget: (input: BudgetInput) => api.post<Budget>('/budgets', input),
    updateBudget: (id: number, amount: number) =>
      api.put<Budget>(`/budgets/${id}`, { amount }),
    deleteBudget: (id: number) => api.delete<{ deleted: boolean }>(`/budgets/${id}`),

    // 対象月のサマリ（収支合計・カテゴリ別支出）
    getSummary: (yearMonth: string) =>
      api.get<Summary>(`/summary?year_month=${yearMonth}`),

    // 月別収支推移（既定で対象月を末尾に過去6ヶ月）
    getTrend: (yearMonth: string, months = 6): Promise<TrendPoint[]> =>
      api
        .get<Trend>(`/summary/trend?year_month=${yearMonth}&months=${months}`)
        .then((r) => r.months),
  }
}
