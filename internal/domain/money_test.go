package domain

import (
	"errors"
	"testing"
)

func TestMoney_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		expectedValue Money
		expectedError error
	}{
		{
			name:          "Valid integer amount",
			payload:       `150`,
			expectedValue: 15000,
			expectedError: nil,
		},
		{
			name:          "Valid decimal amount",
			payload:       `10.50`,
			expectedValue: 1050,
			expectedError: nil,
		},
		{
			name:          "Valid negative decimal amount",
			payload:       `-20.12`,
			expectedValue: -2012,
			expectedError: nil,
		},
		{
			name:          "Valid single decimal place",
			payload:       `10.5`,
			expectedValue: 1050,
			expectedError: nil,
		},
		{
			name:          "Invalid amount with 3 decimal places",
			payload:       `10.452`,
			expectedValue: 0,
			expectedError: ErrAmountDecimalPlaces,
		},
		{
			name:          "Invalid string payload",
			payload:       `"abc"`,
			expectedValue: 0,
			expectedError: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Money
			err := m.UnmarshalJSON([]byte(tt.payload))

			if tt.expectedError == ErrAmountDecimalPlaces {
				if !errors.Is(err, ErrAmountDecimalPlaces) {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if tt.expectedError != nil && err == nil {
				t.Errorf("expected an error, but got nil")
				return
			}

			if err == nil && m != tt.expectedValue {
				t.Errorf("expected value %d, got %d", tt.expectedValue, m)
			}
		})
	}
}

func TestMoney_MarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		moneyValue     Money
		expectedOutput string
	}{
		{
			name:           "Marshal positive value",
			moneyValue:     1050,
			expectedOutput: `10.50`,
		},
		{
			name:           "Marshal negative value",
			moneyValue:     -2012,
			expectedOutput: `-20.12`,
		},
		{
			name:           "Marshal zero value",
			moneyValue:     0,
			expectedOutput: `0.00`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.moneyValue.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error marshaling JSON: %v", err)
			}

			if string(data) != tt.expectedOutput {
				t.Errorf("expected %s, got %s", tt.expectedOutput, string(data))
			}
		})
	}
}
