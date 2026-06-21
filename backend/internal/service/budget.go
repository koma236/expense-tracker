package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
)

// BudgetInput は予算の作成/更新リクエスト。handler が JSON をデコードして渡す。
// CategoryID が nil の場合は「月全体」予算。
type BudgetInput struct {
	CategoryID *int64 `json:"category_id"`
	YearMonth  string `json:"year_month"`
	Amount     int32  `json:"amount"`
}

// BudgetService は予算のCRUDと業務ルールを担う。
type BudgetService struct {
	q *db.Queries
}

func NewBudgetService(q *db.Queries) *BudgetService {
	return &BudgetService{q: q}
}

// List は対象月（YYYY-MM。空なら当月）の予算一覧を返す。
func (s *BudgetService) List(ctx context.Context, yearMonth string) ([]db.ListBudgetsRow, error) {
	ym, err := normalizeYearMonth(yearMonth)
	if err != nil {
		return nil, &ValidationError{Fields: []FieldError{{Field: "year_month", Message: "対象月は YYYY-MM 形式で指定してください"}}}
	}
	return s.q.ListBudgets(ctx, ym)
}

// Create は予算を作成し、作成後の予算を返す。
func (s *BudgetService) Create(ctx context.Context, in BudgetInput) (db.GetBudgetRow, error) {
	ym, categoryID, err := s.validate(ctx, in)
	if err != nil {
		return db.GetBudgetRow{}, err
	}

	// 同一 (対象, 月) の重複チェック。
	var cnt int64
	if categoryID.Valid {
		cnt, err = s.q.CountCategoryBudget(ctx, db.CountCategoryBudgetParams{YearMonth: ym, CategoryID: categoryID})
	} else {
		cnt, err = s.q.CountMonthWideBudget(ctx, ym)
	}
	if err != nil {
		return db.GetBudgetRow{}, err
	}
	if cnt > 0 {
		return db.GetBudgetRow{}, &ConflictError{Message: "同じ対象・月の予算が既に設定されています"}
	}

	id, err := s.q.CreateBudget(ctx, db.CreateBudgetParams{
		CategoryID: categoryID,
		YearMonth:  ym,
		Amount:     in.Amount,
	})
	if err != nil {
		return db.GetBudgetRow{}, err
	}
	return s.q.GetBudget(ctx, id)
}

// Update は既存予算の金額を更新し、更新後の予算を返す。対象・月は変更しない。
func (s *BudgetService) Update(ctx context.Context, id int64, in BudgetInput) (db.GetBudgetRow, error) {
	if in.Amount < 1 {
		return db.GetBudgetRow{}, &ValidationError{Fields: []FieldError{{Field: "amount", Message: "金額は1以上で入力してください"}}}
	}

	rows, err := s.q.UpdateBudget(ctx, db.UpdateBudgetParams{Amount: in.Amount, ID: id})
	if err != nil {
		return db.GetBudgetRow{}, err
	}
	if rows == 0 {
		return db.GetBudgetRow{}, ErrNotFound
	}
	return s.q.GetBudget(ctx, id)
}

// Delete は予算を削除する。存在しなければ ErrNotFound。
func (s *BudgetService) Delete(ctx context.Context, id int64) error {
	rows, err := s.q.DeleteBudget(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// validate は予算入力を検証し、DB 用の値へ変換する。
// 戻り値は (year_month, category_id)。category_id は月全体予算なら Valid=false。
func (s *BudgetService) validate(ctx context.Context, in BudgetInput) (string, sql.NullInt64, error) {
	var fields []FieldError

	ym, err := normalizeYearMonth(in.YearMonth)
	if err != nil || in.YearMonth == "" {
		fields = append(fields, FieldError{Field: "year_month", Message: "対象月は YYYY-MM 形式で指定してください"})
	}

	if in.Amount < 1 {
		fields = append(fields, FieldError{Field: "amount", Message: "金額は1以上で入力してください"})
	}

	var categoryID sql.NullInt64
	if in.CategoryID != nil {
		cat, err := s.q.GetCategory(ctx, *in.CategoryID)
		if errors.Is(err, sql.ErrNoRows) {
			fields = append(fields, FieldError{Field: "category_id", Message: "指定されたカテゴリが存在しません"})
		} else if err != nil {
			return "", sql.NullInt64{}, err
		} else if cat.Type != db.CategoriesTypeExpense {
			fields = append(fields, FieldError{Field: "category_id", Message: "予算は支出カテゴリにのみ設定できます"})
		} else {
			categoryID = sql.NullInt64{Int64: cat.ID, Valid: true}
		}
	}

	if len(fields) > 0 {
		return "", sql.NullInt64{}, &ValidationError{Fields: fields}
	}
	return ym, categoryID, nil
}

// normalizeYearMonth は "YYYY-MM" を検証し正規化する。空なら当月を返す。
func normalizeYearMonth(yearMonth string) (string, error) {
	if yearMonth == "" {
		return time.Now().In(jst).Format("2006-01"), nil
	}
	t, err := time.ParseInLocation("2006-01", yearMonth, jst)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01"), nil
}
