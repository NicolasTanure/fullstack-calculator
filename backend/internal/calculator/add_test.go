package calculator_test

import (
	"errors"
	"math"
	"testing"

	"fullstack-calculator/backend/internal/calculator"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		left, right float64
		want        float64
	}{
		{name: "positive integers", left: 8, right: 2, want: 10},
		{name: "negative integers", left: -8, right: -2, want: -10},
		{name: "mixed signs", left: -8, right: 2, want: -6},
		{name: "decimals", left: 1.25, right: 2.5, want: 3.75},
		{name: "zero operands", left: 0, right: 0, want: 0},
		{name: "zero left operand", left: 0, right: 8, want: 8},
		{name: "zero right operand", left: 8, right: 0, want: 8},
		{name: "maximum finite operand", left: math.MaxFloat64, right: 0, want: math.MaxFloat64},
		{name: "large cancellation", left: math.MaxFloat64, right: -math.MaxFloat64, want: 0},
		{name: "smallest operands", left: math.SmallestNonzeroFloat64, right: math.SmallestNonzeroFloat64, want: 2 * math.SmallestNonzeroFloat64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Add(tt.left, tt.right)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddFloatingPointDecimals(t *testing.T) {
	got, err := calculator.Add(0.1, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-0.3) > 1e-15 {
		t.Errorf("result = %.17g, want approximately 0.3", got)
	}
}

func TestAddInvalidOperands(t *testing.T) {
	tests := []struct {
		name        string
		left, right float64
	}{
		{name: "NaN on left", left: math.NaN(), right: 1},
		{name: "NaN on right", left: 1, right: math.NaN()},
		{name: "infinity on left", left: math.Inf(1), right: 1},
		{name: "infinity on right", left: 1, right: math.Inf(-1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calculator.Add(tt.left, tt.right)
			if !errors.Is(err, calculator.ErrInvalidOperand) {
				t.Errorf("error = %v, want ErrInvalidOperand", err)
			}
		})
	}
}

func TestAddOverflow(t *testing.T) {
	for _, operand := range []float64{math.MaxFloat64, -math.MaxFloat64} {
		_, err := calculator.Add(operand, operand)
		if !errors.Is(err, calculator.ErrNonFiniteResult) {
			t.Errorf("Add(%v, %v) error = %v, want ErrNonFiniteResult", operand, operand, err)
		}
	}
}
