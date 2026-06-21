package service

import (
	"context"
	"math"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
)

// nearThreshold は予算の「接近」と判定する消化率（既定80%）。
const nearThreshold = 0.8

// SummaryService は集計（サマリ・カテゴリ別支出・月別推移）を担う。
type SummaryService struct {
	q *db.Queries
}

func NewSummaryService(q *db.Queries) *SummaryService {
	return &SummaryService{q: q}
}

// Totals は当月の収入・支出・収支差。
type Totals struct {
	Income  int64
	Expense int64
	Balance int64
}

// CategoryAmount はカテゴリ別の集計金額。
type CategoryAmount struct {
	CategoryID   int64
	CategoryName string
	Total        int64
}

// BudgetProgressItem は予算に対する消化状況の1件。
// 月全体予算は CategoryValid=false（CategoryID/CategoryName は無効）。
type BudgetProgressItem struct {
	CategoryID    int64
	CategoryName  string
	CategoryValid bool
	Budget        int64
	Spent         int64
	Ratio         float64
	Status        string // "under" / "near" / "over"
}

// SummaryResult は対象月のサマリ集計結果。
type SummaryResult struct {
	YearMonth         string
	Totals            Totals
	ExpenseByCategory []CategoryAmount
	BudgetProgress    []BudgetProgressItem
}

// TrendPoint は月別収支推移の1点。
type TrendPoint struct {
	YearMonth string
	Income    int64
	Expense   int64
	Balance   int64
}

// Summary は対象月（YYYY-MM。空なら当月）のサマリを返す。
func (s *SummaryService) Summary(ctx context.Context, yearMonth string) (SummaryResult, error) {
	start, end, err := monthRange(yearMonth)
	if err != nil {
		return SummaryResult{}, &ValidationError{Fields: []FieldError{{Field: "year_month", Message: "対象月は YYYY-MM 形式で指定してください"}}}
	}

	totalsRow, err := s.q.GetMonthlyTotals(ctx, db.GetMonthlyTotalsParams{StartDate: start, EndDate: end})
	if err != nil {
		return SummaryResult{}, err
	}

	rows, err := s.q.ListExpenseByCategory(ctx, db.ListExpenseByCategoryParams{StartDate: start, EndDate: end})
	if err != nil {
		return SummaryResult{}, err
	}
	byCategory := make([]CategoryAmount, 0, len(rows))
	for _, r := range rows {
		byCategory = append(byCategory, CategoryAmount{
			CategoryID:   r.CategoryID,
			CategoryName: r.CategoryName,
			Total:        r.Total,
		})
	}

	progress, err := s.budgetProgress(ctx, start.Format("2006-01"), totalsRow.Expense, byCategory)
	if err != nil {
		return SummaryResult{}, err
	}

	return SummaryResult{
		YearMonth: start.Format("2006-01"),
		Totals: Totals{
			Income:  totalsRow.Income,
			Expense: totalsRow.Expense,
			Balance: totalsRow.Income - totalsRow.Expense,
		},
		ExpenseByCategory: byCategory,
		BudgetProgress:    progress,
	}, nil
}

// budgetProgress は対象月の予算に対する消化状況を組み立てる。
// 予算未設定の場合は空（予算のあるカテゴリ／月全体のみ含める）。
func (s *SummaryService) budgetProgress(ctx context.Context, yearMonth string, totalExpense int64, byCategory []CategoryAmount) ([]BudgetProgressItem, error) {
	budgets, err := s.q.ListBudgets(ctx, yearMonth)
	if err != nil {
		return nil, err
	}

	// カテゴリ別支出を引きやすいよう map 化。
	spentByCategory := make(map[int64]int64, len(byCategory))
	for _, c := range byCategory {
		spentByCategory[c.CategoryID] = c.Total
	}

	items := make([]BudgetProgressItem, 0, len(budgets))
	for _, b := range budgets {
		var spent int64
		item := BudgetProgressItem{Budget: int64(b.Amount)}
		if b.CategoryID.Valid {
			item.CategoryID = b.CategoryID.Int64
			item.CategoryName = b.CategoryName.String
			item.CategoryValid = true
			spent = spentByCategory[b.CategoryID.Int64]
		} else {
			spent = totalExpense
		}
		item.Spent = spent
		item.Ratio = ratio(spent, item.Budget)
		item.Status = budgetStatus(item.Ratio)
		items = append(items, item)
	}
	return items, nil
}

// ratio は spent/budget を小数第2位で丸めて返す。budget<=0 は 0。
func ratio(spent, budget int64) float64 {
	if budget <= 0 {
		return 0
	}
	return math.Round(float64(spent)/float64(budget)*100) / 100
}

// budgetStatus は消化率から状態（under/near/over）を判定する。
func budgetStatus(r float64) string {
	switch {
	case r > 1.0:
		return "over"
	case r >= nearThreshold:
		return "near"
	default:
		return "under"
	}
}

// Trend は対象月を末尾に、過去 months ヶ月分の収支推移を古い順で返す。
// 取引が無い月も 0 埋めして連続させる。months が 1 未満なら 6 を用いる。
func (s *SummaryService) Trend(ctx context.Context, yearMonth string, months int) ([]TrendPoint, error) {
	if months < 1 {
		months = 6
	}
	if months > 24 {
		months = 24
	}

	lastStart, _, err := monthRange(yearMonth)
	if err != nil {
		return nil, &ValidationError{Fields: []FieldError{{Field: "year_month", Message: "対象月は YYYY-MM 形式で指定してください"}}}
	}

	// 範囲: (months-1) ヶ月前の月初 〜 対象月の月末。
	firstStart := lastStart.AddDate(0, -(months - 1), 0)
	rangeEnd := lastStart.AddDate(0, 1, -1)

	rows, err := s.q.ListMonthlyTrend(ctx, db.ListMonthlyTrendParams{StartDate: firstStart, EndDate: rangeEnd})
	if err != nil {
		return nil, err
	}
	byMonth := make(map[string]db.ListMonthlyTrendRow, len(rows))
	for _, r := range rows {
		byMonth[r.Ym] = r
	}

	points := make([]TrendPoint, 0, months)
	for i := 0; i < months; i++ {
		ym := firstStart.AddDate(0, i, 0).Format("2006-01")
		if r, ok := byMonth[ym]; ok {
			points = append(points, TrendPoint{
				YearMonth: ym,
				Income:    r.Income,
				Expense:   r.Expense,
				Balance:   r.Income - r.Expense,
			})
		} else {
			points = append(points, TrendPoint{YearMonth: ym})
		}
	}
	return points, nil
}
