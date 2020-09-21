package usecase

import (
	"context"
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

		err := newFakeUseCase().CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with new accounts", func(t *testing.T) {
		accountID1 := uuid.New()
		accountID2 := uuid.New()
		e1 := entities.NewEntry(uuid.New(), entities.DebitOperation, accountID1, entities.NewAccountVersion, 123)
		e2 := entities.NewEntry(uuid.New(), entities.CreditOperation, accountID2, entities.NewAccountVersion, 123)
		entries := []entities.Entry{*e1, *e2}

		err := newFakeUseCase().CreateTransaction(context.Background(), uuid.New(), time.Now().UTC(), entries)

		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with expected version", func(t *testing.T) {
		useCase := newFakeUseCase()

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
		useCase := newFakeUseCase()

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
		assert.Error(t, err, entities.ErrInvalidVersion)
	})
}

func newFakeUseCase() *LedgerUseCase {
	log := logrus.New()

	mockRepository := &ledger.RepositoryMock{}
	mockRepository.OnCreateTransaction = func(context.Context, *entities.Transaction) error {
		return nil
	}

	return NewLedgerUseCase(log, mockRepository)
}
