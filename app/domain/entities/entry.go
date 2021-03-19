package entities

import (
	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type Entry struct {
	ID        uuid.UUID
	Operation vos.OperationType
	Account   *vos.AccountName
	Version   vos.Version
	Amount    int
}

func NewEntry(id uuid.UUID, operation vos.OperationType, accountID string, version vos.Version, amount int) (*Entry, error) {
	if id == uuid.Nil {
		return nil, app.ErrInvalidEntryID
	}

	if operation == vos.InvalidOperation {
		return nil, app.ErrInvalidOperation
	}

	if amount <= 0 {
		return nil, app.ErrInvalidAmount
	}

	acc, err := vos.NewAccountName(accountID)
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
