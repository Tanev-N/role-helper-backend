package games

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	log.Printf("Games API error [%d]: %s - %v", statusCode, message, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error:   message,
		Message: err.Error(),
	})
}

func writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		_, _ = w.Write([]byte("{}"))
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}
