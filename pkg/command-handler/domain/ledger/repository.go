package ledger

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type Repository interface {
	CreateAccount(*entities.Account) (entities.Account, error)
	CreateTransaction(context.Context, *entities.Transaction) error
	LoadObjectsIntoCache(ctx context.Context, objects *entities.CachedAccounts) (entities.Version, error)
	GetAccountBalance(ctx context.Context, accountName entities.AccountName) (*entities.AccountBalance, error)
}
