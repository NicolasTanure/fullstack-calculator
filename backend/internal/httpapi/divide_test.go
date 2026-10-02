package httpapi_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestDivideSuccess(t *testing.T) {
	tests := []struct {
		name string
		body string
		want float64
	}{
		{name: "integers", body: `{"left":8,"right":2}`, want: 4},
		{name: "operand order", body: `{"left":2,"right":8}`, want: 0.25},
		{name: "zero numerator", body: `{"left":0,"right":2}`, want: 0},
		{name: "unrounded decimal", body: `{"left":1,"right":3}`, want: 0.3333333333333333},
		{name: "smallest finite divisor", body: `{"left":5e-324,"right":5e-324}`, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := postCalculation("/api/divide", tt.body, "application/json")
			assertResultResponse(t, response, tt.want)
		})
	}
}

func TestDivideByZeroResponse(t *testing.T) {
	for _, body := range []string{
		`{"left":8,"right":0}`,
		`{"left":8,"right":-0}`,
		`{"left":0,"right":0}`,
	} {
		response := postCalculation("/api/divide", body, "application/json")
		assertErrorResponse(t, response, http.StatusBadRequest, "division_by_zero")
		var result map[string]map[string]string
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if got := result["error"]["message"]; got != "Cannot divide by zero." {
			t.Errorf("message = %q, want Cannot divide by zero.", got)
		}
	}
}

func TestDivideOverflowResponse(t *testing.T) {
	for _, body := range []string{
		`{"left":1e308,"right":0.5}`,
		`{"left":-1e308,"right":0.5}`,
	} {
		response := postCalculation("/api/divide", body, "application/json")
		assertErrorResponse(t, response, http.StatusBadRequest, "non_finite_result")
	}
}
