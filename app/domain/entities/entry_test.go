package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestNewEntry(t *testing.T) {
	newUUID := uuid.New()

	t.Run("Successfully creates an entry with minimum inputs", func(t *testing.T) {
		acc, err := vos.NewAccountPath("assets.bacen.conta_liquidacao.tesouraria")
		assert.Nil(t, err)
		expected := Entry{
			ID:        newUUID,
			Operation: vos.CreditOperation,
			Account:   acc,
			Version:   vos.NextAccountVersion,
			Amount:    123,
		}
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.NextAccountVersion, 123)
		assert.Equal(t, expected, entry)
		assert.Nil(t, err)
	})

	t.Run("Invalid when entry id is invalid", func(t *testing.T) {
		entry, err := NewEntry(uuid.Nil, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.NextAccountVersion, 123)
		assert.Empty(t, entry)
		assert.ErrorIs(t, err, app.ErrInvalidEntryID)
	})

	t.Run("Invalid when operation is invalid", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.InvalidOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.NextAccountVersion, 123)
		assert.Empty(t, entry)
		assert.ErrorIs(t, err, app.ErrInvalidOperation)
	})

	t.Run("Invalid when amount is zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.NextAccountVersion, 0)
		assert.Empty(t, entry)
		assert.ErrorIs(t, err, app.ErrInvalidAmount)
	})

	t.Run("Invalid when amount < zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.NextAccountVersion, -1)
		assert.Empty(t, entry)
		assert.ErrorIs(t, err, app.ErrInvalidAmount)
	})

	t.Run("Invalid when account structure has less than 3 levels", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen", vos.NextAccountVersion, 123)
		assert.Empty(t, entry)
		assert.ErrorIs(t, err, app.ErrInvalidAccountStructure)
	})
}
