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
		acc, err := vos.NewAccountName("assets.bacen.conta_liquidacao.tesouraria")
		assert.Nil(t, err)
		expected := &Entry{
			ID:        newUUID,
			Operation: vos.CreditOperation,
			Account:   acc,
			Version:   vos.AnyAccountVersion,
			Amount:    123,
		}
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.AnyAccountVersion, 123)
		assert.Equal(t, expected, entry)
		assert.Nil(t, err)
	})

	t.Run("Invalid when entry id is invalid", func(t *testing.T) {
		entry, err := NewEntry(uuid.Nil, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.True(t, app.ErrInvalidEntryID.Is(err))
	})

	t.Run("Invalid when operation is invalid", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.InvalidOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.True(t, app.ErrInvalidOperation.Is(err))
	})

	t.Run("Invalid when amount is zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.AnyAccountVersion, 0)
		assert.Nil(t, entry)
		assert.True(t, app.ErrInvalidAmount.Is(err))
	})

	t.Run("Invalid when amount < zero", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen.conta_liquidacao.tesouraria", vos.AnyAccountVersion, -1)
		assert.Nil(t, entry)
		assert.True(t, app.ErrInvalidAmount.Is(err))
	})

	t.Run("Invalid when account structure has less than 3 levels", func(t *testing.T) {
		entry, err := NewEntry(newUUID, vos.CreditOperation, "assets.bacen", vos.AnyAccountVersion, 123)
		assert.Nil(t, entry)
		assert.True(t, app.ErrInvalidAccountStructure.Is(err))
	})
}
