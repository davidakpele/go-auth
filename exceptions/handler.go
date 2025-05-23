package exceptions

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendSuccessResponse sends a success JSON response
func SendSuccessResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(JSONResponse{
		Status:  "success",
		Message: message,
	})
}

// SendErrorResponse sends a standardized JSON error response
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(JSONResponse{
		Status:  "error",
		Message: message,
	})
}