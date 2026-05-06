package domain

import (
	"errors"
	"testing"
)

func TestValidateAccount(t *testing.T) {
	tests := []struct {
		name          string
		account       Account
		expectedError error
	}{
		{
			name:          "Valid Account",
			account:       Account{DocumentNumber: "12345678900"},
			expectedError: nil,
		},
		{
			name:          "Empty Document Number",
			account:       Account{DocumentNumber: ""},
			expectedError: ErrMissingDocumentNumber,
		},
		{
			name:          "Whitespace Document Number",
			account:       Account{DocumentNumber: "   "},
			expectedError: ErrMissingDocumentNumber,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.account.Validate()

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Test '%s' failed: expected error %v, got %v", tc.name, tc.expectedError, err)
			}
		})
	}
}
