package accounts

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"

type Repository interface {
	Create(*entities.Account) error
	Get(id string) (entities.Account, error)
	Search(accountType string, accountOwnerID string, accountOwner string, accountName string, accountMetadata string) (entities.Account, error)
	Update(id string, balance int) error
}
