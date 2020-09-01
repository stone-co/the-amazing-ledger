package accounts

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"

type RepositoryMock struct {
	OnCreate func(*entities.Account) error
	OnGet    func(string) (entities.Account, error)
	OnSearch func(string, string, string, string, string) (entities.Account, error)
	OnUpdate func(string, int) error
}

func (s RepositoryMock) Create(a *entities.Account) error {
	return s.OnCreate(a)
}

func (s RepositoryMock) Get(id string) (entities.Account, error) {
	return s.OnGet(id)
}

func (s RepositoryMock) Search(accountType string, accountOwnerID string, accountOwner string, accountName string, accountMetadata string) (entities.Account, error) {
	return s.OnSearch(accountType, accountOwnerID, accountOwner, accountName, accountMetadata)
}

func (s RepositoryMock) Update(id string, balance int) error {
	return s.OnUpdate(id, balance)
}
