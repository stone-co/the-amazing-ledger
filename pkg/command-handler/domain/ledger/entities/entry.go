package entities

import (
	"errors"

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
		return nil, ErrInvalidData.cause(errors.New("id"))
	}

	if operation == InvalidOperation {
		return nil, ErrInvalidData.cause(errors.New("operation"))
	}

	if amount <= 0 {
		return nil, ErrInvalidData.cause(errors.New("amount"))
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
