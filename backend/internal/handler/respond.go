package handler

import (
	"encoding/json"
	"net/http"
)

// エラーコード（api-spec §1.2 の `code`）。
const (
	codeValidation = "VALIDATION_ERROR"
	codeNotFound   = "NOT_FOUND"
	codeConflict   = "CONFLICT"
	codeInternal   = "INTERNAL_ERROR"
)

// fieldError はフィールド単位のバリデーションエラー。
type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// errorBody は api-spec §1.2 のエラーレスポンス形。
type errorBody struct {
	Error struct {
		Code    string       `json:"code"`
		Message string       `json:"message"`
		Details []fieldError `json:"details,omitempty"`
	} `json:"error"`
}

// writeJSON は status と任意の値を JSON で返す。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// writeError は共通エラー形式で返す。details は任意（バリデーション時のみ）。
func writeError(w http.ResponseWriter, status int, code, message string, details []fieldError) {
	var body errorBody
	body.Error.Code = code
	body.Error.Message = message
	body.Error.Details = details
	writeJSON(w, status, body)
}
