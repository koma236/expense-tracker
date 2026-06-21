package service

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
)

// ConflictError は一意制約・使用中などで操作を実行できない場合に返す（HTTP 409）。
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string { return e.Message }

// CategoryInput はカテゴリの作成/更新リクエスト。handler が JSON をデコードして渡す。
type CategoryInput struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// CategoryService はカテゴリのCRUDと業務ルールを担う。
type CategoryService struct {
	q *db.Queries
}

func NewCategoryService(q *db.Queries) *CategoryService {
	return &CategoryService{q: q}
}

// List は種別で絞り込んだカテゴリ一覧を返す。typ が income/expense 以外なら全件。
func (s *CategoryService) List(ctx context.Context, typ string) ([]db.Category, error) {
	var params db.ListCategoriesByTypeParams
	if typ == "income" || typ == "expense" {
		params.Type = db.NullCategoriesType{CategoriesType: db.CategoriesType(typ), Valid: true}
	}
	return s.q.ListCategoriesByType(ctx, params)
}

// Create はカテゴリを作成し、作成後のカテゴリを返す。
func (s *CategoryService) Create(ctx context.Context, in CategoryInput) (db.Category, error) {
	name, typ, err := s.validate(in)
	if err != nil {
		return db.Category{}, err
	}

	cnt, err := s.q.CountCategoryByNameType(ctx, db.CountCategoryByNameTypeParams{Name: name, Type: typ})
	if err != nil {
		return db.Category{}, err
	}
	if cnt > 0 {
		return db.Category{}, &ConflictError{Message: "同じ種別に同名のカテゴリが既に存在します"}
	}

	id, err := s.q.CreateCategory(ctx, db.CreateCategoryParams{Name: name, Type: typ})
	if err != nil {
		return db.Category{}, err
	}
	return s.q.GetCategory(ctx, id)
}

// Update は既存カテゴリを更新し、更新後のカテゴリを返す。
func (s *CategoryService) Update(ctx context.Context, id int64, in CategoryInput) (db.Category, error) {
	name, typ, err := s.validate(in)
	if err != nil {
		return db.Category{}, err
	}

	if _, err := s.q.GetCategory(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Category{}, ErrNotFound
		}
		return db.Category{}, err
	}

	cnt, err := s.q.CountCategoryByNameTypeExcludingID(ctx, db.CountCategoryByNameTypeExcludingIDParams{Name: name, Type: typ, ID: id})
	if err != nil {
		return db.Category{}, err
	}
	if cnt > 0 {
		return db.Category{}, &ConflictError{Message: "同じ種別に同名のカテゴリが既に存在します"}
	}

	rows, err := s.q.UpdateCategory(ctx, db.UpdateCategoryParams{Name: name, Type: typ, ID: id})
	if err != nil {
		return db.Category{}, err
	}
	if rows == 0 {
		return db.Category{}, ErrNotFound
	}
	return s.q.GetCategory(ctx, id)
}

// Delete はカテゴリを削除する。取引・予算で使用中なら ConflictError、存在しなければ ErrNotFound。
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if _, err := s.q.GetCategory(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	usage, err := s.q.CountCategoryUsage(ctx, db.CountCategoryUsageParams{
		CategoryID:   id,
		CategoryID_2: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		return err
	}
	if usage > 0 {
		return &ConflictError{Message: "使用中のため削除できません"}
	}

	rows, err := s.q.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// validate はカテゴリ入力を検証し、DB 用の値へ変換する。
func (s *CategoryService) validate(in CategoryInput) (string, db.CategoriesType, error) {
	var fields []FieldError

	if in.Name == "" {
		fields = append(fields, FieldError{Field: "name", Message: "名称を入力してください"})
	} else if len([]rune(in.Name)) > 50 {
		fields = append(fields, FieldError{Field: "name", Message: "名称は50文字以内で入力してください"})
	}

	if in.Type != "income" && in.Type != "expense" {
		fields = append(fields, FieldError{Field: "type", Message: "種別は income または expense を指定してください"})
	}

	if len(fields) > 0 {
		return "", "", &ValidationError{Fields: fields}
	}
	return in.Name, db.CategoriesType(in.Type), nil
}
