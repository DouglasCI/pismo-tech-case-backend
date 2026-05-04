package server

import (
	"errors"
	"strings"
)

// Field Validation Errors
var (
	ErrMissingDocumentNumber = errors.New("Field document_number is required")
	ErrInvalidAccountID      = errors.New("Field account_id is invalid")
	ErrZeroAmount            = errors.New("Field amount cannot be exactly zero")
)

// Business Rules Errors
var (
	ErrInvalidNegativeAmount = errors.New("Field amount must be negative for purchases and withdrawals")
	ErrInvalidPositiveAmount = errors.New("Field amount must be positive for credit vouchers")
	ErrUnknownOperationType  = errors.New("Unknown operation type")
)

type OperationType int

const (
	OpNormalPurchase      = 1
	OpPurchaseInstallment = 2
	OpWithdrawal          = 3
	OpCreditVoucher       = 4
)

// Account represents the user account model
type Account struct {
	ID             int    `db:"account_id" json:"account_id"`
	DocumentNumber string `db:"document_number" json:"document_number"`
}

func (a *Account) ValidateAccount() error {
	if strings.TrimSpace(a.DocumentNumber) == "" {
		return ErrMissingDocumentNumber
	}
	return nil
}

// Transaction represents the transaction model
type Transaction struct {
	ID              int           `db:"transaction_id" json:"transaction_id"`
	AccountID       int           `db:"account_id" json:"account_id"`
	OperationTypeID OperationType `db:"operation_type_id" json:"operation_type_id"`
	Amount          float64       `db:"amount" json:"amount"`
	EventDate       string        `db:"event_date" json:"event_date"`
}

func (t *Transaction) validateFields() error {
	if t.AccountID <= 0 {
		return ErrInvalidAccountID
	}
	if t.Amount == 0 {
		return ErrZeroAmount
	}
	return nil
}

func (t *Transaction) validateBusinessRules() error {
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

func (t *Transaction) ValidateTransaction() error {
	if err := t.validateFields(); err != nil {
		return err
	}

	if err := t.validateBusinessRules(); err != nil {
		return err
	}

	return nil
}
