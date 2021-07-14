package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) ListAccountEntries(ctx context.Context, req vos.AccountEntryRequest) (vos.AccountEntryResponse, error) {
	entries, nextPage, err := l.repository.ListAccountEntries(ctx, req)
	if err != nil {
		return vos.AccountEntryResponse{}, fmt.Errorf("failed to get account entries: %w", err)
	}

	return vos.AccountEntryResponse{
		Entries:  entries,
		NextPage: nextPage,
	}, nil
}
