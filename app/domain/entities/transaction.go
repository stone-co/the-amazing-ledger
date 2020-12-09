package entities

import (
	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/errors"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

type Transaction struct {
	ID      uuid.UUID
	Entries []Entry
}

func NewTransaction(id uuid.UUID, entries ...Entry) (*Transaction, error) {
	if id == uuid.Nil {
		return nil, errors.ErrInvalidData
	}

	if len(entries) <= 1 {
		return nil, errors.ErrInvalidEntriesNumber
	}

	balance := 0
	for _, entry := range entries {
		if entry.Operation == vo.DebitOperation {
			balance += entry.Amount
		} else {
			balance -= entry.Amount
		}
	}

	if balance != 0 {
		return nil, errors.ErrInvalidBalance
	}

	t := &Transaction{
		ID:      id,
		Entries: entries,
	}

	return t, nil
}
