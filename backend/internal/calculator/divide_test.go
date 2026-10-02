package calculator_test

import (
	"errors"
	"math"
	"testing"

	"fullstack-calculator/backend/internal/calculator"
)

func TestDivide(t *testing.T) {
	tests := []struct {
		name        string
		left, right float64
		want        float64
	}{
		{name: "positive integers", left: 8, right: 2, want: 4},
		{name: "reversed operands", left: 2, right: 8, want: 0.25},
		{name: "negative left", left: -8, right: 2, want: -4},
		{name: "negative right", left: 8, right: -2, want: -4},
		{name: "both negative", left: -8, right: -2, want: 4},
		{name: "decimals", left: 1.25, right: 2.5, want: 0.5},
		{name: "zero numerator", left: 0, right: 2, want: 0},
		{name: "maximum finite result", left: math.MaxFloat64, right: 1, want: math.MaxFloat64},
		{name: "smallest finite divisor", left: math.SmallestNonzeroFloat64, right: math.SmallestNonzeroFloat64, want: 1},
		{name: "underflow to zero", left: math.SmallestNonzeroFloat64, right: 2, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Divide(tt.left, tt.right)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivideFloatingPointDecimals(t *testing.T) {
	got, err := calculator.Divide(1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-1.0/3.0) > 1e-16 {
		t.Errorf("result = %.17g, want approximately one third", got)
	}
}

func TestDivideByZero(t *testing.T) {
	tests := []struct {
		name        string
		left, right float64
	}{
		{name: "positive zero divisor", left: 8, right: 0},
		{name: "negative zero divisor", left: 8, right: math.Copysign(0, -1)},
		{name: "zero divided by positive zero", left: 0, right: 0},
		{name: "zero divided by negative zero", left: 0, right: math.Copysign(0, -1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calculator.Divide(tt.left, tt.right)
			if !errors.Is(err, calculator.ErrDivisionByZero) {
				t.Errorf("error = %v, want ErrDivisionByZero", err)
			}
		})
	}
}

func TestDivideInvalidOperands(t *testing.T) {
	for _, operands := range [][2]float64{
		{math.NaN(), 1}, {1, math.NaN()},
		{math.Inf(1), 1}, {1, math.Inf(-1)},
	} {
		_, err := calculator.Divide(operands[0], operands[1])
		if !errors.Is(err, calculator.ErrInvalidOperand) {
			t.Errorf("Divide(%v, %v) error = %v, want ErrInvalidOperand", operands[0], operands[1], err)
		}
	}
}

func TestDivideOverflow(t *testing.T) {
	for _, left := range []float64{math.MaxFloat64, -math.MaxFloat64} {
		_, err := calculator.Divide(left, 0.5)
		if !errors.Is(err, calculator.ErrNonFiniteResult) {
			t.Errorf("Divide(%v, 0.5) error = %v, want ErrNonFiniteResult", left, err)
		}
	}
}
