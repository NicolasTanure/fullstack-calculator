package calculator

import "errors"

var ErrDivisionByZero = errors.New("cannot divide by zero")

func Divide(left, right float64) (float64, error) {
	if !isFinite(left) || !isFinite(right) {
		return 0, ErrInvalidOperand
	}
	if right == 0 {
		return 0, ErrDivisionByZero
	}

	result := left / right
	if !isFinite(result) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}
