package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCachedAccounts_LoadOrStore(t *testing.T) {
	c := NewCachedAccounts()

	t.Run("New accounts started with version 1", func(t *testing.T) {
		accountInfo := c.LoadOrStore("account/123")
		assert.Equal(t, NewAccountVersion, accountInfo.Version)
	})

	t.Run("Account info is saved successfully", func(t *testing.T) {
		accountID := "account/321/"
		var version Version = 1234
		c.Store(accountID, version)
		accountInfo := c.LoadOrStore(accountID)
		assert.Equal(t, version, accountInfo.Version)
	})
}
