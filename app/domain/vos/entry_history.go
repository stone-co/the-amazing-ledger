package vos

import (
	"time"

	"github.com/stone-co/the-amazing-ledger/app"
)

type EntryHistory struct {
	Amount    int
	Operation OperationType
	CreatedAt time.Time
}

func NewEntryHistory(operation OperationType, amount int, createdAt time.Time) (*EntryHistory, error) {
	if operation == InvalidOperation {
		return nil, app.ErrInvalidOperation
	}

	if amount <= 0 {
		return nil, app.ErrInvalidAmount
	}

	return &EntryHistory{
		Operation: operation,
		Amount:    amount,
		CreatedAt: createdAt,
	}, nil
}
