package httpapi_test

import (
	"net/http"
	"testing"
)

func TestMultiplySuccess(t *testing.T) {
	tests := []struct {
		name string
		body string
		want float64
	}{
		{name: "integers", body: `{"left":8,"right":2}`, want: 16},
		{name: "negative operand", body: `{"left":-8,"right":2}`, want: -16},
		{name: "zero result", body: `{"left":8,"right":0}`, want: 0},
		{name: "unrounded decimal", body: `{"left":0.1,"right":0.2}`, want: 0.020000000000000004},
		{name: "scientific notation", body: `{"left":1e3,"right":2.5e2}`, want: 250000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := postCalculation("/api/multiply", tt.body, "application/json")
			assertResultResponse(t, response, tt.want)
		})
	}
}

func TestMultiplyOverflowResponse(t *testing.T) {
	for _, body := range []string{
		`{"left":1e308,"right":2}`,
		`{"left":-1e308,"right":2}`,
	} {
		response := postCalculation("/api/multiply", body, "application/json")
		assertErrorResponse(t, response, http.StatusBadRequest, "non_finite_result")
	}
}
