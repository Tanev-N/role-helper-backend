package user

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   err.Error(),
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
	log.Printf("Ошибка: %s - %s", message, err.Error())
}

func writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SuccessResponse{
		Data: data,
	}

	json.NewEncoder(w).Encode(response)
}
