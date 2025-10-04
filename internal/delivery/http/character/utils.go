package character

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	log.Printf("API Error [%d]: %s - %v", statusCode, message, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   message,
		Message: err.Error(),
	}

	json.NewEncoder(w).Encode(response)
}

func writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func isValidationError(err error) bool {
	errLower := strings.ToLower(err.Error())
	return strings.Contains(errLower, "неверн") ||
		strings.Contains(errLower, "не может быть пустым") ||
		strings.Contains(errLower, "должен быть от") ||
		strings.Contains(errLower, "должны быть от") ||
		strings.Contains(errLower, "доступные расы") ||
		strings.Contains(errLower, "доступные классы")
}
