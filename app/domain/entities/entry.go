package entities

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type Entry struct {
	ID        uuid.UUID
	Operation vos.OperationType
	Account   vos.AccountPath
	Version   vos.Version
	Amount    int
}

func NewEntry(id uuid.UUID, operation vos.OperationType, accountID string, version vos.Version, amount int) (Entry, error) {
	if id == uuid.Nil {
		return Entry{}, app.ErrInvalidEntryID
	}

	if operation == vos.InvalidOperation {
		return Entry{}, app.ErrInvalidOperation
	}

	if amount <= 0 {
		return Entry{}, app.ErrInvalidAmount
	}

	acc, err := vos.NewAccountPath(accountID)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to create account path: %w", err)
	}

	return Entry{
		ID:        id,
		Operation: operation,
		Account:   acc,
		Version:   version,
		Amount:    amount,
	}, nil
}
