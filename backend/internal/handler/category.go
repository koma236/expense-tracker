package handler

import (
	"net/http"
	"time"

	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
)

// CategoryHandler はカテゴリ（/api/categories）のHTTPハンドラ。
// プロトタイプでは一覧取得（読み取り）のみを提供する。
type CategoryHandler struct {
	q *db.Queries
}

func NewCategoryHandler(q *db.Queries) *CategoryHandler {
	return &CategoryHandler{q: q}
}

// categoryResponse は api-spec §3.1 のカテゴリレスポンス形。
type categoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// List は GET /api/categories?type=income|expense。
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	typ := r.URL.Query().Get("type")

	var params db.ListCategoriesByTypeParams
	if typ == "income" || typ == "expense" {
		params.Type = db.NullCategoriesType{CategoriesType: db.CategoriesType(typ), Valid: true}
	}

	rows, err := h.q.ListCategoriesByType(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "サーバー内部エラーが発生しました", nil)
		return
	}

	list := make([]categoryResponse, 0, len(rows))
	for _, c := range rows {
		list = append(list, categoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			Type:      string(c.Type),
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
			UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": list})
}
