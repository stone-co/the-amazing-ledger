package vos

import (
	"github.com/stone-co/the-amazing-ledger/app"
)

type Entry struct {
	Account   string
	Operation OperationType
	Amount    int
}

func NewEntry(account string, operation OperationType, amount int) (*Entry, error) {
	if operation == InvalidOperation {
		return nil, app.ErrInvalidOperation
	}

	if amount <= 0 {
		return nil, app.ErrInvalidAmount
	}

	return &Entry{
		Account:   account,
		Operation: operation,
		Amount:    amount,
	}, nil
}

// type Entry struct {
// 	Account   AccountName
// 	Operation OperationType
// 	Amount    int
// }

// func NewEntry(account string, operation OperationType, amount int) (*Entry, error) {
// 	if operation == InvalidOperation {
// 		return nil, app.ErrInvalidOperation
// 	}

// 	if amount <= 0 {
// 		return nil, app.ErrInvalidAmount
// 	}

// 	acc, err := NewAccountName(account)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Entry{
// 		Account:   *acc,
// 		Operation: operation,
// 		Amount:    amount,
// 	}, nil
// }
