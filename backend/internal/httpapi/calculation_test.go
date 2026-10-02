package httpapi_test

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"fullstack-calculator/backend/internal/httpapi"
)

func TestAddSuccess(t *testing.T) {
	validBody := `{"left":8,"right":2}`
	tests := []struct {
		name        string
		body        string
		contentType string
		want        float64
	}{
		{name: "integers", body: validBody, want: 10},
		{name: "mixed signs", body: `{"left":-8,"right":2}`, want: -6},
		{name: "zero left operand", body: `{"left":0,"right":8}`, want: 8},
		{name: "zero right operand", body: `{"left":8,"right":0}`, want: 8},
		{name: "scientific notation", body: `{"left":1e3,"right":2.5e2}`, want: 1250},
		{name: "unrounded result", body: `{"left":0.1,"right":0.2}`, want: 0.30000000000000004},
		{name: "reversed field order", body: `{"right":2,"left":8}`, want: 10},
		{name: "surrounding whitespace", body: "\n " + validBody + " \t\n", want: 10},
		{name: "content type parameters", body: validBody, contentType: "application/json; charset=utf-8", want: 10},
		{name: "exact body limit", body: validBody + strings.Repeat(" ", 1024-len(validBody)), want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentType := tt.contentType
			if contentType == "" {
				contentType = "application/json"
			}
			response := postCalculation("/api/add", tt.body, contentType)
			assertResultResponse(t, response, tt.want)
		})
	}
}

func TestCalculationInvalidInput(t *testing.T) {
	for _, path := range []string{"/api/add", "/api/subtract", "/api/multiply", "/api/divide"} {
		t.Run(path, func(t *testing.T) {
			tests := []struct {
				name string
				body string
			}{
				{name: "empty body", body: ""},
				{name: "whitespace only", body: " \n\t"},
				{name: "malformed JSON", body: `{"left":1,`},
				{name: "empty object", body: `{}`},
				{name: "missing left", body: `{"right":2}`},
				{name: "missing right", body: `{"left":1}`},
				{name: "null left", body: `{"left":null,"right":2}`},
				{name: "null right", body: `{"left":1,"right":null}`},
				{name: "numeric string", body: `{"left":"1","right":2}`},
				{name: "boolean operand", body: `{"left":1,"right":true}`},
				{name: "array operand", body: `{"left":[],"right":2}`},
				{name: "object operand", body: `{"left":1,"right":{}}`},
				{name: "unknown field", body: `{"left":1,"right":2,"extra":3}`},
				{name: "incorrect field casing", body: `{"Left":1,"right":2}`},
				{name: "root null", body: `null`},
				{name: "root array", body: `[1,2]`},
				{name: "root number", body: `1`},
				{name: "root string", body: `"value"`},
				{name: "root boolean", body: `true`},
				{name: "trailing JSON object", body: `{"left":1,"right":2}{}`},
				{name: "trailing JSON value", body: `{"left":1,"right":2}true`},
				{name: "trailing garbage", body: `{"left":1,"right":2}garbage`},
				{name: "operand out of range", body: `{"left":1e309,"right":2}`},
				{name: "NaN token", body: `{"left":NaN,"right":2}`},
				{name: "infinity token", body: `{"left":1,"right":Infinity}`},
				{name: "duplicate left", body: `{"left":1,"left":8,"right":2}`},
				{name: "duplicate right", body: `{"left":8,"right":0,"right":2}`},
				{name: "duplicate identical values", body: `{"left":8,"right":2,"left":8}`},
				{name: "null overwritten by a number", body: `{"left":null,"left":8,"right":2}`},
				{name: "escaped duplicate field", body: `{"left":8,"\u006ceft":1,"right":2}`},
				{name: "leading plus", body: `{"left":+1,"right":2}`},
				{name: "leading zero", body: `{"left":01,"right":2}`},
				{name: "missing integer part", body: `{"left":.5,"right":2}`},
				{name: "missing fractional part", body: `{"left":1.,"right":2}`},
				{name: "incomplete exponent", body: `{"left":1e,"right":2}`},
				{name: "hexadecimal", body: `{"left":0x10,"right":2}`},
				{name: "comma decimal", body: `{"left":1,5,"right":2}`},
				{name: "trailing comma", body: `{"left":1,"right":2,}`},
				{name: "JSON comment", body: `{"left":1,"right":2} /* comment */`},
				{name: "non-JSON trailing whitespace", body: "{\"left\":1,\"right\":2}\u00a0"},
				{name: "negative operand out of range", body: `{"left":1,"right":-1e309}`},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					response := postCalculation(path, tt.body, "application/json")
					assertErrorResponse(t, response, http.StatusBadRequest, "invalid_input")
				})
			}
		})
	}
}

func TestAddOverflowResponse(t *testing.T) {
	for _, body := range []string{
		`{"left":1e308,"right":1e308}`,
		`{"left":-1e308,"right":-1e308}`,
	} {
		response := postCalculation("/api/add", body, "application/json")
		assertErrorResponse(t, response, http.StatusBadRequest, "non_finite_result")
	}
}

