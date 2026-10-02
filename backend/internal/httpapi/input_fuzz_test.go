package httpapi_test

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"testing"
)

func FuzzCalculationInput(f *testing.F) {
	paths := []string{"/api/add", "/api/subtract", "/api/multiply", "/api/divide"}
	for operation := range paths {
		for _, body := range []string{
			`{"left":8,"right":2}`, `{"left":0,"right":-0}`,
			`{"left":1e308,"right":5e-324}`, `{"left":0.1,"right":0.2}`,
			`{"left":null,"left":8,"right":2}`, `{"left":1,"right":2}{}`,
			`{"left":[],"right":2}`, `{"left":1e309,"right":2}`, "", "{",
			strings.Repeat("x", 1025),
		} {
			f.Add(uint8(operation), []byte(body))
		}
	}
	f.Fuzz(func(t *testing.T, operation uint8, body []byte) {
		path := paths[int(operation)%len(paths)]
		response := postCalculation(path, string(body), "application/json")
		if len(body) > 1024 {
			assertErrorResponse(t, response, http.StatusRequestEntityTooLarge, "request_too_large")
			return
		}
		switch response.Code {
		case http.StatusOK:
			var operands map[string]*float64
			if err := json.Unmarshal(body, &operands); err != nil || len(operands) != 2 || operands["left"] == nil || operands["right"] == nil {
				t.Fatalf("accepted invalid operands: %q", body)
			}
			left, right := *operands["left"], *operands["right"]
			var want float64
			switch path {
			case "/api/add":
				want = left + right
			case "/api/subtract":
				want = left - right
			case "/api/multiply":
				want = left * right
			case "/api/divide":
				if right == 0 {
					t.Fatalf("accepted a zero divisor: %q", body)
				}
				want = left / right
			}
			if math.IsInf(want, 0) || math.IsNaN(want) {
				t.Fatalf("accepted a non-finite result: %q", body)
			}
			assertResultResponse(t, response, want)
		case http.StatusBadRequest:
			var envelope struct {
				Error struct{ Code string }
			}
			if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
				t.Fatalf("invalid JSON error: %v", err)
			}
			code := envelope.Error.Code
			if code != "invalid_input" && code != "division_by_zero" && code != "non_finite_result" {
				t.Fatalf("unexpected error code: %q", code)
			}
			assertErrorResponse(t, response, http.StatusBadRequest, code)
		default:
			t.Fatalf("unexpected status %d for %q", response.Code, body)
		}
	})
}
