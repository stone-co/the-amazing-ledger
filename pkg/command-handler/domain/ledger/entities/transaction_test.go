package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	id := uuid.New()

	e11, _ := NewEntry(uuid.New(), DebitOperation, "liability:clients:available:111", AnyAccountVersion, 123)
	e12, _ := NewEntry(uuid.New(), CreditOperation, "liability:clients:available:222", AnyAccountVersion, 123)
	validTwoEntries := []Entry{*e11, *e12}

	e21, _ := NewEntry(uuid.New(), DebitOperation, "liability:clients:available:333", AnyAccountVersion, 400)
	e22, _ := NewEntry(uuid.New(), CreditOperation, "liability:clients:available:444", AnyAccountVersion, 300)
	e23, _ := NewEntry(uuid.New(), CreditOperation, "liability:clients:available:555", AnyAccountVersion, 100)
	validThreeEntries := []Entry{*e21, *e22, *e23}

	t.Run("Invalid entries number when the transaction has no entries", func(t *testing.T) {
		got, err := NewTransaction(id)
		assert.True(t, errors.Is(err, ErrInvalidEntriesNumber))
		assert.Nil(t, got)
	})

	t.Run("Invalid entries number when the transaction has 1 entry", func(t *testing.T) {
		got, err := NewTransaction(id, *e11)
		assert.True(t, errors.Is(err, ErrInvalidEntriesNumber))
		assert.Nil(t, got)
	})

	t.Run("Valid transaction with 2 entries", func(t *testing.T) {
		got, err := NewTransaction(id, validTwoEntries...)
		assert.Nil(t, err)
		assert.Equal(t, &Transaction{
			ID:      id,
			Entries: validTwoEntries,
		}, got)
	})

	t.Run("Valid transaction with 3 entries", func(t *testing.T) {
		got, err := NewTransaction(id, validThreeEntries...)
		assert.Nil(t, err)
		assert.Equal(t, &Transaction{
			ID:      id,
			Entries: validThreeEntries,
		}, got)
	})

	t.Run("Invalid transaction with 2 entries and balance != 0", func(t *testing.T) {
		e1, _ := NewEntry(uuid.New(), DebitOperation, "liability:clients:available:111", AnyAccountVersion, 123)
		e2, _ := NewEntry(uuid.New(), CreditOperation, "liability:clients:available:222", AnyAccountVersion, 234)
		got, err := NewTransaction(id, *e1, *e2)
		assert.True(t, errors.Is(err, ErrInvalidBalance))
		assert.Nil(t, got)
	})

	t.Run("Invalid transaction with 3 entries and balance != 0", func(t *testing.T) {
		e1, _ := NewEntry(uuid.New(), DebitOperation, "liability:clients:available:111", AnyAccountVersion, 400)
		e2, _ := NewEntry(uuid.New(), CreditOperation, "liability:clients:available:222", AnyAccountVersion, 200)
		e3, _ := NewEntry(uuid.New(), CreditOperation, "liability:clients:available:333", AnyAccountVersion, 100)
		got, err := NewTransaction(id, *e1, *e2, *e3)
		assert.True(t, errors.Is(err, ErrInvalidBalance))
		assert.Nil(t, got)
	})

	t.Run("Invalid transaction with empty ID", func(t *testing.T) {
		got, err := NewTransaction(uuid.Nil, validTwoEntries...)
		assert.True(t, errors.Is(err, ErrInvalidData))
		assert.Nil(t, got)
	})
}
