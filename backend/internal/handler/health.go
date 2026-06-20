// Package handler は HTTP リクエストを受け取り、JSON を返す層。
package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler は疎通確認用のハンドラ。
type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

// Health は GET /api/health。DB への Ping 結果を含めて返す。
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	res := healthResponse{Status: "ok", DB: "ok"}
	statusCode := http.StatusOK

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.db.PingContext(ctx); err != nil {
		res.DB = "ng"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(res)
}
