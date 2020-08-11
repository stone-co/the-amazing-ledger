package entities

import "time"

type TransactionType string

type Transaction struct {
	ID           string
	AccountID    string
	OperationID  string
	RequestID    string
	Amount       int
	BalanceAfter int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
