package usecases

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {

	accountBalance, err := l.repository.GetAccountBalance(ctx, accountName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, app.ErrAccountNotFound
		}
		return nil, err
	}

	return accountBalance, nil
}
