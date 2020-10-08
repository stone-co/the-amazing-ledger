package entities

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidData          = errors.New("invalid data")
	ErrInvalidEntriesNumber = errors.New("invalid entries number")
	ErrInvalidBalance       = errors.New("invalid balance")
	ErrIdempotencyKey       = errors.New("idempotency key violation")
	ErrInvalidVersion       = errors.New("invalid version")
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
		if entry.Amount <= 0 {
			return nil, ErrInvalidData
		}

		if entry.Operation == InvalidOperation {
			return nil, ErrInvalidData
		}

		if entry.Operation == InvalidOperation {
			return nil, ErrInvalidData
		}

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
