package entities

import (
	"github.com/google/uuid"
)

type Entry struct {
	ID        uuid.UUID
	Operation OperationType
	AccountID string
	Version   Version
	Amount    int
}

func NewEntry(id uuid.UUID, operation OperationType, accountID string, version Version, amount int) *Entry {
	return &Entry{
		ID:        id,
		Operation: operation,
		AccountID: accountID,
		Version:   version,
		Amount:    amount,
	}
}
