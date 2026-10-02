package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"fullstack-calculator/backend/internal/calculator"
)

const maxBodyBytes = 1024

type calculationResponse struct {
	Result float64 `json:"result"`
}

func calculate(w http.ResponseWriter, r *http.Request, operation func(float64, float64) (float64, error)) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST for this endpoint.")
		return
	}

	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != "application/json" || len(r.Header.Values("Content-Type")) != 1 {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json.")
		return
	}

	limitedBody := http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer limitedBody.Close()
	body, err := io.ReadAll(limitedBody)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "Request body must not exceed 1024 bytes.")
		} else {
			writeError(w, http.StatusBadRequest, "invalid_input", "Could not read request body.")
		}
		return
	}

	left, right, err := decodeOperands(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_input", err.Error())
		return
	}

	result, err := operation(left, right)
	if err != nil {
		switch {
		case errors.Is(err, calculator.ErrDivisionByZero):
			writeError(w, http.StatusBadRequest, "division_by_zero", "Cannot divide by zero.")
		case errors.Is(err, calculator.ErrNonFiniteResult):
			writeError(w, http.StatusBadRequest, "non_finite_result", "The calculation result must be finite.")
		default:
			writeError(w, http.StatusBadRequest, "invalid_input", "Operands must be finite numbers.")
		}
		return
	}
	writeJSON(w, http.StatusOK, calculationResponse{Result: result})
}

func decodeOperands(body []byte) (float64, float64, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	invalidJSON := errors.New("Request body must be a JSON object with numeric left and right fields.")
	opening, err := decoder.Token()
	if err != nil || opening != json.Delim('{') {
		return 0, 0, invalidJSON
	}
	operands := make(map[string]float64, 2)
	for decoder.More() {
		token, err := decoder.Token()
		field, ok := token.(string)
		if err != nil || !ok {
			return 0, 0, invalidJSON
		}
		if field != "left" && field != "right" {
			return 0, 0, errors.New("Provide exactly two numeric fields: left and right.")
		}
		if _, exists := operands[field]; exists {
			return 0, 0, errors.New("Duplicate operand fields are not allowed.")
		}
		var value *float64
		if err := decoder.Decode(&value); err != nil || value == nil {
			return 0, 0, invalidJSON
		}
		operands[field] = *value
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim('}') {
		return 0, 0, invalidJSON
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return 0, 0, errors.New("Request body must contain a single JSON object.")
	}
	if len(operands) != 2 {
		return 0, 0, errors.New("Provide exactly two numeric fields: left and right.")
	}
	return operands["left"], operands["right"], nil
}
