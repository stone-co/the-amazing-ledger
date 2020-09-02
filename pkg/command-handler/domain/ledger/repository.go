package ledger

import (
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type Repository interface {
	CreateAccount(*entities.Account) error
	GetAccount(id string) (entities.Account, error)
	SearchAccount(accountType string, accountOwnerID string, accountOwner string, accountName string, accountMetadata string) (entities.Account, error)
	UpdateBalance(id string, balance int) error
	CreateTransaction(*[]entities.Entry) error
}
