package domain

import "errors"

var (
	ErrInvalidAccountID      = errors.New("Field account_id is invalid")
	ErrZeroAmount            = errors.New("Field amount cannot be exactly zero")
	ErrInvalidNegativeAmount = errors.New("Field amount must be negative for this operation")
	ErrInvalidPositiveAmount = errors.New("Field amount must be positive for this operation")
	ErrUnknownOperationType  = errors.New("Unknown operation type")
)

type OperationType int

const (
	OpNormalPurchase      OperationType = 1
	OpPurchaseInstallment OperationType = 2
	OpWithdrawal          OperationType = 3
	OpCreditVoucher       OperationType = 4
)

type Transaction struct {
	ID              int           `db:"transaction_id" json:"transaction_id"`
	AccountID       int           `db:"account_id" json:"account_id"`
	OperationTypeID OperationType `db:"operation_type_id" json:"operation_type_id"`
	Amount          Money         `db:"amount" json:"amount"`
	EventDate       string        `db:"event_date" json:"event_date"`
}

func (t *Transaction) Validate() error {
	if t.AccountID <= 0 {
		return ErrInvalidAccountID
	}
	if t.Amount == 0 {
		return ErrZeroAmount
	}

	switch t.OperationTypeID {
	case OpNormalPurchase, OpPurchaseInstallment, OpWithdrawal:
		if t.Amount > 0 {
			return ErrInvalidNegativeAmount
		}
	case OpCreditVoucher:
		if t.Amount < 0 {
			return ErrInvalidPositiveAmount
		}
	default:
		return ErrUnknownOperationType
	}
	return nil
}
