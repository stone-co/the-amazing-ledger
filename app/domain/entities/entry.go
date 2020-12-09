package entities

import (
	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/errors"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

type Entry struct {
	ID        uuid.UUID
	Operation vo.OperationType
	Account   *vo.AccountName
	Version   vo.Version
	Amount    int
}

func NewEntry(id uuid.UUID, operation vo.OperationType, accountID string, version vo.Version, amount int) (*Entry, error) {
	if id == uuid.Nil {
		return nil, errors.ErrInvalidData.Cause("id")
	}

	if operation == vo.InvalidOperation {
		return nil, errors.ErrInvalidData.Cause("operation")
	}

	if amount <= 0 {
		return nil, errors.ErrInvalidData.Cause("amount")
	}

	acc, err := vo.NewAccountName(accountID)
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
