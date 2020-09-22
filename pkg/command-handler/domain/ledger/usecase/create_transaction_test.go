package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_CreateTransaction(t *testing.T) {
	t.Run("Successfully creates a transaction with minimum inputs", func(t *testing.T) {
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, uuid.New(), entities.AnyAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, uuid.New(), entities.AnyAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := newFakeUseCase(nil).CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with new accounts", func(t *testing.T) {
		accountID1 := uuid.New()
		accountID2 := uuid.New()
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := newFakeUseCase(nil).CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with expected version", func(t *testing.T) {
		useCase := newFakeUseCase(nil)

		accountID1 := uuid.New()
		accountID2 := uuid.New()
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)

		assert.Nil(t, err)

		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.Nil(t, err)
	})

	t.Run("Fail with invalid expected version", func(t *testing.T) {
		useCase := newFakeUseCase(nil)

		accountID1 := uuid.New()
		accountID2 := uuid.New()
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)

		assert.Nil(t, err)

		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 4, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.True(t, errors.Is(err, entities.ErrInvalidVersion))
	})

	t.Run("When transaction fail, the counter isn't incremented", func(t *testing.T) {
		useCase := newFakeUseCase(nil)

		accountID1 := uuid.New()
		accountID2 := uuid.New()
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)

		assert.Nil(t, err)

		lastVersion := useCase.GetLastVersion()

		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 4, 123)
		entries = []entities.Entry{*e1, *e2}

		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.True(t, errors.Is(err, entities.ErrInvalidVersion))

		assert.Equal(t, lastVersion, useCase.GetLastVersion())
	})

	t.Run("Object version does not change when the transaction fails", func(t *testing.T) {
		useCase := newFakeUseCase(nil)

		lastVersion := useCase.GetLastVersion()
		accountID1 := uuid.New()
		accountID2 := uuid.New()
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}
		err := useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.True(t, errors.Is(err, entities.ErrInvalidVersion))
		assert.Equal(t, lastVersion, useCase.GetLastVersion())

		lastVersion = useCase.GetLastVersion()
		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 4, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.Nil(t, err)
		assert.NotEqual(t, lastVersion, useCase.GetLastVersion())
	})

	t.Run("Object version does not change when the idempotency fails, but global counter is changed", func(t *testing.T) {
		useCase := newFakeUseCase(nil)

		accountID1 := uuid.New()
		accountID2 := uuid.New()
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}
		err := useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.Nil(t, err)

		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.Nil(t, err)

		useCase.repository = &ledger.RepositoryMock{
			OnCreateTransaction: func(context.Context, *entities.Transaction) error {
				return entities.ErrIdempotencyKey
			},
		}

		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 2, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.True(t, errors.Is(err, entities.ErrInvalidVersion))

		useCase.repository = &ledger.RepositoryMock{
			OnCreateTransaction: func(context.Context, *entities.Transaction) error {
				return nil
			},
		}

		e1 = entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, 4, 123)
		e2 = entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{*e1, *e2}
		err = useCase.CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)
		assert.Nil(t, err)
	})
}

func newFakeUseCase(resultInMocks error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &ledger.RepositoryMock{}
	mockRepository.OnCreateTransaction = func(context.Context, *entities.Transaction) error {
		return resultInMocks
	}

	return NewLedgerUseCase(log, mockRepository)
}
