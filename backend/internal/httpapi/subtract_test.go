package httpapi_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestSubtractSuccess(t *testing.T) {
	validBody := `{"left":8,"right":2}`
	tests := []struct {
		name string
		body string
		want float64
	}{
		{name: "integers", body: validBody, want: 6},
		{name: "operand order", body: `{"left":2,"right":8}`, want: -6},
		{name: "negative right", body: `{"left":8,"right":-2}`, want: 10},
		{name: "zero result", body: `{"left":8,"right":8}`, want: 0},
		{name: "unrounded decimal", body: `{"left":0.3,"right":0.2}`, want: 0.09999999999999998},
		{name: "reversed field order", body: `{"right":2,"left":8}`, want: 6},
		{name: "exact body limit", body: validBody + strings.Repeat(" ", 1024-len(validBody)), want: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := postCalculation("/api/subtract", tt.body, "application/json; charset=utf-8")
			assertResultResponse(t, response, tt.want)
		})
	}
}

func TestSubtractOverflowResponse(t *testing.T) {
	for _, body := range []string{
		`{"left":1e308,"right":-1e308}`,
		`{"left":-1e308,"right":1e308}`,
	} {
		response := postCalculation("/api/subtract", body, "application/json")
		assertErrorResponse(t, response, http.StatusBadRequest, "non_finite_result")
	}
}
