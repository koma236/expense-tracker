package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
	"github.com/koma236/expense-tracker/backend/internal/service"
)

// CategoryHandler はカテゴリ（/api/categories）のHTTPハンドラ。
type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// categoryResponse は api-spec §3.1 のカテゴリレスポンス形。
type categoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func newCategoryResponse(c db.Category) categoryResponse {
	return categoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Type:      string(c.Type),
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

// List は GET /api/categories?type=income|expense。
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.List(r.Context(), r.URL.Query().Get("type"))
	if err != nil {
		h.handleError(w, err)
		return
	}

	list := make([]categoryResponse, 0, len(rows))
	for _, c := range rows {
		list = append(list, newCategoryResponse(c))
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": list})
}

// Create は POST /api/categories。
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	in, ok := decodeCategory(w, r)
	if !ok {
		return
	}
	cat, err := h.svc.Create(r.Context(), in)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, newCategoryResponse(cat))
}

// Update は PUT /api/categories/{id}。
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	in, ok := decodeCategory(w, r)
	if !ok {
		return
	}
	cat, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newCategoryResponse(cat))
}

// Delete は DELETE /api/categories/{id}。
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

// decodeCategory はリクエストボディをカテゴリ入力にデコードする。
func decodeCategory(w http.ResponseWriter, r *http.Request) (service.CategoryInput, bool) {
	var in service.CategoryInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, codeValidation, "リクエストボディの形式が不正です", nil)
		return service.CategoryInput{}, false
	}
	return in, true
}

// handleError は service のエラーを HTTP レスポンスに変換する。
func (h *CategoryHandler) handleError(w http.ResponseWriter, err error) {
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
