package entities

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID      uuid.UUID
	Entries []Entry
}

func NewTransaction(id uuid.UUID, entries ...Entry) (*Transaction, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidData
	}

	if len(entries) <= 1 {
		return nil, ErrInvalidEntriesNumber
	}

	balance := 0
	for _, entry := range entries {
		if entry.Operation == DebitOperation {
			balance += entry.Amount
		} else {
			balance -= entry.Amount
		}
	}

	if balance != 0 {
		return nil, ErrInvalidBalance
	}

	t := &Transaction{
		ID:      id,
		Entries: entries,
	}

	return t, nil
}
