package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/errors"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	newUUID := uuid.New()

	t.Run("Successfully creates an entry with minimum inputs", func(t *testing.T) {
		acc, err := vo.NewAccountName("assets:bacen:conta_liquidacao:tesouraria")
		assert.Nil(t, err)
		expected := &Entry{
			ID:        newUUID,
			Operation: vo.CreditOperation,
			Account:   acc,
			Version:   vo.AnyAccountVersion,
			Amount:    123,
		}
		entry, err := NewEntry(newUUID, vo.CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", vo.AnyAccountVersion, 123)
		assert.Equal(t, expected, entry)
		assert.Nil(t, err)
	})

	t.Run("Invalid when entry id is invalid", func(t *testing.T) {
		entry, err := NewEntry(uuid.Nil, vo.CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", vo.AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.True(t, errors.ErrInvalidData.Is(err))
	})

	t.Run("Invalid when operation is invalid", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vo.InvalidOperation, "assets:bacen:conta_liquidacao:tesouraria", vo.AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.True(t, errors.ErrInvalidData.Is(err))
	})

	t.Run("Invalid when amount is zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vo.CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", vo.AnyAccountVersion, 0)
		assert.Nil(t, entry)
		assert.True(t, errors.ErrInvalidData.Is(err))
	})

	t.Run("Invalid when amount < zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vo.CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", vo.AnyAccountVersion, -1)
		assert.Nil(t, entry)
		assert.True(t, errors.ErrInvalidData.Is(err))
	})

	t.Run("Invalid when account structure has less than 4 levels", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vo.CreditOperation, "assets:bacen:conta_liquidacao", vo.AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.True(t, errors.ErrInvalidAccountStructure.Is(err))
	})
}
