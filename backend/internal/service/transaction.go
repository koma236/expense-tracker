package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
)

// jst は当アプリの基準タイムゾーン（Asia/Tokyo）。
var jst = mustLoadJST()

func mustLoadJST() *time.Location {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.FixedZone("JST", 9*60*60)
	}
	return loc
}

// ErrNotFound は対象リソースが存在しない場合に返す。
var ErrNotFound = errors.New("not found")

// FieldError はフィールド単位のバリデーションエラー。
type FieldError struct {
	Field   string
	Message string
}

// ValidationError は複数のフィールドエラーをまとめた業務エラー。
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string { return "入力内容に誤りがあります" }

// TransactionInput は取引の作成/更新リクエスト。handler が JSON をデコードして渡す。
type TransactionInput struct {
	OccurredOn string  `json:"occurred_on"`
	Amount     int32   `json:"amount"`
	Type       string  `json:"type"`
	CategoryID int64   `json:"category_id"`
	Memo       *string `json:"memo"`
}

// ListFilter は一覧取得の絞り込み条件。
type ListFilter struct {
	YearMonth  string // "YYYY-MM"。空なら当月。
	Type       string // "income" / "expense" / ""（全件）
	CategoryID int64  // 0 なら絞り込みなし
}

// TransactionService は取引のCRUDと業務ルールを担う。
type TransactionService struct {
	q *db.Queries
}

func NewTransactionService(q *db.Queries) *TransactionService {
	return &TransactionService{q: q}
}

// List は条件に合致する取引を新しい順で返す。
func (s *TransactionService) List(ctx context.Context, f ListFilter) ([]db.ListTransactionsRow, error) {
	start, end, err := monthRange(f.YearMonth)
	if err != nil {
		return nil, &ValidationError{Fields: []FieldError{{Field: "year_month", Message: "対象月は YYYY-MM 形式で指定してください"}}}
	}

	params := db.ListTransactionsParams{StartDate: start, EndDate: end}
	if f.Type == "income" || f.Type == "expense" {
		params.Type = db.NullTransactionsType{TransactionsType: db.TransactionsType(f.Type), Valid: true}
	}
	if f.CategoryID > 0 {
		params.CategoryID = sql.NullInt64{Int64: f.CategoryID, Valid: true}
	}
	return s.q.ListTransactions(ctx, params)
}

// Get は単一の取引を返す。存在しなければ ErrNotFound。
func (s *TransactionService) Get(ctx context.Context, id int64) (db.GetTransactionRow, error) {
	row, err := s.q.GetTransaction(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return db.GetTransactionRow{}, ErrNotFound
	}
	return row, err
}

// Create は取引を作成し、作成後の取引を返す。
func (s *TransactionService) Create(ctx context.Context, in TransactionInput) (db.GetTransactionRow, error) {
	params, err := s.validate(ctx, in)
	if err != nil {
		return db.GetTransactionRow{}, err
	}
	id, err := s.q.CreateTransaction(ctx, db.CreateTransactionParams{
		OccurredOn: params.OccurredOn,
		Amount:     params.Amount,
		Type:       params.Type,
		CategoryID: params.CategoryID,
		Memo:       params.Memo,
	})
	if err != nil {
		return db.GetTransactionRow{}, err
	}
	return s.q.GetTransaction(ctx, id)
}

// Update は既存の取引を更新し、更新後の取引を返す。
func (s *TransactionService) Update(ctx context.Context, id int64, in TransactionInput) (db.GetTransactionRow, error) {
	params, err := s.validate(ctx, in)
	if err != nil {
		return db.GetTransactionRow{}, err
	}
	rows, err := s.q.UpdateTransaction(ctx, db.UpdateTransactionParams{
		OccurredOn: params.OccurredOn,
		Amount:     params.Amount,
		Type:       params.Type,
		CategoryID: params.CategoryID,
		Memo:       params.Memo,
		ID:         id,
	})
	if err != nil {
		return db.GetTransactionRow{}, err
	}
	if rows == 0 {
		return db.GetTransactionRow{}, ErrNotFound
	}
	return s.q.GetTransaction(ctx, id)
}

// Delete は取引を削除する。存在しなければ ErrNotFound。
func (s *TransactionService) Delete(ctx context.Context, id int64) error {
	rows, err := s.q.DeleteTransaction(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// validatedTransaction はバリデーション済みの取引値。
type validatedTransaction struct {
	OccurredOn time.Time
	Amount     int32
	Type       db.TransactionsType
	CategoryID int64
	Memo       sql.NullString
}

// validate は取引入力を検証し、DB 用の値へ変換する。
func (s *TransactionService) validate(ctx context.Context, in TransactionInput) (validatedTransaction, error) {
	var fields []FieldError
	var out validatedTransaction

	occurred, err := time.ParseInLocation("2006-01-02", in.OccurredOn, jst)
	if err != nil {
		fields = append(fields, FieldError{Field: "occurred_on", Message: "取引日は YYYY-MM-DD 形式で入力してください"})
	} else {
		out.OccurredOn = occurred
	}

	if in.Amount < 1 {
		fields = append(fields, FieldError{Field: "amount", Message: "金額は1以上で入力してください"})
	} else {
		out.Amount = in.Amount
	}

	if in.Type != "income" && in.Type != "expense" {
		fields = append(fields, FieldError{Field: "type", Message: "種別は income または expense を指定してください"})
	} else {
		out.Type = db.TransactionsType(in.Type)
	}

	if in.Memo != nil {
		if len([]rune(*in.Memo)) > 255 {
			fields = append(fields, FieldError{Field: "memo", Message: "メモは255文字以内で入力してください"})
		} else if *in.Memo != "" {
			out.Memo = sql.NullString{String: *in.Memo, Valid: true}
		}
	}

	// カテゴリの存在チェックと種別一致チェック。
	if in.CategoryID <= 0 {
		fields = append(fields, FieldError{Field: "category_id", Message: "カテゴリを指定してください"})
	} else {
		cat, err := s.q.GetCategory(ctx, in.CategoryID)
		if errors.Is(err, sql.ErrNoRows) {
			fields = append(fields, FieldError{Field: "category_id", Message: "指定されたカテゴリが存在しません"})
		} else if err != nil {
			return validatedTransaction{}, err
		} else {
			out.CategoryID = cat.ID
			if in.Type == "income" || in.Type == "expense" {
				if string(cat.Type) != in.Type {
					fields = append(fields, FieldError{Field: "category_id", Message: "カテゴリの種別が取引種別と一致しません"})
				}
			}
		}
	}

	if len(fields) > 0 {
		return validatedTransaction{}, &ValidationError{Fields: fields}
	}
	return out, nil
}

// monthRange は "YYYY-MM" から当月の開始日・終了日（DATE境界）を返す。空なら当月。
func monthRange(yearMonth string) (time.Time, time.Time, error) {
	var first time.Time
	if yearMonth == "" {
		now := time.Now().In(jst)
		first = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, jst)
	} else {
		t, err := time.ParseInLocation("2006-01", yearMonth, jst)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		first = t
	}
	last := first.AddDate(0, 1, -1) // 当月末日
	return first, last, nil
}
