package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	newUUID := uuid.New()

	t.Run("Successfully creates an entry with minimum inputs", func(t *testing.T) {
		acc, err := NewAccountName("assets:bacen:conta_liquidacao:tesouraria")
		assert.Nil(t, err)
		expected := &Entry{
			ID:        newUUID,
			Operation: CreditOperation,
			Account:   acc,
			Version:   AnyAccountVersion,
			Amount:    123,
		}
		entry, err := NewEntry(newUUID, CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", AnyAccountVersion, 123)
		assert.Equal(t, expected, entry)
		assert.Nil(t, err)
	})

	t.Run("Invalid when entry id is invalid", func(t *testing.T) {
		entry, err := NewEntry(uuid.Nil, CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.Equal(t, ErrInvalidData, err)
	})

	t.Run("Invalid when operation is invalid", func(t *testing.T) {
		entry, err := NewEntry(newUUID, InvalidOperation, "assets:bacen:conta_liquidacao:tesouraria", AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.Equal(t, ErrInvalidData, err)
	})

	t.Run("Invalid when amount is zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", AnyAccountVersion, 0)
		assert.Nil(t, entry)
		assert.Equal(t, ErrInvalidData, err)
	})

	t.Run("Invalid when amount < zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, CreditOperation, "assets:bacen:conta_liquidacao:tesouraria", AnyAccountVersion, -1)
		assert.Nil(t, entry)
		assert.Equal(t, ErrInvalidData, err)
	})

	t.Run("Invalid when account structure has less than 4 levels", func(t *testing.T) {
		entry, err := NewEntry(newUUID, CreditOperation, "assets:bacen:conta_liquidacao", AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.Equal(t, ErrInvalidAccountStructure, err)
	})
}
