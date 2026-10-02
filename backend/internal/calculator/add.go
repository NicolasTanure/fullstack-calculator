package calculator

func Add(left, right float64) (float64, error) {
	if !isFinite(left) || !isFinite(right) {
		return 0, ErrInvalidOperand
	}

	result := left + right
	if !isFinite(result) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}
