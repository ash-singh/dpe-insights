package helpers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse error response.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// SuccessResponse success response.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// GetError return error response.
func GetError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	response := ErrorResponse{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: err.Error(),
	}

	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(message)
}
