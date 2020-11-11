package ledger

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type Repository interface {
	CreateAccount(*entities.Account) (entities.Account, error)
	GetAccount(id string) (entities.Account, error)
	SearchAccount(*entities.Account) (entities.Account, error)
	UpdateBalance(id string, balance int) error
	CreateTransaction(context.Context, *entities.Transaction) error
	LoadObjectsIntoCache(ctx context.Context, objects *entities.CachedAccounts) (entities.Version, error)
	GetAccountInfo(ctx context.Context, accountName *entities.AccountName) (*entities.AccountInfo, error)
}
