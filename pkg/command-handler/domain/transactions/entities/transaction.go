package entities

import "time"

const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

type TransactionType string

type Transaction struct {
	ID           string
	AccountID    string
	OperationID  string
	RequestID    string
	Type         TransactionType
	Amount       int
	BalanceAfter int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
