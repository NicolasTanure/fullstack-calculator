package calculator_test

import (
	"errors"
	"math"
	"testing"

	"fullstack-calculator/backend/internal/calculator"
)

func TestMultiply(t *testing.T) {
	tests := []struct {
		name        string
		left, right float64
		want        float64
	}{
		{name: "positive integers", left: 8, right: 2, want: 16},
		{name: "negative left", left: -8, right: 2, want: -16},
		{name: "negative right", left: 8, right: -2, want: -16},
		{name: "both negative", left: -8, right: -2, want: 16},
		{name: "decimals", left: 1.25, right: 2.5, want: 3.125},
		{name: "zero left", left: 0, right: 8, want: 0},
		{name: "zero right", left: 8, right: 0, want: 0},
		{name: "maximum finite result", left: math.MaxFloat64, right: 1, want: math.MaxFloat64},
		{name: "smallest finite result", left: math.SmallestNonzeroFloat64, right: 1, want: math.SmallestNonzeroFloat64},
		{name: "underflow to zero", left: math.SmallestNonzeroFloat64, right: 0.5, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Multiply(tt.left, tt.right)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiplyFloatingPointDecimals(t *testing.T) {
	got, err := calculator.Multiply(0.1, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-0.02) > 1e-17 {
		t.Errorf("result = %.17g, want approximately 0.02", got)
	}
}

func TestMultiplyInvalidOperands(t *testing.T) {
	for _, operands := range [][2]float64{
		{math.NaN(), 1}, {1, math.NaN()},
		{math.Inf(1), 0}, {0, math.Inf(-1)},
	} {
		_, err := calculator.Multiply(operands[0], operands[1])
		if !errors.Is(err, calculator.ErrInvalidOperand) {
			t.Errorf("Multiply(%v, %v) error = %v, want ErrInvalidOperand", operands[0], operands[1], err)
		}
	}
}

func TestMultiplyOverflow(t *testing.T) {
	for _, left := range []float64{math.MaxFloat64, -math.MaxFloat64} {
		_, err := calculator.Multiply(left, 2)
		if !errors.Is(err, calculator.ErrNonFiniteResult) {
			t.Errorf("Multiply(%v, 2) error = %v, want ErrNonFiniteResult", left, err)
		}
	}
}
