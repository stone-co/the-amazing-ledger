package usecase

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stone-co/the-amazing-ledger/app/domain/errors"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, accountName vo.AccountName) (*vo.AccountBalance, error) {

	accountBalance, err := l.repository.GetAccountBalance(ctx, accountName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrAccountNotFound
		}
		return nil, err
	}

	return accountBalance, nil
}
