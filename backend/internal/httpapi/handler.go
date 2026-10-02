package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"fullstack-calculator/backend/internal/calculator"
)

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/add":
			calculate(w, r, calculator.Add)
		case "/api/subtract":
			calculate(w, r, calculator.Subtract)
		case "/api/multiply":
			calculate(w, r, calculator.Multiply)
		case "/api/divide":
			calculate(w, r, calculator.Divide)
		default:
			notFound(w, r)
		}
	})
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusNotFound, "not_found", "Endpoint not found.")
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{
		Error: apiError{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Could not write JSON response: %v", err)
	}
}
