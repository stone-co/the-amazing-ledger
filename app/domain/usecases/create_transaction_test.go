package usecases

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/mocks"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_CreateTransaction(t *testing.T) {
	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	t.Run("Successfully creates a transaction with minimum inputs", func(t *testing.T) {
		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.AnyAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.AnyAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := newFakeCreateTransactionUseCase(nil).CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with new accounts", func(t *testing.T) {
		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := newFakeCreateTransactionUseCase(nil).CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with expected version", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
	})

	t.Run("Fail with invalid expected version", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 4, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, app.ErrInvalidVersion.Is(err))
	})

	t.Run("When transaction fail, the counter isn't incremented", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)

		assert.Nil(t, err)

		lastVersion := useCase.GetLastVersion()

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 4, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, app.ErrInvalidVersion.Is(err))

		assert.Equal(t, lastVersion, useCase.GetLastVersion())
	})

	t.Run("Global object version does not change when the transaction fails", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		lastVersion := useCase.GetLastVersion()
		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}
		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, app.ErrInvalidVersion.Is(err))
		assert.Equal(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())
	})

	t.Run("Object version does not change when the idempotency fails, but global counter is changed", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NewAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}
		err := useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)

		lastVersion := useCase.GetLastVersion()
		idempotencyKey := uuid.New()
		e1, _ = entities.NewEntry(idempotencyKey, vos.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(idempotencyKey, vos.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		useCase.repository = &mocks.Repository{
			OnCreateTransaction: func(context.Context, *entities.Transaction) error {
				return app.ErrIdempotencyKeyViolation
			},
		}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.True(t, app.ErrIdempotencyKeyViolation.Is(err))
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		useCase.repository = &mocks.Repository{
			OnCreateTransaction: func(context.Context, *entities.Transaction) error {
				return nil
			},
		}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), entries)
		assert.Nil(t, err)
	})
}