func TestCalculationNumericBoundaries(t *testing.T) {
	tests := []struct {
		name, path, body string
		want             float64
		errorCode        string
	}{
		{name: "maximum finite JSON operand", path: "/api/add", body: `{"left":1.7976931348623157e308,"right":0}`, want: math.MaxFloat64},
		{name: "minimum nonzero JSON operand", path: "/api/add", body: `{"left":5e-324,"right":0}`, want: math.SmallestNonzeroFloat64},
		{name: "maximum finite cancellation", path: "/api/subtract", body: `{"left":1.7976931348623157e308,"right":1.7976931348623157e308}`, want: 0},
		{name: "zero times a large finite value", path: "/api/multiply", body: `{"left":0,"right":1.7976931348623157e308}`, want: 0},
		{name: "multiplication underflow", path: "/api/multiply", body: `{"left":5e-324,"right":0.5}`, want: 0},
		{name: "division underflow", path: "/api/divide", body: `{"left":5e-324,"right":2}`, want: 0},
		{name: "negative zero numerator", path: "/api/divide", body: `{"left":-0,"right":2}`, want: math.Copysign(0, -1)},
		{name: "divisor underflows to positive zero", path: "/api/divide", body: `{"left":1,"right":1e-999}`, errorCode: "division_by_zero"},
		{name: "divisor underflows to negative zero", path: "/api/divide", body: `{"left":1,"right":-1e-999}`, errorCode: "division_by_zero"},
		{name: "tiny divisor overflows result", path: "/api/divide", body: `{"left":1,"right":5e-324}`, errorCode: "non_finite_result"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := postCalculation(tt.path, tt.body, "application/json")
			if tt.errorCode != "" {
				assertErrorResponse(t, response, http.StatusBadRequest, tt.errorCode)
			} else {
				assertResultResponse(t, response, tt.want)
			}
		})
	}
}

func TestCalculationBodyTooLarge(t *testing.T) {
	for _, path := range []string{"/api/add", "/api/subtract", "/api/multiply", "/api/divide"} {
		t.Run(path, func(t *testing.T) {
			validBody := `{"left":1,"right":2}`
			tests := []struct {
				name string
				body string
			}{
				{name: "valid JSON with excess whitespace", body: validBody + strings.Repeat(" ", 1025-len(validBody))},
				{name: "oversized malformed body", body: strings.Repeat("x", 1025)},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					response := postCalculation(path, tt.body, "application/json")
					assertErrorResponse(t, response, http.StatusRequestEntityTooLarge, "request_too_large")
				})
			}
		})
	}
}

func TestCalculationContentType(t *testing.T) {
	for _, path := range []string{"/api/add", "/api/subtract", "/api/multiply", "/api/divide"} {
		t.Run(path, func(t *testing.T) {
			tests := []struct {
				name        string
				contentType string
			}{
				{name: "missing content type", contentType: ""},
				{name: "incompatible media type", contentType: "text/plain"},
				{name: "malformed parameter", contentType: "application/json; charset"},
				{name: "multiple media types", contentType: "application/json, text/plain"},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					response := postCalculation(path, `{"left":1,"right":2}`, tt.contentType)
					assertErrorResponse(t, response, http.StatusUnsupportedMediaType, "unsupported_media_type")
				})
			}
		})
	}
}

func TestCalculationUnsupportedMethod(t *testing.T) {
	for _, path := range []string{"/api/add", "/api/subtract", "/api/multiply", "/api/divide"} {
		t.Run(path, func(t *testing.T) {
			for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions} {
				t.Run(method, func(t *testing.T) {
					request := httptest.NewRequest(method, path, nil)
					response := httptest.NewRecorder()
					httpapi.NewHandler().ServeHTTP(response, request)
					assertErrorResponse(t, response, http.StatusMethodNotAllowed, "method_not_allowed")
					if got := response.Header().Get("Allow"); got != http.MethodPost {
						t.Errorf("Allow = %q, want POST", got)
					}
				})
			}
		})
	}
}

func TestCalculationAmbiguousContentType(t *testing.T) {
	for _, path := range []string{"/api/add", "/api/subtract", "/api/multiply", "/api/divide"} {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"left":8,"right":2}`))
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Content-Type", "text/plain")
			response := httptest.NewRecorder()
			httpapi.NewHandler().ServeHTTP(response, request)
			assertErrorResponse(t, response, http.StatusUnsupportedMediaType, "unsupported_media_type")
		})
	}
}

func TestCalculationBodyReadFailure(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/add", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Body = io.NopCloser(iotest.ErrReader(errors.New("transport read failed")))
	response := httptest.NewRecorder()
	httpapi.NewHandler().ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusBadRequest, "invalid_input")
	var body map[string]map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if got := body["error"]["message"]; got != "Could not read request body." {
		t.Errorf("message = %q, want Could not read request body.", got)
	}
}

func postCalculation(path, body, contentType string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	response := httptest.NewRecorder()
	httpapi.NewHandler().ServeHTTP(response, request)
	return response
}

func assertResultResponse(t *testing.T, response *httptest.ResponseRecorder, result float64) {
	t.Helper()
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	want := map[string]any{"result": result}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, status, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	errorBody, ok := body["error"].(map[string]any)
	if !ok || len(body) != 1 || len(errorBody) != 2 {
		t.Fatalf("unexpected error envelope: %v", body)
	}
	if got := errorBody["code"]; got != code {
		t.Errorf("error code = %v, want %q", got, code)
	}
	if message, ok := errorBody["message"].(string); !ok || message == "" {
		t.Errorf("error message must be a non-empty string, got %v", errorBody["message"])
	}
}
