package entities

import (
	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type Transaction struct {
	ID      uuid.UUID
	Entries []Entry
}

func NewTransaction(id uuid.UUID, entries ...Entry) (*Transaction, error) {
	if id == uuid.Nil {
		return nil, app.ErrInvalidTransactionID
	}

	if len(entries) <= 1 {
		return nil, app.ErrInvalidEntriesNumber
	}

	balance := 0
	for _, entry := range entries {
		if entry.Operation == vos.DebitOperation {
			balance += entry.Amount
		} else {
			balance -= entry.Amount
		}
	}

	if balance != 0 {
		return nil, app.ErrInvalidBalance
	}

	t := &Transaction{
		ID:      id,
		Entries: entries,
	}

	return t, nil
}
