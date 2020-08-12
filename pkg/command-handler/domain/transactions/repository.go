package transactions

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/transactions/entities"

type Repository interface {
	Create(*[]entities.Transaction) error
}
