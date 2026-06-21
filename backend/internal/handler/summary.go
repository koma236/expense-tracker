package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/koma236/expense-tracker/backend/internal/service"
)

// SummaryHandler は集計（/api/summary）のHTTPハンドラ。
type SummaryHandler struct {
	svc *service.SummaryService
}

func NewSummaryHandler(svc *service.SummaryService) *SummaryHandler {
	return &SummaryHandler{svc: svc}
}

// totalsResponse は api-spec §5.1 の totals。
type totalsResponse struct {
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
	Balance int64 `json:"balance"`
}

// expenseByCategoryItem は api-spec §5.1 の expense_by_category[] 要素。
type expenseByCategoryItem struct {
	Category categoryRef `json:"category"`
	Total    int64       `json:"total"`
}

// summaryResponse は GET /api/summary のレスポンス形。
// budget_progress は F-4 で実装するため、現状は常に空配列を返す。
type summaryResponse struct {
	YearMonth         string                  `json:"year_month"`
	Totals            totalsResponse          `json:"totals"`
	ExpenseByCategory []expenseByCategoryItem `json:"expense_by_category"`
	BudgetProgress    []any                   `json:"budget_progress"`
}

// Get は GET /api/summary。
func (h *SummaryHandler) Get(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.Summary(r.Context(), r.URL.Query().Get("year_month"))
	if err != nil {
		h.handleError(w, err)
		return
	}

	items := make([]expenseByCategoryItem, 0, len(res.ExpenseByCategory))
	for _, c := range res.ExpenseByCategory {
		items = append(items, expenseByCategoryItem{
			Category: categoryRef{ID: c.CategoryID, Name: c.CategoryName},
			Total:    c.Total,
		})
	}

	writeJSON(w, http.StatusOK, summaryResponse{
		YearMonth: res.YearMonth,
		Totals: totalsResponse{
			Income:  res.Totals.Income,
			Expense: res.Totals.Expense,
			Balance: res.Totals.Balance,
		},
		ExpenseByCategory: items,
		BudgetProgress:    []any{},
	})
}

// trendPointResponse は月別収支推移の1点。
type trendPointResponse struct {
	YearMonth string `json:"year_month"`
	Income    int64  `json:"income"`
	Expense   int64  `json:"expense"`
	Balance   int64  `json:"balance"`
}

// Trend は GET /api/summary/trend。
func (h *SummaryHandler) Trend(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	months := 0
	if m := q.Get("months"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			months = v
		}
	}

	points, err := h.svc.Trend(r.Context(), q.Get("year_month"), months)
	if err != nil {
		h.handleError(w, err)
		return
	}

	list := make([]trendPointResponse, 0, len(points))
	for _, p := range points {
		list = append(list, trendPointResponse{
			YearMonth: p.YearMonth,
			Income:    p.Income,
			Expense:   p.Expense,
			Balance:   p.Balance,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"months": list})
}

// handleError は service のエラーを HTTP レスポンスに変換する。
func (h *SummaryHandler) handleError(w http.ResponseWriter, err error) {
	var ve *service.ValidationError
	switch {
	case errors.As(err, &ve):
		details := make([]fieldError, 0, len(ve.Fields))
		for _, f := range ve.Fields {
			details = append(details, fieldError{Field: f.Field, Message: f.Message})
		}
		writeError(w, http.StatusBadRequest, codeValidation, ve.Error(), details)
	default:
		writeError(w, http.StatusInternalServerError, codeInternal, "サーバー内部エラーが発生しました", nil)
	}
}
