package calculator_test

import (
	"errors"
	"math"
	"testing"

	"fullstack-calculator/backend/internal/calculator"
)

func TestSubtract(t *testing.T) {
	tests := []struct {
		name        string
		left, right float64
		want        float64
	}{
		{name: "positive integers", left: 8, right: 2, want: 6},
		{name: "reversed operands", left: 2, right: 8, want: -6},
		{name: "negative integers", left: -8, right: -2, want: -6},
		{name: "subtract negative", left: 8, right: -2, want: 10},
		{name: "negative left", left: -8, right: 2, want: -10},
		{name: "decimals", left: 1.25, right: 2.5, want: -1.25},
		{name: "zero left", left: 0, right: 8, want: -8},
		{name: "zero right", left: 8, right: 0, want: 8},
		{name: "equal operands", left: 8, right: 8, want: 0},
		{name: "largest finite cancellation", left: math.MaxFloat64, right: math.MaxFloat64, want: 0},
		{name: "smallest finite difference", left: 0, right: math.SmallestNonzeroFloat64, want: -math.SmallestNonzeroFloat64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Subtract(tt.left, tt.right)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubtractFloatingPointDecimals(t *testing.T) {
	got, err := calculator.Subtract(0.3, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-0.1) > 1e-15 {
		t.Errorf("result = %.17g, want approximately 0.1", got)
	}
}

func TestSubtractInvalidOperands(t *testing.T) {
	for _, operands := range [][2]float64{
		{math.NaN(), 1}, {1, math.NaN()},
		{math.Inf(1), 1}, {1, math.Inf(-1)},
	} {
		_, err := calculator.Subtract(operands[0], operands[1])
		if !errors.Is(err, calculator.ErrInvalidOperand) {
			t.Errorf("Subtract(%v, %v) error = %v, want ErrInvalidOperand", operands[0], operands[1], err)
		}
	}
}

func TestSubtractOverflow(t *testing.T) {
	for _, left := range []float64{math.MaxFloat64, -math.MaxFloat64} {
		_, err := calculator.Subtract(left, -left)
		if !errors.Is(err, calculator.ErrNonFiniteResult) {
			t.Errorf("Subtract(%v, %v) error = %v, want ErrNonFiniteResult", left, -left, err)
		}
	}
}
