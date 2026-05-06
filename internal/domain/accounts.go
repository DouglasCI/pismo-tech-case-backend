package domain

import (
	"errors"
	"strings"
)

var ErrMissingDocumentNumber = errors.New("Field document_number is required")

type Account struct {
	ID             int    `db:"account_id" json:"account_id"`
	DocumentNumber string `db:"document_number" json:"document_number"`
}

func (a *Account) Validate() error {
	if strings.TrimSpace(a.DocumentNumber) == "" {
		return ErrMissingDocumentNumber
	}
	return nil
}
