package response

import (
	"encoding/json"
	"net/http"
)

type successBody struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

type errorBody struct {
	Success bool `json:"success"`
	Error   any  `json:"error"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func OK(w http.ResponseWriter, data any, message string) {
	write(w, http.StatusOK, successBody{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(w http.ResponseWriter, data any, message string) {
	write(w, http.StatusCreated, successBody{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Error(w http.ResponseWriter, status int, code, message string, details any) {
	write(w, status, errorBody{
		Success: false,
		Error: errorPayload{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
