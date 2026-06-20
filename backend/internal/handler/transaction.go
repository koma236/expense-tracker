package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
	"github.com/koma236/expense-tracker/backend/internal/service"
)

// TransactionHandler は取引（/api/transactions）のHTTPハンドラ。
type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// categoryRef はレスポンス内のカテゴリ参照（api-spec）。
type categoryRef struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// transactionResponse は api-spec §2 の取引レスポンス形。
type transactionResponse struct {
	ID         int64       `json:"id"`
	OccurredOn string      `json:"occurred_on"`
	Amount     int32       `json:"amount"`
	Type       string      `json:"type"`
	Category   categoryRef `json:"category"`
	Memo       *string     `json:"memo"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
}

func newTransactionResponse(id int64, occurredOn time.Time, amount int32, typ, categoryName string, categoryID int64, memo sql.NullString, createdAt, updatedAt time.Time) transactionResponse {
	var memoPtr *string
	if memo.Valid {
		memoPtr = &memo.String
	}
	return transactionResponse{
		ID:         id,
		OccurredOn: occurredOn.Format("2006-01-02"),
		Amount:     amount,
		Type:       typ,
		Category:   categoryRef{ID: categoryID, Name: categoryName},
		Memo:       memoPtr,
		CreatedAt:  createdAt.Format(time.RFC3339),
		UpdatedAt:  updatedAt.Format(time.RFC3339),
	}
}

// List は GET /api/transactions。
func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := service.ListFilter{
		YearMonth: q.Get("year_month"),
		Type:      q.Get("type"),
	}
	if cid := q.Get("category_id"); cid != "" {
		if v, err := strconv.ParseInt(cid, 10, 64); err == nil {
			filter.CategoryID = v
		}
	}

	rows, err := h.svc.List(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	list := make([]transactionResponse, 0, len(rows))
	for _, row := range rows {
		list = append(list, newTransactionResponse(
			row.ID, row.OccurredOn, row.Amount, string(row.Type),
			row.CategoryName, row.CategoryID, row.Memo, row.CreatedAt, row.UpdatedAt,
		))
	}
	writeJSON(w, http.StatusOK, map[string]any{"transactions": list})
}

// Get は GET /api/transactions/{id}。
func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	row, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, getRowToResponse(row))
}

// Create は POST /api/transactions。
func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	in, ok := decodeTransaction(w, r)
	if !ok {
		return
	}
	row, err := h.svc.Create(r.Context(), in)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, getRowToResponse(row))
}

// Update は PUT /api/transactions/{id}。
func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	in, ok := decodeTransaction(w, r)
	if !ok {
		return
	}
	row, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, getRowToResponse(row))
}

// Delete は DELETE /api/transactions/{id}。
func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func getRowToResponse(row db.GetTransactionRow) transactionResponse {
	return newTransactionResponse(
		row.ID, row.OccurredOn, row.Amount, string(row.Type),
		row.CategoryName, row.CategoryID, row.Memo, row.CreatedAt, row.UpdatedAt,
	)
}

// decodeTransaction はリクエストボディを取引入力にデコードする。
func decodeTransaction(w http.ResponseWriter, r *http.Request) (service.TransactionInput, bool) {
	var in service.TransactionInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, codeValidation, "リクエストボディの形式が不正です", nil)
		return service.TransactionInput{}, false
	}
	return in, true
}

// parseID はパスパラメータ {id} を解釈する。
func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusNotFound, codeNotFound, "リソースが見つかりません", nil)
		return 0, false
	}
	return id, true
}

// handleError は service のエラーを HTTP レスポンスに変換する。
func (h *TransactionHandler) handleError(w http.ResponseWriter, err error) {
	var ve *service.ValidationError
	switch {
	case errors.As(err, &ve):
		details := make([]fieldError, 0, len(ve.Fields))
		for _, f := range ve.Fields {
			details = append(details, fieldError{Field: f.Field, Message: f.Message})
		}
		writeError(w, http.StatusBadRequest, codeValidation, ve.Error(), details)
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, codeNotFound, "リソースが見つかりません", nil)
	default:
		writeError(w, http.StatusInternalServerError, codeInternal, "サーバー内部エラーが発生しました", nil)
	}
}
