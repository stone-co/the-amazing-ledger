package usecases

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/mocks"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

const _company = "abc"

func TestLedgerUseCase_CreateTransaction(t *testing.T) {
	accountID1 := "liability.clients.available." + strings.ReplaceAll(uuid.New().String(), "-", "_")
	accountID2 := "liability.clients.available." + strings.ReplaceAll(uuid.New().String(), "-", "_")

	t.Run("Successfully creates a transaction with minimum inputs", func(t *testing.T) {
		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NextAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NextAccountVersion, 123)
		entries := []entities.Entry{e1, e2}

		tx, err := entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		err = newFakeCreateTransactionUseCase(nil).CreateTransaction(context.Background(), tx)
		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with new accounts", func(t *testing.T) {
		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NextAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NextAccountVersion, 123)
		entries := []entities.Entry{e1, e2}

		tx, err := entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		err = newFakeCreateTransactionUseCase(nil).CreateTransaction(context.Background(), tx)
		assert.Nil(t, err)
	})

	t.Run("Successfully creates a transaction with expected version", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NextAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NextAccountVersion, 123)
		entries := []entities.Entry{e1, e2}

		tx, err := entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		err = useCase.CreateTransaction(context.Background(), tx)
		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{e1, e2}

		tx, err = entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		err = useCase.CreateTransaction(context.Background(), tx)
		assert.Nil(t, err)
	})

	t.Run("Object version does not change when the idempotency fails, but global counter is changed", func(t *testing.T) {
		useCase := newFakeCreateTransactionUseCase(nil)

		e1, _ := entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, vos.NextAccountVersion, 123)
		e2, _ := entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, vos.NextAccountVersion, 123)
		entries := []entities.Entry{e1, e2}

		tx, err := entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		err = useCase.CreateTransaction(context.Background(), tx)
		assert.Nil(t, err)

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 2, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 3, 123)
		entries = []entities.Entry{e1, e2}

		tx, err = entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		err = useCase.CreateTransaction(context.Background(), tx)
		assert.Nil(t, err)

		idempotencyKey := uuid.New()
		e1, _ = entities.NewEntry(idempotencyKey, vos.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(idempotencyKey, vos.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{e1, e2}

		tx, err = entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		useCase.repository = &mocks.Repository{
			OnCreateTransaction: func(context.Context, entities.Transaction) error {
				return app.ErrIdempotencyKeyViolation
			},
		}
		err = useCase.CreateTransaction(context.Background(), tx)
		assert.True(t, app.ErrIdempotencyKeyViolation.Is(err))

		e1, _ = entities.NewEntry(uuid.New(), vos.DebitOperation, accountID1, 4, 123)
		e2, _ = entities.NewEntry(uuid.New(), vos.CreditOperation, accountID2, 5, 123)
		entries = []entities.Entry{e1, e2}

		tx, err = entities.NewTransaction(uuid.New(), entries...)
		assert.Nil(t, err)
		tx.Company = _company
		tx.Event = 1
		tx.CompetenceDate = time.Now()

		useCase.repository = &mocks.Repository{
			OnCreateTransaction: func(context.Context, entities.Transaction) error {
				return nil
			},
		}
		err = useCase.CreateTransaction(context.Background(), tx)
		assert.Nil(t, err)
	})
}
