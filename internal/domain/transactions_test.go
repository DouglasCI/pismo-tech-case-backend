package domain

import (
	"errors"
	"testing"
)

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
				OperationTypeID: OpNormalPurchase,
				Amount:          -10.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Normal Purchase (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpNormalPurchase,
				Amount:          10.0,
			},
			expectedError: ErrInvalidNegativeAmount,
		},
		{
			name: "Valid Purchase With Installments (Negative Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpPurchaseInstallment,
				Amount:          -20.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Purchase With Installments (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpPurchaseInstallment,
				Amount:          20.0,
			},
			expectedError: ErrInvalidNegativeAmount,
		},
		{
			name: "Valid Withdrawal (Negative Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpPurchaseInstallment,
				Amount:          -30.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Withdrawal (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpPurchaseInstallment,
				Amount:          30.0,
			},
			expectedError: ErrInvalidNegativeAmount,
		},
		{
			name: "Valid Credit Voucher (Positive Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpCreditVoucher,
				Amount:          40.0,
			},
			expectedError: nil,
		},
		{
			name: "Invalid Credit Voucher (Negative Amount)",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpCreditVoucher,
				Amount:          -40.0,
			},
			expectedError: ErrInvalidPositiveAmount,
		},
		{
			name: "Missing Account ID",
			transaction: Transaction{
				AccountID:       0,
				OperationTypeID: OpWithdrawal,
				Amount:          -10.0,
			},
			expectedError: ErrInvalidAccountID,
		},
		{
			name: "Amount Exactly Zero",
			transaction: Transaction{
				AccountID:       1,
				OperationTypeID: OpNormalPurchase,
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
			err := tc.transaction.Validate()

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Test '%s' failed: expected error %v, got %v", tc.name, tc.expectedError, err)
			}
		})
	}
}
