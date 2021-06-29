package usecases

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
)

func TestLedgerUseCase_QueryAggregatedBalance(t *testing.T) {
	t.Run("should return aggregated balance successfully", func(t *testing.T) {
		query, err := vos.NewAccountQuery("liability.stone.clients.*")
		assert.NoError(t, err)

		queryBalance := vos.NewQueryBalance(query, 20)
		mockedRepository := &mocks.RepositoryMock{
			QueryAggregatedBalanceFunc: func(ctx context.Context, account vos.AccountQuery) (vos.QueryBalance, error) {
				return queryBalance, nil
			},
		}
		usecase := NewLedgerUseCase(logrus.New(), mockedRepository)

		got, err := usecase.QueryAggregatedBalance(context.Background(), query)
		assert.NoError(t, err)
		assert.Equal(t, queryBalance.Balance, got.Balance)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		query, err := vos.NewAccountQuery("liability.stone.clients.*")
		assert.NoError(t, err)

		mockedRepository := &mocks.RepositoryMock{
			QueryAggregatedBalanceFunc: func(ctx context.Context, account vos.AccountQuery) (vos.QueryBalance, error) {
				return vos.QueryBalance{}, app.ErrAccountNotFound
			},
		}
		usecase := NewLedgerUseCase(logrus.New(), mockedRepository)

		got, err := usecase.QueryAggregatedBalance(context.Background(), query)
		assert.Empty(t, got)
		assert.ErrorIs(t, err, app.ErrAccountNotFound)
	})
}
