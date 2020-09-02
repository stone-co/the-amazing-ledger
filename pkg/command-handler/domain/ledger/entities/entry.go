package entities

import "time"

type EntryType string

type Entry struct {
	ID            string
	AccountID     string
	TransactionID string
	RequestID     string
	Amount        int
	BalanceAfter  int
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
