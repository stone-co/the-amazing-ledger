package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type Transaction struct {
	ID             uuid.UUID
	Entries        []Entry
	Event          uint32
	Company        string
	CompetenceDate time.Time
}

func NewTransaction(id uuid.UUID, entries ...Entry) (Transaction, error) {
	if id == uuid.Nil {
		return Transaction{}, app.ErrInvalidTransactionID
	}

	if len(entries) <= 1 {
		return Transaction{}, app.ErrInvalidEntriesNumber
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
		return Transaction{}, app.ErrInvalidBalance
	}

	t := Transaction{
		ID:      id,
		Entries: entries,
	}

	return t, nil
}
