package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
)

func TestLedgerUseCase_CreateTransaction(t *testing.T) {
	accountID1 := testdata.GenerateAccountPath()
	accountID2 := testdata.GenerateAccountPath()

	testCases := []struct {
		name        string
		repoSetup   *mocks.RepositoryMock
		entries     func(t *testing.T) []entities.Entry
		expectedErr error
	}{
		{
			name: "Should create a transaction successfully",
			repoSetup: &mocks.RepositoryMock{
				CreateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) error {
					return nil
				},
			},
			entries: func(t *testing.T) []entities.Entry {
				e1, err := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NextAccountVersion, 123)
				assert.NoError(t, err)

				e2, err := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NextAccountVersion, 123)
				assert.NoError(t, err)

				return []entities.Entry{e1, e2}
			},
			expectedErr: nil,
		},
		{
			name: "Should return an error if entry version is invalid",
			repoSetup: &mocks.RepositoryMock{
				CreateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) error {
					return app.ErrInvalidVersion
				},
			},
			entries: func(t *testing.T) []entities.Entry {
				e1, err := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.Version(1), 123)
				assert.NoError(t, err)

				e2, err := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.Version(3), 123)
				assert.NoError(t, err)

				return []entities.Entry{e1, e2}
			},
			expectedErr: app.ErrInvalidVersion,
		},
		{
			name: "Should return an error if violates idempotency key",
			repoSetup: &mocks.RepositoryMock{
				CreateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) error {
					return app.ErrIdempotencyKeyViolation
				},
			},
			entries: func(t *testing.T) []entities.Entry {
				idempotencyKey := uuid.New()

				e1, err := entities.NewEntry(idempotencyKey, vos.DebitOperation, accountID1, vos.NextAccountVersion, 123)
				assert.NoError(t, err)

				e2, err := entities.NewEntry(idempotencyKey, vos.CreditOperation, accountID2, vos.NextAccountVersion, 123)
				assert.NoError(t, err)

				return []entities.Entry{e1, e2}
			},
			expectedErr: app.ErrIdempotencyKeyViolation,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewLedgerUseCase(logrus.New(), tt.repoSetup)

			tx, err := entities.NewTransaction(uuid.New(), 1, "abc", time.Now(), tt.entries(t)...)
			assert.NoError(t, err)

			err = usecase.CreateTransaction(context.Background(), tx)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
