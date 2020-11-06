package usecase

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (l *LedgerUseCase) GetAccountInfo(ctx context.Context, accountID string) (*entities.AccountInfo, error) {

	accountInfo, err := l.repository.GetAccountInfo(ctx, accountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	return accountInfo, nil
}
