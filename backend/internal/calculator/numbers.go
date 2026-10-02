package calculator

import (
	"errors"
	"math"
)

var (
	ErrInvalidOperand  = errors.New("operands must be finite numbers")
	ErrNonFiniteResult = errors.New("the calculation result must be finite")
)

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
