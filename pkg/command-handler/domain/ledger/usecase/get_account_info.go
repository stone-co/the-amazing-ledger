package usecase

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (l *LedgerUseCase) GetAccountInfo(ctx context.Context, accountPath string) (*entities.AccountInfo, error) {

	accountName, err := entities.NewAccountName(accountPath)

	if err != nil {
		return nil, err
	}

	accountInfo, err := l.repository.GetAccountInfo(ctx, accountName)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	return accountInfo, nil
}
