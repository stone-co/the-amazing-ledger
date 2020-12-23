package vos

import (
	"github.com/stone-co/the-amazing-ledger/app"
)

type Statement struct {
	Account   string
	Operation OperationType
	Amount    int
}

func NewStatement(account string, operation OperationType, amount int) (*Statement, error) {
	if operation == InvalidOperation {
		return nil, app.ErrInvalidOperation
	}

	if amount <= 0 {
		return nil, app.ErrInvalidAmount
	}

	return &Statement{
		Account:   account,
		Operation: operation,
		Amount:    amount,
	}, nil
}
