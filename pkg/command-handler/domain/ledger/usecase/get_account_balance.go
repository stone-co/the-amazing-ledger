package usecase

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, accountName entities.AccountName) (*entities.AccountBalance, error) {

	accountBalance, err := l.repository.GetAccountBalance(ctx, accountName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrAccountNotFound
		}
		return nil, err
	}

	return accountBalance, nil
}
