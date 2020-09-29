package entities

import (
	"github.com/google/uuid"
)

type VersionType uint64

const (
	AnyAccountVersion VersionType = 0
	NewAccount        VersionType = 1
)

type Entry struct {
	ID        uuid.UUID
	Operation OperationType
	AccountID uuid.UUID
	Version   VersionType
	Amount    int
}

func NewEntry(id uuid.UUID, operation OperationType, accountID uuid.UUID, version VersionType, amount int) *Entry {
	return &Entry{
		ID:        id,
		Operation: operation,
		AccountID: accountID,
		Version:   version,
		Amount:    amount,
	}
}
