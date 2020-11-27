package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func TestLedgerUseCase_CreateTransaction(t *testing.T) {
	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	t.Run("Successfully creates a transaction with minimum inputs", func(t *testing.T) {
		e1, _ := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.AnyAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.AnyAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := newFakeCreateTransactionUseCase(nil).CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with new accounts", func(t *testing.T) {
		e1, _ := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := newFakeCreateTransactionUseCase(nil).CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with expected version", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
	})

	t.Run("Fail with invalid expected version", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 4, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, entities.ErrInvalidVersion.Is(err))
	})

	t.Run("When transaction fail, the counter isn't incremented", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)

		lastVersion := useCase.GetLastVersion()

		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 4, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, entities.ErrInvalidVersion.Is(err))

		assert.Equal(t, lastVersion, useCase.GetLastVersion())
	})

	t.Run("Global object version does not change when the transaction fails", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		lastVersion := useCase.GetLastVersion()
		e1, _ := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}
		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, entities.ErrInvalidVersion.Is(err))
		assert.Equal(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())
	})

	t.Run("Object version does not change when the idempotency fails, but global counter is changed", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}
		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)

		lastVersion := useCase.GetLastVersion()
		idempotencyKey := uuid.New()
		e1, _ = entities.NewEntry(idempotencyKey, entities.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(idempotencyKey, entities.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		useCase.repository = &ledger.RepositoryMock{
			OnCreateTransaction: func(context.Context, *entities.Transaction) error {
				return entities.ErrIdempotencyKey
			},
		}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, entities.ErrIdempotencyKey.Is(err))
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		e1, _ = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		useCase.repository = &ledger.RepositoryMock{
			OnCreateTransaction: func(context.Context, *entities.Transaction) error {
				return nil
			},
		}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
	})
}
