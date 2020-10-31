package entities

import (
	"github.com/google/uuid"
)

type Entry struct {
	ID        uuid.UUID
	Operation OperationType
	Account   *AccountName
	Version   Version
	Amount    int
}

func NewEntry(id uuid.UUID, operation OperationType, accountID string, version Version, amount int) (*Entry, error) {

	if id == uuid.Nil {
		return nil, ErrInvalidData
	}

	if operation != CreditOperation && operation != DebitOperation {
		return nil, ErrInvalidData
	}

	if amount <= 0 {
		return nil, ErrInvalidData
	}

	acc, err := NewAccountName(accountID)
	if err != nil {
		return nil, err
	}

	return &Entry{
		ID:        id,
		Operation: operation,
		Account:   acc,
		Version:   version,
		Amount:    amount,
	}, nil
}
