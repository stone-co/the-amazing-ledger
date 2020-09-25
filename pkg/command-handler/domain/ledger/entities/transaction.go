package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidData          = errors.New("invalid data")
	ErrInvalidEntriesNumber = errors.New("invalid entries number")
	ErrInvalidBalance       = errors.New("invalid balance")
	ErrIdempotencyKey       = errors.New("idempotency key violation")
)

type Transaction struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Entries   []Entry
}

func NewTransaction(id uuid.UUID, createdAt time.Time, entries ...Entry) (*Transaction, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidData
	}

	if createdAt.IsZero() {
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
		ID:        id,
		CreatedAt: createdAt,
		Entries:   entries,
	}

	return t, nil
}
