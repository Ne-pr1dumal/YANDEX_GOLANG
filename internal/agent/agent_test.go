package agent

import (
	"errors"
	"math"
	"testing"
)

func TestCalculations(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		a           float64
		b           float64
		expected    float64
		expectErr   bool
		expectedErr error
	}{
		// Valid cases
		{
			name:      "Add positive integers",
			operation: "+",
			a:         10,
			b:         20,
			expected:  30,
		},
		{
			name:      "Add decimal numbers",
			operation: "+",
			a:         2.5,
			b:         3.75,
			expected:  6.25,
		},
		{
			name:      "Subtract larger number",
			operation: "-",
			a:         5,
			b:         15,
			expected:  -10,
		},
		{
			name:      "Multiply by negative",
			operation: "*",
			a:         -4,
			b:         25,
			expected:  -100,
		},
		{
			name:      "Divide by small decimal",
			operation: "/",
			a:         100,
			b:         0.1,
			expected:  1000,
		},

		// Edge cases
		{
			name:      "Multiply by zero",
			operation: "*",
			a:         math.Pi,
			b:         0,
			expected:  0,
		},
		{
			name:      "Divide zero by number",
			operation: "/",
			a:         0,
			b:         5,
			expected:  0,
		},
		{
			name:        "Divide by zero",
			operation:   "/",
			a:           1,
			b:           0,
			expectErr:   true,
			expectedErr: ErrDivisionByZero,
		},

		// Error cases
		{
			name:        "Invalid operator symbol",
			operation:   "^",
			a:           2,
			b:           3,
			expectErr:   true,
			expectedErr: ErrInvalidOperator,
		},
		{
			name:        "Empty operator",
			operation:   "",
			a:           1,
			b:           1,
			expectErr:   true,
			expectedErr: ErrInvalidOperator,
		},
		{
			name:        "Multi-character operator",
			operation:   "++",
			a:           1,
			b:           1,
			expectErr:   true,
			expectedErr: ErrInvalidOperator,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculate(tt.operation, tt.a, tt.b)

			if tt.expectErr {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}

				if !errors.Is(err, tt.expectedErr) && err.Error() != tt.expectedErr.Error() {
					t.Errorf("Error mismatch:\nExpected: %v\nGot:      %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Result mismatch:\nExpected: %.10f\nGot:      %.10f", tt.expected, result)
			}
		})
	}
}

func TestCalculateExpression(t *testing.T) {
	t.Run("NotImplemented check", func(t *testing.T) {
		_, err := CalculateExpression("2+2")
		if err == nil || err.Error() != "not implemented" {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Empty expression", func(t *testing.T) {
		_, err := CalculateExpression("")
		if err == nil {
			t.Error("Expected error for empty expression")
		}
	})
}
