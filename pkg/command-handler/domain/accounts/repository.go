package accounts

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"

type Repository interface {
	Create(*entities.Account) error
	Get(*string) (entities.Account, error)
	Update(*string, *int) error
}
