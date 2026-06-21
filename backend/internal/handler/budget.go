package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
	"github.com/koma236/expense-tracker/backend/internal/service"
)

// BudgetHandler は予算（/api/budgets）のHTTPハンドラ。
type BudgetHandler struct {
	svc *service.BudgetService
}

func NewBudgetHandler(svc *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

// budgetResponse は api-spec §4 の予算レスポンス形。
// category は月全体予算の場合 null。
type budgetResponse struct {
	ID        int64        `json:"id"`
	Category  *categoryRef `json:"category"`
	YearMonth string       `json:"year_month"`
	Amount    int32        `json:"amount"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

// List は GET /api/budgets?year_month=YYYY-MM。
func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.List(r.Context(), r.URL.Query().Get("year_month"))
	if err != nil {
		h.handleError(w, err)
		return
	}

	list := make([]budgetResponse, 0, len(rows))
	for _, b := range rows {
		list = append(list, budgetResponse{
			ID:        b.ID,
			Category:  categoryRefFromNull(b.CategoryID.Int64, b.CategoryID.Valid, b.CategoryName.String),
			YearMonth: b.YearMonth,
			Amount:    b.Amount,
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
			UpdatedAt: b.UpdatedAt.Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"budgets": list})
}

// Create は POST /api/budgets。
func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	in, ok := decodeBudget(w, r)
	if !ok {
		return
	}
	b, err := h.svc.Create(r.Context(), in)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, budgetRowToResponse(b))
}

// Update は PUT /api/budgets/{id}。金額のみ更新する。
func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	in, ok := decodeBudget(w, r)
	if !ok {
		return
	}
	b, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, budgetRowToResponse(b))
}

// Delete は DELETE /api/budgets/{id}。
func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func budgetRowToResponse(b db.GetBudgetRow) budgetResponse {
	return budgetResponse{
		ID:        b.ID,
		Category:  categoryRefFromNull(b.CategoryID.Int64, b.CategoryID.Valid, b.CategoryName.String),
		YearMonth: b.YearMonth,
		Amount:    b.Amount,
		CreatedAt: b.CreatedAt.Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt.Format(time.RFC3339),
	}
}

// categoryRefFromNull は NULL 可なカテゴリ参照を組み立てる。月全体予算なら nil。
func categoryRefFromNull(id int64, valid bool, name string) *categoryRef {
	if !valid {
		return nil
	}
	return &categoryRef{ID: id, Name: name}
}

// decodeBudget はリクエストボディを予算入力にデコードする。
func decodeBudget(w http.ResponseWriter, r *http.Request) (service.BudgetInput, bool) {
	var in service.BudgetInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, codeValidation, "リクエストボディの形式が不正です", nil)
		return service.BudgetInput{}, false
	}
	return in, true
}

// handleError は service のエラーを HTTP レスポンスに変換する。
func (h *BudgetHandler) handleError(w http.ResponseWriter, err error) {
	var ve *service.ValidationError
	var ce *service.ConflictError
	switch {
	case errors.As(err, &ve):
		details := make([]fieldError, 0, len(ve.Fields))
		for _, f := range ve.Fields {
			details = append(details, fieldError{Field: f.Field, Message: f.Message})
		}
		writeError(w, http.StatusBadRequest, codeValidation, ve.Error(), details)
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, codeNotFound, "リソースが見つかりません", nil)
	case errors.As(err, &ce):
		writeError(w, http.StatusConflict, codeConflict, ce.Error(), nil)
	default:
		writeError(w, http.StatusInternalServerError, codeInternal, "サーバー内部エラーが発生しました", nil)
	}
}
