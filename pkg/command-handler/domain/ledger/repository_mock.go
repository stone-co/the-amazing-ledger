package ledger

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type RepositoryMock struct {
	OnCreateAccount        func(*entities.Account) (entities.Account, error)
	OnGetAccount           func(string) (entities.Account, error)
	OnSearchAccount        func(*entities.Account) (entities.Account, error)
	OnUpdateBalance        func(string, int) error
	OnCreateTransaction    func(context.Context, *entities.Transaction) error
	OnLoadObjectsIntoCache func(ctx context.Context, cachedAccounts *entities.CachedAccounts) (entities.Version, error)
	OnGetAccountBalance    func(ctx context.Context, accountName entities.AccountName) (*entities.AccountBalance, error)
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

func (s RepositoryMock) CreateTransaction(ctx context.Context, transaction *entities.Transaction) error {
	return s.OnCreateTransaction(ctx, transaction)
}

func (s RepositoryMock) LoadObjectsIntoCache(ctx context.Context, cachedAccounts *entities.CachedAccounts) (entities.Version, error) {
	return s.OnLoadObjectsIntoCache(ctx, cachedAccounts)
}

func (s RepositoryMock) GetAccountBalance(ctx context.Context, accountName entities.AccountName) (*entities.AccountBalance, error) {
	return s.OnGetAccountBalance(ctx, accountName)
}
