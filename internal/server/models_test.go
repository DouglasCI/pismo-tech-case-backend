package server

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
			err := tc.account.ValidateAccount()

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Test '%s' failed: expected error %v, got %v", tc.name, tc.expectedError, err)
			}
		})
	}
}

func TestValidateTransaction(t *testing.T) {
	tests := []struct {
		name          string
		transaction   Transaction
		expectedError error
	}{
		{
			name: "Valid Normal Purchase (Negative Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: NormalPurchase,
				Amount:          -10.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Normal Purchase (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: NormalPurchase,
				Amount:          10.0,
			},
			expectedError: ErrInvalidNegativeAmount,
		},
		{
			name: "Valid Purchase With Installments (Negative Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: PurchaseInstallment,
				Amount:          -20.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Purchase With Installments (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: PurchaseInstallment,
				Amount:          20.0,
			},
			expectedError: ErrInvalidNegativeAmount,
		},
		{
			name: "Valid Withdrawal (Negative Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: PurchaseInstallment,
				Amount:          -30.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Withdrawal (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: PurchaseInstallment,
				Amount:          30.0,
			},
			expectedError: ErrInvalidNegativeAmount,
		},
		{
			name: "Valid Credit Voucher (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: CreditVoucher,
				Amount:          40.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Credit Voucher (Negative Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: CreditVoucher,
				Amount:          -40.0,
			},
			expectedError: ErrInvalidPositiveAmount,
		},
		{
			name: "Missing Account ID",
			transaction: Transaction{
				AccountID:       0,
				OperationTypeID: Withdrawal,
				Amount:          -10.0,
			},
			expectedError: ErrInvalidAccountID,
		},
		{
			name: "Amount Exactly Zero",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: NormalPurchase,
				Amount:          0.0,
			},
			expectedError: ErrZeroAmount,
		},
		{
			name: "Unknown Operation Type",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: 999,
				Amount:          10.0,
			},
			expectedError: ErrUnknownOperationType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.transaction.ValidateTransaction()

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Test '%s' failed: expected error %v, got %v", tc.name, tc.expectedError, err)
			}
		})
	}
}
