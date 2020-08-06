package accounts

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"

type RepositoryMock struct {
	OnCreate func(*entities.Account) error
	OnGet    func(*string) (entities.Account, error)
}

func (s RepositoryMock) Create(a *entities.Account) error {
	return s.OnCreate(a)
}

func (s RepositoryMock) Get(id *string) (entities.Account, error) {
	return s.OnGet(id)
}
