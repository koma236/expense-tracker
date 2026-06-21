package service

import (
	"context"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
)

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

// SummaryResult は対象月のサマリ集計結果。
type SummaryResult struct {
	YearMonth         string
	Totals            Totals
	ExpenseByCategory []CategoryAmount
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

	return SummaryResult{
		YearMonth: start.Format("2006-01"),
		Totals: Totals{
			Income:  totalsRow.Income,
			Expense: totalsRow.Expense,
			Balance: totalsRow.Income - totalsRow.Expense,
		},
		ExpenseByCategory: byCategory,
	}, nil
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
