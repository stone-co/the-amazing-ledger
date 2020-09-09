package ledger

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"

type RepositoryMock struct {
	OnCreateAccount     func(*entities.Account) (entities.Account, error)
	OnGetAccount        func(string) (entities.Account, error)
	OnSearchAccount     func(*entities.Account) (entities.Account, error)
	OnUpdateBalance     func(string, int) error
	OnCreateTransaction func(*[]entities.Entry) error
}

func (s RepositoryMock) CreateAccount(a *entities.Account) (entities.Account, error) {
	return s.OnCreateAccount(a)
}

func (s RepositoryMock) GetAccount(id string) (entities.Account, error) {
	return s.OnGetAccount(id)
}

func (s RepositoryMock) SearchAccount(a *entities.Account) (entities.Account, error) {
	return s.OnSearchAccount(a)
}

func (s RepositoryMock) UpdateBalance(id string, balance int) error {
	return s.OnUpdateBalance(id, balance)
}

func (s RepositoryMock) CreateTransaction(entries *[]entities.Entry) error {
	return s.OnCreateTransaction(entries)
}
